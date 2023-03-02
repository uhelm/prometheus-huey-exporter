package exporter

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Executions *prometheus.CounterVec
	Completed  *prometheus.CounterVec
	Locked     *prometheus.CounterVec
	Duration   *prometheus.HistogramVec
}

func SetupMetrics(prefix string) *Metrics {
	m := &Metrics{
		Executions: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "scheduler",
			Name:      "task_execution_total",
			Help:      "The Number of times a scheduler task has been executed.",
		}, []string{"task_name"}),
		Completed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "scheduler",
			Name:      "task_completed_total",
			Help:      "The Number of times a scheduler task has been completed.",
		}, []string{"task_name", "success"}),
		Locked: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "scheduler",
			Name:      "task_locked_total",
			Help:      "The Number of times a scheduler task failed to acquire a lock.",
		}, []string{"task_name"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "scheduler",
			Name:      "task_duration_seconds",
			Help:      "Task duration in seconds.",
		}, []string{"task_name", "success"}),
	}

	prometheus.MustRegister(m.Executions)
	prometheus.MustRegister(m.Completed)
	prometheus.MustRegister(m.Locked)
	prometheus.MustRegister(m.Duration)
	return m
}
