package f1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/AidanFogarty/pitwall/internal/logging"
	"github.com/cenkalti/backoff/v4"
	"github.com/philippseith/signalr"
)

var (
	negotiateURL = "https://livetiming.formula1.com/signalrcore/negotiate"
	signalrURL   = "https://livetiming.formula1.com/signalrcore"
)

const (
	maxMessageSize   = 10 * 1024 * 1024 // 10mb
	connectTimeout   = 30 * time.Second
	subscribeTimeout = 30 * time.Second
	reconnectDelay   = time.Second
)

type FeedHandler func(topic string, data any, timestamp string) error

type F1Receiver struct {
	signalr.Hub
	handler FeedHandler
	log     *logging.Logger

	events    atomic.Int64
	lastEvent atomic.Int64
	lastTopic atomic.Value
}

func (r *F1Receiver) Feed(topic string, data any, timestamp string) {
	r.events.Add(1)
	r.lastEvent.Store(time.Now().UnixNano())
	r.lastTopic.Store(topic)

	if r.handler == nil {
		return
	}

	if err := r.handler(topic, data, timestamp); err != nil {
		r.log.Logf("handler error topic=%s err=%v", topic, err)
	}
}

func (r *F1Receiver) sinceLastEvent() time.Duration {
	last := r.lastEvent.Load()
	if last == 0 {
		return -1
	}

	return time.Since(time.Unix(0, last))
}

type F1LiveClient struct {
	topics  []string
	timeout time.Duration

	handler  FeedHandler
	receiver *F1Receiver
	client   signalr.Client
	log      *logging.Logger

	epoch           atomic.Int64
	subscribedEpoch int64
}

func NewF1LiveClient(topics []string) *F1LiveClient {
	return &F1LiveClient{
		topics:  topics,
		timeout: 60 * time.Second,
		log:     logging.NewDiscard(),
	}
}

func (c *F1LiveClient) SetLogFile(path string) error {
	log, err := logging.NewFile(path)
	if err != nil {
		return err
	}

	c.log = log
	return nil
}

func (c *F1LiveClient) SetHandler(handler FeedHandler) {
	// TODO: Need to refactor so this set handler method is not needed
	c.handler = handler
	if c.receiver != nil {
		c.receiver.handler = handler
	}
}

func (c *F1LiveClient) Start(ctx context.Context) (*InitialLiveState, error) {
	c.log.Logf("starting topics=%d", len(c.topics))

	c.receiver = &F1Receiver{log: c.log}

	client, err := signalr.NewClient(ctx,
		signalr.WithConnector(func() (signalr.Connection, error) {
			return c.connect(ctx)
		}),
		signalr.WithReceiver(c.receiver),
		signalr.MaximumReceiveMessageSize(maxMessageSize),
		signalr.Logger(c.log, false),
		signalr.WithBackoff(func() backoff.BackOff {
			return backoff.WithMaxRetries(backoff.NewConstantBackOff(reconnectDelay), 1)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	c.client = client

	client.Start()

	if err := c.waitConnected(ctx); err != nil {
		return nil, err
	}

	epoch := c.epoch.Load()
	initialState, err := c.subscribe(ctx)
	if err != nil {
		return nil, err
	}

	c.subscribedEpoch = epoch
	c.log.Logf("subscribed epoch=%d", epoch)

	return initialState, nil
}

func (c *F1LiveClient) waitConnected(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	err := <-c.client.WaitForState(ctx, signalr.ClientConnected)
	if err == nil {
		return nil
	}

	if clientErr := c.client.Err(); clientErr != nil {
		err = clientErr
	}
	c.client.Stop()
	c.log.Logf("connect failed err=%v", err)

	return fmt.Errorf("failed to connect: %w", err)
}

func (c *F1LiveClient) connect(ctx context.Context) (signalr.Connection, error) {
	epoch := c.epoch.Add(1)
	c.log.Logf("connecting epoch=%d", epoch)

	awsalbcors, err := c.getAWSCookie()
	if err != nil {
		c.log.Logf("cookie failed epoch=%d err=%v", epoch, err)
		return nil, fmt.Errorf("failed to get cookie: %w", err)
	}

	connection, err := signalr.NewHTTPConnection(
		ctx,
		signalrURL,
		signalr.WithHTTPHeaders(func() http.Header {
			header := http.Header{}
			header.Set("Cookie", fmt.Sprintf("AWSALBCORS=%s", awsalbcors))
			return header
		}),
	)
	if err != nil {
		c.log.Logf("connection failed epoch=%d err=%v", epoch, err)
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return connection, nil
}

func (c *F1LiveClient) subscribe(ctx context.Context) (*InitialLiveState, error) {
	select {
	case result := <-c.client.Invoke("Subscribe", c.topics):
		if result.Error != nil {
			return nil, fmt.Errorf("failed to subscribe: %w", result.Error)
		}
		return parseInitialState(result.Value)
	case <-time.After(subscribeTimeout):
		return nil, fmt.Errorf("subscribe timed out after %s", subscribeTimeout)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func parseInitialState(value any) (*InitialLiveState, error) {
	if value == nil {
		return nil, fmt.Errorf("no initial live state")
	}

	stateMap, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid initial state format")
	}

	// Remove the _kf field from DriverList if it exists
	if driverList, ok := stateMap["DriverList"].(map[string]any); ok {
		delete(driverList, "_kf")
	}

	jsonData, err := json.Marshal(stateMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal initial state: %w", err)
	}

	var initialState InitialLiveState
	if err := json.Unmarshal(jsonData, &initialState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state: %w", err)
	}

	return &initialState, nil
}

func (c *F1LiveClient) getAWSCookie() (string, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequest("OPTIONS", negotiateURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "AWSALBCORS" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("AWSALBCORS cookie not found")
}
