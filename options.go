package exporter

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

const extraDoc = `
This application can be configured also using environment variables

  * HUEY_EXPORTER_LOG_LEVEL for the log level
  * HUEY_EXPORTER_LOG_FORMAT for the log format
  * HUEY_EXPORTER_REDIS_ADDR for the Redis instance address
  * HUEY_EXPORTER_REDIS_CHANNEL for the Redis channel to subscribe to
  * HUEY_EXPORTER_WEB_LISTEN_ADDRESS for the HTTP address to listen to
  * HUEY_EXPORTER_METRICS_PATH for the metrics endpoint
  * HUEY_EXPORTER_METRICS_PREFIX for the prefix to apply to all metrics

Command line options have priority over the env variables

`

// Default values for configuration options
const (
	DefaultLogLevel     = slog.LevelInfo
	DefaultLogFormat    = "text"
	DefaultRedisAddress = "localhost:6379"
	DefaultRedisChannel = "events"
	DefaultHTTPAddress  = ":9234"
	DefaultMetricsPath  = "/metrics"
	DefaultPrintVersion = false
)

// Option contains all the configuration parameters
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
		logLevelName      = "log.level"
		logFormatName     = "log.format"
		redisAddrName     = "redis.addr"
		redisChannelName  = "redis.channel"
		httpAddrName      = "web.listen"
		metricsPrefixName = "metrics.prefix"
		metricsPathName   = "metrics.path"
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
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s\n", fs.Name())
		fmt.Fprint(fs.Output(), extraDoc)
		fmt.Fprintln(fs.Output(), "Options:")
		fs.PrintDefaults()
	}

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

// ParseOptions reads configuration parameters from environment variables
// and then command line parameters args
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
