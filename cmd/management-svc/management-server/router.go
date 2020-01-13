package management_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
)

func SetupRouter(srv *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.ChainMiddlewares(true)...)

	r.Route("/admin", func(r chi.Router) {
		r.Post("/issue-book", srv.issueBook)
		r.Get("/get-history/{id}", srv.getHistory)
		r.Get("/complete-history", srv.getCompleteHistory)
		r.Get("/return-book/{id}", srv.returnBook)
	})
	r.Route("/user", func(r chi.Router) {
		r.Get("/check-availability/{id}", srv.checkAvailability)
	})
	return r
}
