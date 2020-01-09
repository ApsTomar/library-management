package server

import (
	"github.com/go-chi/chi"
	datastore "github.com/library/data-store"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	DB datastore.DbUtil
}

func NewServer(db datastore.DbUtil) *Server {
	return &Server{DB: db}
}

func (s *Server) ListenAndServe(r *chi.Mux, service string, port string) error {
	logrus.WithFields(logrus.Fields{
		"service": service,
	}).Info(service+" binding on ", ":"+port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
