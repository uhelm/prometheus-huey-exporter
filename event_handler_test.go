package exporter

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

func Test_eventHandler_HandleEvent(t *testing.T) {
	m := SetupMetrics("foo")
	logger := slog.Default()

	h := NewEventHandler(m, logger)

	t.Run("malformed event", func(t *testing.T) {
		msg := &redis.Message{Payload: "something"}

		err := h.HandleEvent(msg)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run(string(Executing), func(t *testing.T) {

	})

	t.Run(string(Complete), func(t *testing.T) {

	})

	t.Run(string(Error), func(t *testing.T) {

	})

	t.Run(string(Locked), func(t *testing.T) {

	})
}
