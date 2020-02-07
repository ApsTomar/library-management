package middleware

import (
	"github.com/go-chi/chi"
	"github.com/library/envConfig"
	"github.com/library/metrics"
)

func ChainMiddlewares(authMware bool, metrics *metrics.Metrics, env *envConfig.Env) chi.Middlewares {
	if authMware {
		return chi.Chain(MetricsCollector(metrics, env), CheckAuth(), AllowOptions(), RequestTracing())
	} else {
		return chi.Chain(MetricsCollector(metrics, env), AllowOptions(), RequestTracing())
	}
}
