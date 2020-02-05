package management_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server, prom *prometheus.Registry) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics)...)
		r.Post("/issue-book", srv.issueBook)
		r.Get("/get-history/{id}", srv.getHistory)
		r.Get("/complete-history", srv.getCompleteHistory)
		r.Get("/return-book/{id}", srv.returnBook)
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics)...)
		r.Get("/check-availability/{id}", srv.checkAvailability)
	})
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return r
}
