package main

import (
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/models"
	"net/http"
)

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin", func(r chi.Router) {
		r.Post("/login", login(models.AdminAccount))
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/registration", register())
		r.Post("/login", login(models.UserAccount))
	})
	return r
}

func main() {
	env = &envConfig.Env{}
	err := envconfig.Process("library", env)
	if err != nil {
		glog.Fatal(err)
	}
	dataStore = data_store.DbConnect(env)

	r := router()
	err = http.ListenAndServe(":"+env.UserSvcPort, r)
	if err != nil {
		glog.Fatal(err)
	}
}
