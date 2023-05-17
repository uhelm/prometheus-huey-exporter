package exporter

import (
	"encoding/json"
	"testing"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

func Test_eventHandler_HandleEvent(t *testing.T) {
	m := SetupMetrics("foo")
	logger := slog.Default()

	h := NewEventHandler(m, logger)

	evt := event{
		TaskName: "my-task",
		TaskId:   "aabbccddee",
	}

	t.Run("malformed event", func(t *testing.T) {
		msg := &redis.Message{Payload: "something"}

		err := h.HandleEvent(msg)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run(string(executingEvent), func(t *testing.T) {
		evt.Event = executingEvent

		msg := evtToMessage(evt)

		_ = h.HandleEvent(msg)
	})

	t.Run(string(completeEvent), func(t *testing.T) {
		evt.Event = executingEvent

		msg := evtToMessage(evt)

		_ = h.HandleEvent(msg)
	})

	t.Run(string(errorEvent), func(t *testing.T) {
		evt.Event = executingEvent

		msg := evtToMessage(evt)

		_ = h.HandleEvent(msg)
	})

	t.Run(string(lockedEvent), func(t *testing.T) {
		evt.Event = executingEvent

		msg := evtToMessage(evt)

		_ = h.HandleEvent(msg)
	})
}

func evtToMessage(evt event) *redis.Message {
	data, err := json.Marshal(evt)
	if err != nil {
		return nil
	}
	return &redis.Message{
		Payload: string(data),
	}
}
