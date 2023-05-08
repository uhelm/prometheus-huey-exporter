package exporter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type EventHandler interface {
	HandleEvent(msg *redis.Message) error
}

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
	case Executing:
		h.metrics.Executions.WithLabelValues(evt.TaskName).Inc()
		h.startExecutionMap[evt.TaskId] = time.Now()

	case Complete, Error:
		labelValues := []string{evt.TaskName, fmt.Sprint(evt.Event == Complete)}
		h.metrics.Completed.WithLabelValues(labelValues...).Inc()
		if startTime, ok := h.startExecutionMap[evt.TaskId]; ok {
			h.metrics.Duration.WithLabelValues(labelValues...).Observe(time.Since(startTime).Seconds())
			delete(h.startExecutionMap, evt.TaskId)
		}

	case Locked:
		h.metrics.Locked.WithLabelValues(evt.TaskName).Inc()
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

type evtLoggingMdw struct {
	logger *slog.Logger
	next   EventHandler
}

func (mdw *evtLoggingMdw) HandleEvent(msg *redis.Message) error {
	mdw.logger.Debug(fmt.Sprintf("event: %+v", msg))

	return mdw.next.HandleEvent(msg)
}
