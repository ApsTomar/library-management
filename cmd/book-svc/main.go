package main

import (
	"flag"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/data-store"
	"github.com/library/efk"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

const efkTag = "book_svc.logs"

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
	logger    *fluent.Fluent
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AllowOptions)
	r.Use(middleware.AuthMiddleware()...)

	r.Route("/admin/add", func(r chi.Router) {
		r.Post("/author", addAuthor)
		r.Post("/book", addBook)
		r.Post("/subject", addSubject)

	})
	r.Route("/get", func(r chi.Router) {
		r.Get("/books", getBooks)
		r.Get("/authors", getAuthors)
		r.Get("/subjects", getSubjects)
		r.Get("/books-by-name/{name}", getBooksByName)
		r.Get("/book-by-id/{id}", getBookByBookID)
		r.Get("/books-by-author/{id}", getBooksByAuthorID)
		r.Get("/books-by-subject/{id}", getBooksBySubjectID)
		r.Get("/author-by-name/{name}", getAuthorByName)
		r.Get("/author-by-id/{id}", getAuthorByID)
	})
	return r
}

func main() {
	flag.Parse()
	env = &envConfig.Env{}
	err := envconfig.Process("LIBRARY", env)
	if err != nil {
		glog.Fatal(err)
	}
	logger = efk.NewLogger(env)
	defer logger.Close()

	middleware.SetJwtSigningKey(env.JwtSigningKey)
	dataStore = data_store.DbConnect(env)

	r := router()
	logrus.WithFields(logrus.Fields{
		"service": "book-service",
	}).Info("book-service binding on ", ":"+env.BookSvcPort)

	err = http.ListenAndServe(":"+env.BookSvcPort, r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("server start")
	}
}
