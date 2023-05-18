/*
The prometheus-nats-exporter exposes metrics about Huey2 tasks executions
*/
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/oklog/run"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"

	exporter "github.com/mcosta74/prometheus-huey-exporter"
)

var (
	version = "0.3.0"
	commit  = ""
	date    = ""
)

func main() {
	// Load Program Options
	opts, err := exporter.ParseOptions(os.Args[1:])
	if err != nil {
		fmt.Printf("error parsing options: %v\n", err)
		os.Exit(1)
	}

	if opts.PrintVersion {
		fmt.Println(version)
		os.Exit(1)
	}

	logger := getLogger(opts)

	logger.Info("service started")
	defer logger.Info("service stopped")

	// Establish connection with Redis
	rc, err := connectToRedis(opts.RedisAddress)
	if err != nil {
		logger.Error("connection with Redis failed", "addr", opts.RedisAddress, "err", err)
		os.Exit(1)
	}
	defer rc.Close()

	logger.Info("connected with Redis", "addr", opts.RedisAddress)

	var (
		metrics       = exporter.SetupMetrics(opts.MetricsPrefix)
		eventHandler  = exporter.NewEventHandler(metrics, logger)
		eventListener = exporter.NewEventListener(rc, opts.RedisChannel, logger.With("component", "event-listener"), eventHandler)
		httpHandler   = exporter.MakeHTTPHandler(opts.MetricsPath)
	)

	// Setup go-routines
	var g run.Group
	{
		// Signal Handler
		g.Add(run.SignalHandler(context.Background(), syscall.SIGTERM, syscall.SIGINT))
	}

	{
		// Prometheus HTTP
		webListener, err := net.Listen("tcp", opts.HTTPAddress)
		if err != nil {
			logger.Error("error listening on TCP", "addr", opts.HTTPAddress, "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			logger.Info("starting HTTP metrics server", "addr", fmt.Sprintf("http://%s%s", webListener.Addr(), opts.MetricsPath))
			return http.Serve(webListener, httpHandler)
		}, func(err error) {
			webListener.Close()
		})
	}

	{
		// Redis event listener
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return eventListener.Run(ctx)
		}, func(err error) {
			cancel()
		})
	}
	logger.Info("service shutdown", "err", g.Run())
}

func getLogger(opts exporter.Options) *slog.Logger {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(opts.LogLevel)); err != nil {
		lvl = slog.LevelInfo
	}
	hOpts := slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				// use only filename for source
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)

			case slog.TimeKey:
				// Use UTC time for timestamp
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
			return a
		},
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stderr, &hOpts)
	if opts.LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &hOpts)
	}
	return slog.New(handler)
}

func connectToRedis(addr string) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return rc, nil
}
