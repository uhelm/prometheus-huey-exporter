package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type EventListener interface {
	Run(ctx context.Context) error
}

type listener struct {
	rdb               *redis.Client
	channel           string
	logger            *slog.Logger
	m                 *Metrics
	startExecutionMap map[string]time.Time
}

func NewEventListener(rdb *redis.Client, channel string, logger *slog.Logger, metrics *Metrics) EventListener {
	return &listener{
		rdb:               rdb,
		channel:           channel,
		logger:            logger,
		m:                 metrics,
		startExecutionMap: make(map[string]time.Time),
	}
}

func (l *listener) Run(ctx context.Context) error {
	pubsub := l.rdb.Subscribe(ctx, l.channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}

	ch := pubsub.Channel()

	for {
		select {
		case evt := <-ch:
			if err = l.handleEvent(evt); err != nil {
				l.logger.Error("running handleEvent()", "err", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (l *listener) handleEvent(msg *redis.Message) error {
	var evt event
	if err := json.Unmarshal([]byte(msg.Payload), &evt); err != nil {
		return err
	}

	l.logger.Debug(fmt.Sprintf("event: %+v", evt))

	switch evt.Event {
	case Executing:
		l.m.Executions.WithLabelValues(evt.TaskName).Inc()
		l.startExecutionMap[evt.TaskId] = time.Now()

	case Complete, Error:
		labelValues := []string{evt.TaskName, fmt.Sprint(evt.Event == Complete)}
		l.m.Completed.WithLabelValues(labelValues...).Inc()
		if startTime, ok := l.startExecutionMap[evt.TaskId]; ok {
			l.m.Duration.WithLabelValues(labelValues...).Observe(time.Since(startTime).Seconds())
			delete(l.startExecutionMap, evt.TaskId)
		}

	case Locked:
		l.m.Locked.WithLabelValues(evt.TaskName).Inc()
	}
	return nil
}

type eventType string

const (
	Canceled    eventType = "canceled"
	Complete    eventType = "complete"
	Error       eventType = "error"
	Executing   eventType = "executing"
	Expired     eventType = "expired"
	Locked      eventType = "locked"
	Retrying    eventType = "retrying"
	Revoked     eventType = "revoked"
	Scheduled   eventType = "scheduled"
	Interrupted eventType = "interrupted"
)

type event struct {
	Event    eventType `json:"event,omitempty"`
	TaskName string    `json:"task_name,omitempty"`
	TaskId   string    `json:"task_id,omitempty"`
}
