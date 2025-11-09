package f1

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type F1JsonConfig struct {
	DataDir string
	Session string

	// For debugging replay sessions
	SkipDelayEventCount int
}

func NewF1JsonConfig(dataDir string, session string) *F1JsonConfig {
	return &F1JsonConfig{
		DataDir: dataDir,
		Session: session,
	}
}

type F1EventHandler func(ctx context.Context, data F1Event) error

type F1JsonClient struct {
	config  *F1JsonConfig
	handler F1EventHandler
}

func NewF1JsonClient(config *F1JsonConfig, handler F1EventHandler) *F1JsonClient {
	return &F1JsonClient{
		config:  config,
		handler: handler,
	}
}

func (client *F1JsonClient) Start(ctx context.Context) error {
	events, err := client.parseFile()
	if err != nil {
		return err
	}

	return client.replayEvents(ctx, events)
}

func (client *F1JsonClient) parseFile() ([]F1Event, error) {
	filePath := fmt.Sprintf("%s/%s/live.txt", client.config.DataDir, client.config.Session)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open session file: %s: %w", filePath, err)
	}
	defer file.Close()

	var events []F1Event
	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++

		line := scanner.Text()
		var event F1Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			fmt.Printf("warning: failed to parse line %v", line)
		}
		events = append(events, event)
	}
	return events, nil
}

func (client *F1JsonClient) findSessionStartIndex(events []F1Event) int {
	for i, event := range events {
		if event.Type != "SessionStatus" {
			continue
		}

		var sessionStatusData SessionStatus
		if err := json.Unmarshal(event.Data, &sessionStatusData); err != nil {
			continue
		}

		if sessionStatusData.Status == "Started" {
			return i
		}
	}

	return -1
}

func (client *F1JsonClient) replayEvents(ctx context.Context, events []F1Event) error {
	if len(events) == 0 {
		return fmt.Errorf("no events to replay")
	}

	adjustedStartTime := client.fastForwardEvents(ctx, events)
	return client.processTimedEvents(ctx, events, adjustedStartTime)
}

func (client *F1JsonClient) fastForwardEvents(ctx context.Context, events []F1Event) time.Time {
	skipCount := client.config.SkipDelayEventCount

	if skipCount == -1 {
		skipCount = client.findSessionStartIndex(events)
		if skipCount == -1 {
			return time.Now()
		}
	}

	for i := 0; i <= skipCount; i++ {
		if err := client.processEvent(ctx, events[i]); err != nil {
			fmt.Printf("Warning: error processing event %d: %v\n", i, err)
		}
	}

	lastSkippedEvent := events[skipCount]
	return time.Now().Add(-time.Duration(lastSkippedEvent.Offset) * time.Millisecond)
}

func (client *F1JsonClient) processTimedEvents(ctx context.Context, events []F1Event, startTime time.Time) error {
	sessionStartIdx := client.findSessionStartIndex(events)
	startIndex := sessionStartIdx + 1

	if sessionStartIdx == -1 {
		startIndex = 0
	}

	if startIndex >= len(events) {
		return nil
	}

	for i := startIndex; i < len(events); i++ {
		event := events[i]

		timeSinceStart := time.Duration(event.Offset) * time.Millisecond
		timeToWait := timeSinceStart - time.Since(startTime)

		timer := time.NewTimer(timeToWait)
		<-timer.C

		if err := client.processEvent(ctx, event); err != nil {
			fmt.Printf("Warning: error processing event %d: %v\n", i, err)
		}
	}

	return nil
}

func (client *F1JsonClient) processEvent(ctx context.Context, event F1Event) error {
	err := client.handler(ctx, event)
	if err != nil {
		return fmt.Errorf("event handler failed for %s: %w", event.Type, err)
	}

	return nil
}
