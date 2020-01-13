package user_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
)

func SetupRouter(srv *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.ChainMiddlewares(false)...)

	r.Post("/register", srv.register())
	r.Post("/login", srv.login())

	return r
}
