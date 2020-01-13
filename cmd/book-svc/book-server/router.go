package book_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
)

func SetupRouter(srv *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AllowOptions)
	r.Use(middleware.RequestTracing)
	r.Use(middleware.ChainMiddlewares(true)...)

	r.Route("/admin/add", func(r chi.Router) {
		r.Post("/author", srv.addAuthor)
		r.Post("/book", srv.addBook)
		r.Post("/subject", srv.addSubject)

	})
	r.Route("/get", func(r chi.Router) {
		r.Get("/books", srv.getBooks)
		r.Get("/authors", srv.getAuthors)
		r.Get("/subjects", srv.getSubjects)
		r.Get("/books-by-name/{name}",srv.getBooksByName)
		r.Get("/book-by-id/{id}", srv.getBookByBookID)
		r.Get("/books-by-author/{id}", srv.getBooksByAuthorID)
		r.Get("/books-by-subject/{id}", srv.getBooksBySubjectID)
		r.Get("/author-by-name/{name}", srv.getAuthorByName)
		r.Get("/author-by-id/{id}", srv.getAuthorByID)
	})
	return r
}

