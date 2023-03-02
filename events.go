package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis/v8"
)

type EventListener interface {
	Run(ctx context.Context) error
}

type listener struct {
	rdb               *redis.Client
	channel           string
	logger            log.Logger
	m                 *metrics
	startExecutionMap map[string]time.Time
}

func NewEventListener(rdb *redis.Client, channel string, logger log.Logger, m *metrics) EventListener {
	return &listener{
		rdb:               rdb,
		channel:           channel,
		logger:            logger,
		m:                 m,
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
				level.Error(l.logger).Log("during", "handleEvent", "err", err)
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

	switch evt.Event {
	case SIGNAL_EXECUTING:
		l.m.executions.WithLabelValues(evt.TaskName).Inc()
		l.startExecutionMap[evt.TaskId] = time.Now()

	case SIGNAL_COMPLETE, SIGNAL_ERROR:
		labelValues := []string{evt.TaskName, fmt.Sprint(evt.Event == SIGNAL_COMPLETE)}
		l.m.completed.WithLabelValues(labelValues...).Inc()
		if startTime, ok := l.startExecutionMap[evt.TaskId]; ok {
			l.m.duration.WithLabelValues(labelValues...).Observe(time.Since(startTime).Seconds())
			delete(l.startExecutionMap, evt.TaskId)
		}

	case SIGNAL_LOCKED:
		l.m.locked.WithLabelValues(evt.TaskName).Inc()
	}

	level.Info(l.logger).Log("map", fmt.Sprintf("%+v", l.startExecutionMap))

	return nil
}

type eventType string

const (
	SIGNAL_CANCELED    eventType = "canceled"
	SIGNAL_COMPLETE    eventType = "complete"
	SIGNAL_ERROR       eventType = "error"
	SIGNAL_EXECUTING   eventType = "executing"
	SIGNAL_EXPIRED     eventType = "expired"
	SIGNAL_LOCKED      eventType = "locked"
	SIGNAL_RETRYING    eventType = "retrying"
	SIGNAL_REVOKED     eventType = "revoked"
	SIGNAL_SCHEDULED   eventType = "scheduled"
	SIGNAL_INTERRUPTED eventType = "interrupted"
)

type event struct {
	Event    eventType `json:"event,omitempty"`
	TaskName string    `json:"task_name,omitempty"`
	TaskId   string    `json:"task_id,omitempty"`
}
