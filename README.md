# prometheus-huey-exporter

![build](https://github.com/mcosta74/prometheus-huey-exporter/actions/workflows/build.yml/badge.svg)
![GitHub](https://img.shields.io/github/license/mcosta74/prometheus-huey-exporter)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/mcosta74/prometheus-huey-exporter) [![Go Reference](https://pkg.go.dev/badge/github.com/mcosta74/prometheus-huey-exporter.svg)](https://pkg.go.dev/github.com/mcosta74/prometheus-huey-exporter)

Expose metrics from the [huey](https://huey.readthedocs.io/en/latest/) task queue


## Usage

### Huey configuration

Create a custom [signal](https://huey.readthedocs.io/en/latest/signals.html) handler that catch all the signals and publish on a specific Redis channel

```py
@huey.signal()
def metrics(signal, task, exc=None):
    # conn = djhuey.HUEY.storage.conn
    conn = get_redis_connection()

    data = {
        'event': signal,
        'task_name': task.name,
        'task_id': task.id,
    }

    conn.publish(channel_name, json.dumps(data))
```


### Run the exporter

Start the exporter passing the same `channel_name` used in the python code

```sh
> ./prometheus-huey-exporter -redis.channel <channel_name>
```

### Configuration parameters

| Flag Name             | Environment Variable Name          | Description                                          |
| --------------------- | ---------------------------------- | ---------------------------------------------------- |
| `-log.level`          | `HUEY_EXPORTER_LOG_LEVEL`          | Log level                                            |
| `-log.format`         | `HUEY_EXPORTER_LOG_FORMAT`         | Log format                                           |
| `-redis.addr`         | `HUEY_EXPORTER_REDIS_ADDR`         | Address of the Redis instance to connect to          |
| `-redis.channel`      | `HUEY_EXPORTER_REDIS_CHANNEL`      | Redis channel to subscribe to listen for events      |
| `-metrics.namespace`  | `HUEY_EXPORTER_METRICS_NAMESPACE`  | Namespace for metrics                                |
| `-web.telemetry-path` | `HUEY_EXPORTER_WEB_PATH`           | Path under which to expose metrics                   |
| `-web.listen-addr`    | `HUEY_EXPORTER_WEB_LISTEN_ADDRESS` | Address to listen on for web interface and telemetry |

### Other Flags
| Flag Name  | Description               |
| ---------- | ------------------------- |
| `-version` | Show the version and quit |
| `-h`       | Show the help and quit    |