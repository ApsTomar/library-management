package main

import (
	"github.com/kelseyhightower/envconfig"
	user_server "github.com/library/cmd/user-svc/user-server"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/library/models"
	password_hash "github.com/library/password-hash"
	"github.com/sirupsen/logrus"
	"net/http/httptest"
	"os"
	"testing"
)

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

	srv = user_server.NewServer(env, dataStore, nil)
	r := user_server.SetupRouter(srv)
	testServer = httptest.NewServer(r)
	_ = m.Run()
	if err := dataStore.ClearUserSvcData(adminEmail, userEmail); err != nil {
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
