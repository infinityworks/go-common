package metrics

import (
	"log"
	"net/http"
	"time"

	"os"

	"github.com/infinityworksltd/go-common/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// FunctionDurations - Create a summary to track elapsed time of our key functions
	FunctionDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "function_durations_seconds",
			Help:       "Function timings for Rancher Exporter",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"pkg", "fnc"})
)

// Init registers the prometheus metrics for the measurement of the exporter itsself.
func Init(config config.AppConfig) {
	prometheus.MustRegister(FunctionDurations)
	StartMetrics(config)
}

func RecordFunction(start time.Time, pkgName string, fncName string) {
	elapsed := float64(time.Since(start))
	FunctionDurations.WithLabelValues(pkgName, fncName).Observe(elapsed)
}

func StartMetrics(config config.AppConfig) {
	// Send metrics to Prometheus Handler
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{})
	http.Handle("/metrics", prometheus.InstrumentHandler("prometheus", handler))

	go func() {
		err := http.ListenAndServe(config.MetricsPort(), nil)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()
}
