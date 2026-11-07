package f1

import (
	"context"
	"fmt"
	"time"

	"github.com/philippseith/signalr"
)

const (
	statusInterval    = 10 * time.Second
	reconcileInterval = 2 * time.Second
	staleThreshold    = 30 * time.Second
)

func (c *F1LiveClient) Watch(ctx context.Context, onResync func(*InitialLiveState), onState func(string)) {
	states := make(chan signalr.ClientState, 16)
	cancel := c.client.ObserveStateChanged(states)
	defer cancel()

	statusTicker := time.NewTicker(statusInterval)
	defer statusTicker.Stop()

	reconcileTicker := time.NewTicker(reconcileInterval)
	defer reconcileTicker.Stop()

	onState(stateName(c.client.State()))

	var lastCount int64
	for {
		select {
		case <-ctx.Done():
			c.log.Logf("watch stopping err=%v", ctx.Err())
			return
		case state := <-states:
			c.log.Logf("client state=%s err=%v", stateName(state), c.client.Err())
			onState(stateName(state))
			c.resubscribe(ctx, onResync)
		case <-reconcileTicker.C:
			c.resubscribe(ctx, onResync)
		case <-statusTicker.C:
			lastCount = c.logStatus(lastCount)
		}
	}
}

func (c *F1LiveClient) logStatus(lastCount int64) int64 {
	count := c.receiver.events.Load()
	topic, _ := c.receiver.lastTopic.Load().(string)
	since := c.receiver.sinceLastEvent()
	state := stateName(c.client.State())

	c.log.Logf("status state=%s err=%v epoch=%d subscribed_epoch=%d events_total=%d events_delta=%d since_last_event=%s last_topic=%s",
		state, c.client.Err(), c.epoch.Load(), c.subscribedEpoch,
		count, count-lastCount, since.Round(time.Millisecond), topic)

	if since > staleThreshold {
		c.log.Logf("WARN no events for %s state=%s", since.Round(time.Second), state)
	}

	return count
}

func (c *F1LiveClient) resubscribe(ctx context.Context, onResync func(*InitialLiveState)) {
	if c.client.State() != signalr.ClientConnected {
		return
	}

	epoch := c.epoch.Load()
	if epoch == c.subscribedEpoch {
		return
	}

	c.log.Logf("resubscribing epoch=%d", epoch)
	state, err := c.subscribe(ctx)
	if err != nil {
		c.log.Logf("resubscribe failed epoch=%d err=%v", epoch, err)
		return
	}

	c.subscribedEpoch = epoch
	c.log.Logf("resubscribed epoch=%d", epoch)
	onResync(state)
}

func stateName(state signalr.ClientState) string {
	switch state {
	case signalr.ClientCreated:
		return "created"
	case signalr.ClientConnecting:
		return "connecting"
	case signalr.ClientConnected:
		return "connected"
	case signalr.ClientClosed:
		return "closed"
	}

	return fmt.Sprintf("unknown(%d)", int(state))
}
