package exporter

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type mockHandler struct {
	payloads []string
}

func (h *mockHandler) HandleEvent(msg *redis.Message) error {
	h.payloads = append(h.payloads, msg.Payload)
	return nil
}

func Test_listener_Run(t *testing.T) {
	rs := miniredis.RunT(t)

	rc := redis.NewClient(&redis.Options{
		Addr: rs.Addr(),
	})
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		t.Errorf("failed to connect: %v", err)
	}
	h := &mockHandler{}

	l := NewEventListener(rc, "foo", slog.Default(), h)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = l.Run(ctx)
	}()

	// let time to connect and subscribe
	time.Sleep(100 * time.Millisecond)

	rs.Publish("foo", "hello world")
	time.Sleep(time.Second)
	cancel()

	want := []string{"hello world"}
	if !reflect.DeepEqual(h.payloads, want) {
		t.Errorf("listener::Run(), want %v, got %v", want, h.payloads)
	}
}
