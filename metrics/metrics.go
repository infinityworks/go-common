package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// FunctionDurations - Create a summary to track elapsed time of our key functions
	duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Request duration seconds for HTTP Request",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 1, 2},
	}, []string{"method", "name"})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Total number of HTTP requests",
		},
		[]string{"status", "name"},
	)
)

// Init registers the prometheus metrics for the measurement of the exporter itsself.
func Init() {
	prometheus.MustRegister(duration)
	prometheus.MustRegister(counter)
}

func Instrument(time float64, statusCode int, method string, name string) {
	l := prometheus.Labels{
		"status": fmt.Sprint(statusCode),
		"name":   name,
	}

	duration.WithLabelValues(method, name).Observe(time)
	counter.With(l).Inc()
}
