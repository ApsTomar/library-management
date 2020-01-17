package middleware

import (
	"github.com/library/metrics"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func MetricsCollector(metrics *metrics.Metrics) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			logResponseWriter := NewLogResponseWriter(w)
			handler.ServeHTTP(logResponseWriter, r)
			handler := r.URL.Path
			handler = strings.TrimPrefix(handler, "/")
			metrics.RequestCounter.WithLabelValues(strconv.Itoa(logResponseWriter.StatusCode), handler).Add(float64(1))
			metrics.LatencyCalculator.WithLabelValues(strconv.Itoa(logResponseWriter.StatusCode), handler).
				Observe(float64(time.Since(startTime).Milliseconds()))
		})
	}
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{w, http.StatusOK}
}

func (rw *LogResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
