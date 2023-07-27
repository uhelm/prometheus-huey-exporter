package exporter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

// An EventHandler is responsible to handle all the messages received from
// a [EventListener]
type EventHandler interface {
	HandleEvent(msg *redis.Message) error
}

// NewEventHandler creates and EventHandler that will handle the message
// updating the metrics
func NewEventHandler(metrics *Metrics, logger *slog.Logger) EventHandler {
	eHandler := &eventHandler{
		metrics:           metrics,
		startExecutionMap: make(map[string]time.Time),
	}

	return &evtLoggingMdw{
		logger: logger,
		next:   eHandler,
	}
}

type eventHandler struct {
	metrics           *Metrics
	startExecutionMap map[string]time.Time
}

func (h *eventHandler) HandleEvent(msg *redis.Message) error {
	var evt event
	if err := json.Unmarshal([]byte(msg.Payload), &evt); err != nil {
		return err
	}

	switch evt.Event {
	case executingEvent:
		h.metrics.Executions.WithLabelValues(evt.TaskName).Inc()
		h.startExecutionMap[evt.TaskId] = time.Now()

	case completeEvent, errorEvent:
		labelValues := []string{evt.TaskName, fmt.Sprint(evt.Event == completeEvent)}
		h.metrics.Completed.WithLabelValues(labelValues...).Inc()
		if startTime, ok := h.startExecutionMap[evt.TaskId]; ok {
			duration := time.Since(startTime).Seconds()
			h.metrics.Duration.WithLabelValues(labelValues...).Observe(duration)
			h.metrics.LastDuration.WithLabelValues(labelValues...).Set(duration)
			delete(h.startExecutionMap, evt.TaskId)
		}

	case lockedEvent:
		h.metrics.Locked.WithLabelValues(evt.TaskName).Inc()
	}
	return nil
}

type eventType string

const (
	canceledEvent    eventType = "canceled"
	completeEvent    eventType = "complete"
	errorEvent       eventType = "error"
	executingEvent   eventType = "executing"
	expiredEvent     eventType = "expired"
	lockedEvent      eventType = "locked"
	retryingEvent    eventType = "retrying"
	revokedEvent     eventType = "revoked"
	scheduledEvent   eventType = "scheduled"
	interruptedEvent eventType = "interrupted"
)

type event struct {
	Event    eventType `json:"event,omitempty"`
	TaskName string    `json:"task_name,omitempty"`
	TaskId   string    `json:"task_id,omitempty"`
}

type evtLoggingMdw struct {
	logger *slog.Logger
	next   EventHandler
}

func (mdw *evtLoggingMdw) HandleEvent(msg *redis.Message) error {
	mdw.logger.Debug(fmt.Sprintf("event: %+v", msg))

	return mdw.next.HandleEvent(msg)
}
