package main

import (
	"github.com/kelseyhightower/envconfig"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/library/server"
	"github.com/sirupsen/logrus"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	adminToken string
	userToken  string
	testServer *httptest.Server
	testAuthorID string
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
	adminToken, userToken, err = setupAuthInfo(env)

	srv = server.NewServer(dataStore)
	r := setupRouter()
	testServer = httptest.NewServer(r)
	_ = m.Run()
	if err := dataStore.ClearBookSvcData("testAuthor", "testSubject", "testBook"); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("cleaning testDB")
	}
}
