package book_server

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	datastore "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	prom        *prometheus.Registry
	promMetrics *metrics.Metrics
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
	prom = prometheus.NewRegistry()
	promMetrics = metrics.NewMetrics("book_svc")
	prom.MustRegister(promMetrics.RequestCounter)
	prom.MustRegister(promMetrics.LatencyCalculator)

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
