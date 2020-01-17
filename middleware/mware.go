package middleware

import (
	"github.com/go-chi/chi"
	"github.com/library/metrics"
)

func ChainMiddlewares(authMware bool, metrics *metrics.Metrics) chi.Middlewares {
	if authMware {
		return chi.Chain(MetricsCollector(metrics), CheckAuth(), AllowOptions(), RequestTracing())
	} else {
		return chi.Chain(MetricsCollector(metrics), AllowOptions(), RequestTracing())
	}
}
