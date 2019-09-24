package main

import (
	"flag"
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
	r.Route("/admin", func(r chi.Router) {
		r.Get("/issue-book", issueBook)
		r.Get("/get-history/{name}", getHistory)
		r.Get("/complete-history", getCompleteHistory)
		r.Get("/return-book/{id}", returnBook)
	})
	r.Route("/user", func(r chi.Router) {
		r.Get("/check-availability/{id}", checkAvailability)
	})
	return r
}

func main() {
	flag.Parse()
	env = &envConfig.Env{}
	err := envconfig.Process("library", env)
	if err != nil {
		glog.Fatal(err)
	}
	dataStore = data_store.DbConnect(env)
	middleware.SetJwtSigningKey(env.JwtSigningKey)

	r := router()
	glog.Infof("Management-service binding on %s", ":"+env.ManagementSvcPort)
	err = http.ListenAndServe(":"+env.ManagementSvcPort, r)
	if err != nil {
		glog.Fatal(err)
	}
}
