package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestCounter    *prometheus.CounterVec
	LatencyCalculator *prometheus.HistogramVec
}

func NewMetrics(svc string) *Metrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "library",
			Subsystem: svc,
			Name:      "request_counter",
			Help:      "Count of all http requests in " + svc,
		}, []string{"code", "handler"},
	)

	latencyCounter := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "library",
			Subsystem: svc,
			Name:      "latency_calculator",
			Help:      "Latency of handlers in " + svc,
			Buckets:   []float64{10, 50, 100, 150, 200, 500},
		}, []string{"code", "handler"},
	)
	return &Metrics{
		RequestCounter:    requestCounter,
		LatencyCalculator: latencyCounter,
	}
}
