package exporter

import (
	"flag"
	"os"

	"golang.org/x/exp/slog"
)

const (
	DefaultLogLevel     = slog.LevelInfo
	DefaultLogFormat    = "text"
	DefaultRedisAddress = "localhost:6379"
	DefaultRedisChannel = "events"
	DefaultHTTPAddress  = ":9234"
	DefaultMetricsPath  = "/metrics"
	DefaultPrintVersion = false
)

type Options struct {
	LogLevel      string
	LogFormat     string
	RedisAddress  string
	RedisChannel  string
	HTTPAddress   string
	MetricsPath   string
	MetricsPrefix string
	PrintVersion  bool
}

func (o *Options) loadFromEnv() error {
	o.LogLevel = getEnvString("HUEY_EXPORTER_LOG_LEVEL", o.LogLevel)
	o.LogFormat = getEnvString("HUEY_EXPORTER_LOG_FORMAT", o.LogFormat)
	o.RedisAddress = getEnvString("HUEY_EXPORTER_REDIS_ADDR", o.RedisAddress)
	o.RedisChannel = getEnvString("HUEY_EXPORTER_REDIS_CHANNEL", o.RedisChannel)
	o.HTTPAddress = getEnvString("HUEY_EXPORTER_WEB_LISTEN_ADDRESS", o.HTTPAddress)
	o.MetricsPath = getEnvString("HUEY_EXPORTER_METRICS_PATH", o.MetricsPath)
	o.MetricsPrefix = getEnvString("HUEY_EXPORTER_METRICS_PREFIX", o.MetricsPrefix)
	return nil
}

func (o *Options) loadFromArgs(args []string) error {
	var (
		logLevelName      = "log-level"
		logFormatName     = "log-format"
		redisAddrName     = "redis-addr"
		redisChannelName  = "redis-channel"
		httpAddrName      = "http-addr"
		metricsPrefixName = "metrics-prefix"
		metricsPathName   = "metrics-path"
		printVersionName  = "version"
	)

	fs := flag.NewFlagSet("prometheus-huey-exporter", flag.ExitOnError)
	fs.String(logLevelName, o.LogLevel, "log level (debug, info, warn, error)")
	fs.String(logFormatName, o.LogFormat, "log format (text, json)")
	fs.String(redisAddrName, o.RedisAddress, "Redis address")
	fs.String(redisChannelName, o.RedisChannel, "Redis channel to listen for events")
	fs.String(httpAddrName, o.HTTPAddress, "HTTP address to listen")
	fs.String(metricsPrefixName, o.MetricsPrefix, "prefix to apply to all metrics")
	fs.String(metricsPathName, o.MetricsPath, "HTTP path for the metrics endpoint")
	fs.Bool(printVersionName, o.PrintVersion, "print the version")

	if err := fs.Parse(args); err != nil {
		return err
	}

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case logLevelName:
			o.LogLevel = f.Value.String()

		case logFormatName:
			o.LogFormat = f.Value.String()

		case redisAddrName:
			o.RedisAddress = f.Value.String()

		case redisChannelName:
			o.RedisChannel = f.Value.String()

		case httpAddrName:
			o.HTTPAddress = f.Value.String()

		case metricsPrefixName:
			o.MetricsPrefix = f.Value.String()

		case metricsPathName:
			o.MetricsPath = f.Value.String()

		case printVersionName:
			o.PrintVersion = f.Value.String() == "true"
		}
	})
	return nil
}

func ParseOptions(args []string) (Options, error) {
	opts := Options{
		LogLevel:     DefaultLogLevel.String(),
		LogFormat:    DefaultLogFormat,
		RedisAddress: DefaultRedisAddress,
		RedisChannel: DefaultRedisChannel,
		HTTPAddress:  DefaultHTTPAddress,
		MetricsPath:  DefaultMetricsPath,
	}

	if err := opts.loadFromEnv(); err != nil {
		return Options{}, err
	}

	if err := opts.loadFromArgs(args); err != nil {
		return Options{}, err
	}
	return opts, nil
}

func getEnvString(key, defaultVal string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return defaultVal
}
