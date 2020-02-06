package book_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin/add", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics)...)
		r.Post("/author", srv.addAuthor)
		r.Post("/book", srv.addBook)
		r.Post("/subject", srv.addSubject)

	})
	r.Route("/get", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics)...)
		r.Get("/books", srv.getBooks)
		r.Get("/authors", srv.getAuthors)
		r.Get("/subjects", srv.getSubjects)
		r.Get("/books-by-name/{name}", srv.getBooksByName)
		r.Get("/book-by-id/{id}", srv.getBookByBookID)
		r.Get("/books-by-author/{id}", srv.getBooksByAuthorID)
		r.Get("/books-by-subject/{id}", srv.getBooksBySubjectID)
		r.Get("/author-by-name/{name}", srv.getAuthorByName)
		r.Get("/author-by-id/{id}", srv.getAuthorByID)
	})
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return r
}
