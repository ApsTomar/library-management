package main

import (
	"flag"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/data-store"
	"github.com/library/efk"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const efkTag = "user_svc.logs"

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
	logger    *fluent.Fluent
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AllowOptions)

	r.Post("/register", register())
	r.Post("/login", login())

	return r
}

func main() {
	flag.Parse()
	env = &envConfig.Env{}
	err := envconfig.Process("LIBRARY", env)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("processing env")
	}
	dataStore = data_store.DbConnect(env)
	logger = efk.NewLogger(env)
	defer logger.Close()

	middleware.SetJwtSigningKey(env.JwtSigningKey)

	r := router()
	logrus.WithFields(logrus.Fields{
		"service": "user-service",
	}).Info("user-service binding on ", ":"+env.UserSvcPort)

	err = http.ListenAndServe(":"+env.UserSvcPort, r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("server start")
	}
}
