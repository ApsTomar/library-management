package main

import (
	"flag"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/kelseyhightower/envconfig"
	user_server "github.com/library/cmd/user-svc/user-server"
	"github.com/library/data-store"
	"github.com/library/efk"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	dataStore data_store.DbUtil
	env       *envConfig.Env
	logger    *fluent.Fluent
	srv       *user_server.Server
	testRun   bool
)

func init() {
	testRun = false
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
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
	dataStore = data_store.DbConnect(env, testRun)
	logger = efk.NewLogger(env)
	defer logger.Close()

	middleware.SetJwtSigningKey(env.JwtSigningKey)

	srv = user_server.NewServer(env, dataStore, logger)
	err = srv.ListenAndServe("user-service", env.UserSvcPort)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("user-server start")
	}
}
