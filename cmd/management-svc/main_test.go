package main

import (
	"github.com/kelseyhightower/envconfig"
	management_server "github.com/library/cmd/management-svc/management-server"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/sirupsen/logrus"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	testServer *httptest.Server
	adminToken string
	userToken  string
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func TestMain(m *testing.M) {
	env = &envConfig.Env{}
	err := envconfig.Process("LIBRARY", env)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("processing env")
	}
	dataStore = data_store.DbConnect(env, true)
	middleware.SetJwtSigningKey(env.JwtSigningKey)
	if err := cleanMockData(dataStore.Db); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("cleaning mock data")
	}
	adminToken, userToken, err = setupAuthToken(env, dataStore.Db)
	err = setupMockData(dataStore.Db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("setting up mock data")
	}
	srv = management_server.NewServer(env, dataStore, nil)
	r := management_server.SetupRouter(srv)
	testServer = httptest.NewServer(r)
	_ = m.Run()
	if err := cleanMockData(dataStore.Db); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("cleaning mock data")
	}
}
