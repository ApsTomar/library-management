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
	"github.com/library/models"
	"net/http"
)

const efkTag = "user_svc.logs"

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
	logger    *fluent.Fluent
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin", func(router chi.Router) {
		router.Post("/login", login(models.AdminAccount))
	})

	r.Route("/user", func(router chi.Router) {
		router.Post("/register", register())
		router.Post("/login", login(models.UserAccount))
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
	logger = efk.NewLogger(env)
	defer logger.Close()

	dataStore = data_store.DbConnect(env)
	middleware.SetJwtSigningKey(env.JwtSigningKey)

	r := router()
	glog.Infof("User-service binding on %s", ":"+env.UserSvcPort)
	err = http.ListenAndServe(":"+env.UserSvcPort, r)
	if err != nil {
		glog.Fatal(err)
	}
}
