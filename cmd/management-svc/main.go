package main

import (
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"net/http"
)

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware()...)

	r.Route("/admin/add", func(r chi.Router) {
		r.Post("/author", addAuthor)
		r.Post("/book", addBook)
		r.Post("/subject", addSubject)

	})
	r.Route("/get", func(r chi.Router) {
		r.Post("books", getBooks)
		r.Post("/authors", getAuthors)
		r.Post("/subjects", getSubjects)
		r.Post("/book-by-id/{id}", getBookByBookID)
		r.Post("/books-by-author/{id}", getBooksByAuthorID)
		r.Post("books-by-subject/{id}", getBooksBySubjectID)
	})
	return r
}

func main() {
	env = &envConfig.Env{}
	err := envconfig.Process("library", env)
	if err != nil {
		glog.Fatal(err)
	}
	middleware.SetJwtSigningKey(env.JwtSigningKey)
	dataStore = data_store.DbConnect(env)

	r := router()
	err = http.ListenAndServe(":"+env.BookSvcPort, r)
	if err != nil {
		glog.Fatal(err)
	}
}
