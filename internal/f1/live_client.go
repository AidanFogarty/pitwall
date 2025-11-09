package f1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/philippseith/signalr"
)

const (
	negotiateURL   = "https://livetiming.formula1.com/signalrcore/negotiate"
	signalrURL     = "https://livetiming.formula1.com/signalrcore"
	maxMessageSize = 10 * 1024 * 1024 // 10mb
)

type FeedHandler func(topic string, data any, timestamp string) error

type F1Receiver struct {
	signalr.Hub
	handler FeedHandler
}

func (r *F1Receiver) Receive(topic string, data any, timestamp string) {
	if r.handler != nil {
		r.handler(topic, data, timestamp)
	}
}

type F1LiveClient struct {
	topics  []string
	timeout time.Duration

	handler  FeedHandler
	receiver *F1Receiver
}

func NewF1LiveClient(topics []string) *F1LiveClient {
	return &F1LiveClient{
		topics:  topics,
		timeout: 60 * time.Second,
	}
}

func (c *F1LiveClient) SetHandler(handler FeedHandler) {
	// TODO: Need to refactor so this set handler method is not needed
	c.handler = handler
	if c.receiver != nil {
		c.receiver.handler = handler
	}
}

func (c *F1LiveClient) Start(ctx context.Context) (*InitialLiveState, error) {
	awsalbcors, err := c.getAWSCookie()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookie: %w", err)
	}

	receiver := &F1Receiver{}
	c.receiver = receiver

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
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	client, err := signalr.NewClient(ctx,
		signalr.WithConnection(connection),
		signalr.WithReceiver(receiver),
		signalr.MaximumReceiveMessageSize(maxMessageSize),
		signalr.Logger(nil, false),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	client.Start()

	resultCh := client.Invoke("Subscribe", c.topics)
	result := <-resultCh

	if result.Error != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", result.Error)
	}

	if result.Value == nil {
		return nil, fmt.Errorf("no initial live state")
	}

	stateMap, ok := result.Value.(map[string]any)
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
