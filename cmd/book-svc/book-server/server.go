package book_server

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	datastore "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	DB        datastore.DbUtil
	Env       *envConfig.Env
	EfkLogger *fluent.Fluent
	TracingID string
	EfkTag    string
	TestRun   bool
}

func NewServer(env *envConfig.Env, db datastore.DbUtil, logger *fluent.Fluent) *Server {
	return &Server{
		DB:        db,
		Env:       env,
		EfkLogger: logger,
		TracingID: "",
		EfkTag:    "book_svc.logs",
		TestRun:   false,
	}
}

func (srv *Server) ListenAndServe(service string, port string) error {
	r := SetupRouter(srv)
	logrus.WithFields(logrus.Fields{
		"service": service,
	}).Info(service+" binding on ", ":"+port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
