# prometheus-huey-exporter

![build](https://github.com/mcosta74/prometheus-huey-exporter/actions/workflows/build.yml/badge.svg)
![GitHub](https://img.shields.io/github/license/mcosta74/prometheus-huey-exporter)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/mcosta74/prometheus-huey-exporter) [![Go Reference](https://pkg.go.dev/badge/github.com/mcosta74/prometheus-huey-exporter.svg)](https://pkg.go.dev/github.com/mcosta74/prometheus-huey-exporter)

Expose metrics from the [huey](https://huey.readthedocs.io/en/latest/) task queue


## Usage

### Huey configuration

Create a custom [signal](https://huey.readthedocs.io/en/latest/signals.html) handler that catch all the signals and publish on a specific Redis channel

```py
@djhuey.signal()
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

It's also possible to configure some values

- `-http.addr`: HTTP address 
- `-redis.addr`: address to the Redis instance

use

```sh
> ./prometheus-huey-exporter -h
```

for the full list
