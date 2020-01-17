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
			metricName := getMetricName(r.URL.Path)
			metrics.RequestCounter.WithLabelValues(strconv.Itoa(logResponseWriter.StatusCode), metricName).Add(float64(1))
			metrics.LatencyCalculator.WithLabelValues(strconv.Itoa(logResponseWriter.StatusCode), metricName).
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

func getMetricName(urlPath string) string {
	metricName := strings.TrimPrefix(urlPath, "/")
	if strings.Contains(metricName, "/") {
		parts := strings.Split(metricName, "/")
		if len(parts) > 1 {
			i := 0
			if parts[i] == "admin" || parts[i] == "user" {
				i++
			}
			metricName = parts[i]
			if parts[i] == "get" || parts[i] == "add" {
				metricName += "/" + parts[i+1]
			}
		}
	}
	return metricName
}
