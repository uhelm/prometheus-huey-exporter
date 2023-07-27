package exporter

import "github.com/prometheus/client_golang/prometheus"

// Metrics stores all the metrics exposed by the exporter
type Metrics struct {
	Executions   *prometheus.CounterVec
	Completed    *prometheus.CounterVec
	Locked       *prometheus.CounterVec
	Duration     *prometheus.HistogramVec
	LastDuration *prometheus.GaugeVec
}

// SetupMetrics takes care of initializing all the metrics (the names are prefixed
// with prefix) and register them to the default Registerer
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
		LastDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: prefix,
			Subsystem: "scheduler",
			Name:      "last_task_duration_seconds",
			Help:      "Last task duration in seconds.",
		}, []string{"task_name", "success"}),
	}

	prometheus.MustRegister(m.Executions)
	prometheus.MustRegister(m.Completed)
	prometheus.MustRegister(m.Locked)
	prometheus.MustRegister(m.Duration)
	prometheus.MustRegister(m.LastDuration)
	return m
}
