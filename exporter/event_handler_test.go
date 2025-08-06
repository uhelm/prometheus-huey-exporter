package exporter

import (
	"encoding/json"
	"testing"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/redis/go-redis/v9"
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

		want := testutil.ToFloat64(m.Executions.WithLabelValues(evt.TaskName)) + 1

		_ = h.HandleEvent(msg)

		got := testutil.ToFloat64(m.Executions.WithLabelValues(evt.TaskName))

		if got != want {
			t.Errorf("executing: want:%v, got:%v", want, got)
		}
	})

	t.Run(string(completeEvent), func(t *testing.T) {
		evt.Event = completeEvent

		msg := evtToMessage(evt)

		successWant := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "true")) + 1
		failureWant := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "false"))

		_ = h.HandleEvent(msg)

		successGot := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "true"))
		failureGot := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "false"))

		if successGot != successWant {
			t.Errorf("completed: successGot:%v, successWant:%v", successWant, successGot)
		}

		if failureGot != failureWant {
			t.Errorf("completed: failureGot:%v, failureWant:%v", successWant, successGot)
		}
	})

	t.Run(string(errorEvent), func(t *testing.T) {
		evt.Event = errorEvent

		msg := evtToMessage(evt)

		successWant := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "true"))
		failureWant := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "false")) + 1

		_ = h.HandleEvent(msg)

		successGot := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "true"))
		failureGot := testutil.ToFloat64(m.Completed.WithLabelValues(evt.TaskName, "false"))

		if successGot != successWant {
			t.Errorf("completed: successGot:%v, successWant:%v", successWant, successGot)
		}

		if failureGot != failureWant {
			t.Errorf("completed: failureGot:%v, failureWant:%v", successWant, successGot)
		}
	})

	t.Run(string(lockedEvent), func(t *testing.T) {
		evt.Event = lockedEvent

		msg := evtToMessage(evt)

		want := testutil.ToFloat64(m.Locked.WithLabelValues(evt.TaskName)) + 1

		_ = h.HandleEvent(msg)

		got := testutil.ToFloat64(m.Locked.WithLabelValues(evt.TaskName))

		if got != want {
			t.Errorf("executing: want:%v, got:%v", want, got)
		}
	})

	t.Run(string(canceledEvent), func(t *testing.T) {
		evt.Event = canceledEvent

		msg := evtToMessage(evt)

		want := testutil.ToFloat64(m.Canceled.WithLabelValues(evt.TaskName)) + 1

		_ = h.HandleEvent(msg)

		got := testutil.ToFloat64(m.Canceled.WithLabelValues(evt.TaskName))

		if got != want {
			t.Errorf("executing: want:%v, got:%v", want, got)
		}
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
