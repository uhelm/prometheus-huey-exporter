package exporter

import (
	"context"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type EventListener interface {
	Run(ctx context.Context) error
}

type listener struct {
	rdb     *redis.Client
	channel string
	logger  *slog.Logger
	handler EventHandler
}

func NewEventListener(rdb *redis.Client, channel string, logger *slog.Logger, handler EventHandler) EventListener {
	return &listener{
		rdb:     rdb,
		channel: channel,
		logger:  logger,
		handler: handler,
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
			if err = l.handler.HandleEvent(evt); err != nil {
				l.logger.Error("running handleEvent()", "err", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
