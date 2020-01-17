package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestCounter    *prometheus.CounterVec
	LatencyCalculator *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "library",
			Subsystem: "user_svc",
			Name:      "requests_count",
			Help:      "Count of all user service requests",
		}, []string{"code", "handler"},
	)

	latencyCounter := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "library",
			Subsystem: "management_svc",
			Name:      "book_issue_latency",
			Help:      "Latency of book_issue from library",
			Buckets:   []float64{50, 100, 150, 200, 500},
		}, []string{"code", "handler"},
	)
	return &Metrics{
		RequestCounter:    requestCounter,
		LatencyCalculator: latencyCounter,
	}
}
