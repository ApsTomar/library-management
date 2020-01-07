package middleware

import (
	"github.com/go-chi/chi"
)

func ChainMiddlewares(authMware bool) chi.Middlewares {
	if authMware {
		return chi.Chain(CheckAuth(), AllowOptions, RequestTracing)
	} else {
		return chi.Chain(AllowOptions, RequestTracing)
	}
}
