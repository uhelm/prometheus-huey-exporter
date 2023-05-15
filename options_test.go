package exporter_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	exporter "github.com/mcosta74/prometheus-huey-exporter"
)

func TestParseOptions(t *testing.T) {
	varsMap := map[string]string{
		"HUEY_EXPORTER_LOG_LEVEL":          "warn",
		"HUEY_EXPORTER_LOG_FORMAT":         "json",
		"HUEY_EXPORTER_REDIS_ADDR":         "redis.example.com:12345",
		"HUEY_EXPORTER_REDIS_CHANNEL":      "foo",
		"HUEY_EXPORTER_WEB_LISTEN_ADDRESS": ":8081",
		"HUEY_EXPORTER_METRICS_PATH":       "/hello",
		"HUEY_EXPORTER_METRICS_PREFIX":     "bar",
	}

	argsMap := map[string]string{
		"-log.level":      "warn",
		"-log.format":     "json",
		"-redis.addr":     "redis.example.com:12345",
		"-redis.channel":  "foo",
		"-http.addr":      ":8081",
		"-metrics.prefix": "/hello",
		"-metrics.path":   "bar",
	}

	t.Run("default values", func(t *testing.T) {
		got, err := exporter.ParseOptions([]string{})
		if err != nil {
			t.Errorf("ParseOptions() error=%v, want nil", err)
		}

		want := exporter.Options{
			LogLevel:      exporter.DefaultLogLevel.String(),
			LogFormat:     exporter.DefaultLogFormat,
			RedisAddress:  exporter.DefaultRedisAddress,
			RedisChannel:  exporter.DefaultRedisChannel,
			HTTPAddress:   exporter.DefaultHTTPAddress,
			MetricsPath:   exporter.DefaultMetricsPath,
			MetricsPrefix: "",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ParseOptions() = %+v, want %+v", got, want)
		}
	})

	t.Run("environment variable only", func(t *testing.T) {
		for k, v := range varsMap {
			os.Setenv(k, v)
		}

		got, err := exporter.ParseOptions([]string{})
		if err != nil {
			t.Errorf("ParseOptions() error=%v, want nil", err)
		}

		want := exporter.Options{
			LogLevel:      varsMap["HUEY_EXPORTER_LOG_LEVEL"],
			LogFormat:     varsMap["HUEY_EXPORTER_LOG_FORMAT"],
			RedisAddress:  varsMap["HUEY_EXPORTER_REDIS_ADDR"],
			RedisChannel:  varsMap["HUEY_EXPORTER_REDIS_CHANNEL"],
			HTTPAddress:   varsMap["HUEY_EXPORTER_WEB_LISTEN_ADDRESS"],
			MetricsPath:   varsMap["HUEY_EXPORTER_METRICS_PATH"],
			MetricsPrefix: varsMap["HUEY_EXPORTER_METRICS_PREFIX"],
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ParseOptions() = %+v, want %+v", got, want)
		}
		os.Clearenv()
	})

	t.Run("args only", func(t *testing.T) {
		args := make([]string, 0, len(argsMap))
		for k, v := range argsMap {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}

		got, err := exporter.ParseOptions(args)
		if err != nil {
			t.Errorf("ParseOptions() error=%v, want nil", err)
		}

		want := exporter.Options{
			LogLevel:      argsMap["-log.level"],
			LogFormat:     argsMap["-log.format"],
			RedisAddress:  argsMap["-redis.addr"],
			RedisChannel:  argsMap["-redis.channel"],
			HTTPAddress:   argsMap["-http.addr"],
			MetricsPath:   argsMap["-metrics.path"],
			MetricsPrefix: argsMap["-metrics.prefix"],
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ParseOptions() = %+v, want %+v", got, want)
		}

	})

	t.Run("args override env", func(t *testing.T) {
		os.Setenv("HUEY_EXPORTER_LOG_LEVEL", "info")
		args := []string{"-log.level=error"}

		got, err := exporter.ParseOptions(args)
		if err != nil {
			t.Errorf("ParseOptions() error=%v, want nil", err)
		}

		want := exporter.Options{
			LogLevel:      "error",
			LogFormat:     exporter.DefaultLogFormat,
			RedisAddress:  exporter.DefaultRedisAddress,
			RedisChannel:  exporter.DefaultRedisChannel,
			HTTPAddress:   exporter.DefaultHTTPAddress,
			MetricsPath:   exporter.DefaultMetricsPath,
			MetricsPrefix: "",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ParseOptions() = %+v, want %+v", got, want)
		}
	})
}
