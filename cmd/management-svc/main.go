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
	"github.com/library/server"
	"github.com/sirupsen/logrus"
	"os"
)

const efkTag = "management_svc.logs"

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
	logger    *fluent.Fluent
	srv       *server.Server
	tracingID string
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.ChainMiddlewares(true)...)

	r.Route("/admin", func(r chi.Router) {
		r.Get("/issue-book", issueBook)
		r.Get("/get-history/{id}", getHistory)
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
	err := envconfig.Process("LIBRARY", env)
	if err != nil {
		glog.Fatal(err)
	}
	logger = efk.NewLogger(env)
	defer logger.Close()

	dataStore = data_store.DbConnect(env)
	middleware.SetJwtSigningKey(env.JwtSigningKey)

	srv = server.NewServer(dataStore)
	r := router()
	err = srv.ListenAndServe(r, env.ManagementSvcPort)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("server start")
	}
}
