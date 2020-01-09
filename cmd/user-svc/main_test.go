package main

import (
	"github.com/kelseyhightower/envconfig"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/library/models"
	password_hash "github.com/library/password-hash"
	"github.com/library/server"
	"github.com/sirupsen/logrus"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

var (
	adminEmail string
	userEmail  string
	testServer *httptest.Server
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
	err = createAdminAccount()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("creating admin account")
	}
	middleware.SetJwtSigningKey(env.JwtSigningKey)

	srv = server.NewServer(dataStore)
	r := setupRouter()
	testServer = httptest.NewServer(r)
	_ = m.Run()
	if err := dataStore.ClearDb(adminEmail, userEmail); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("cleaning testDB")
	}
}

func createAdminAccount() error {
	password := "password"
	hashedPwd, err := password_hash.HashPassword(password)
	if err != nil {
		return err
	}
	admin := &models.Account{
		Name:         "IntegrationAdmin",
		Email:        "integration@admin.com",
		AccountRole:  models.AdminAccount,
		Password:     password,
		PasswordHash: hashedPwd,
	}
	return dataStore.CreateUserAccount(*admin)
}
