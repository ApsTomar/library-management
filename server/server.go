package server

import (
	"github.com/go-chi/chi"
	data_store "github.com/library/data-store"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	DB data_store.DbUtil
}

func NewServer(db data_store.DbUtil) *Server {
	return &Server{DB: db}
}

func (s *Server) ListenAndServe(r *chi.Mux, port string)error {
	logrus.WithFields(logrus.Fields{
		"service": "user-service",
	}).Info("user-service binding on ", ":"+port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
