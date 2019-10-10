package main

import (
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Management-Service", func() {
	var (
		db         *gorm.DB
		r          *chi.Mux
		adminToken string
		userToken  string
		err        error
	)

	BeforeSuite(func() {
		env = &envConfig.Env{}
		err = envconfig.Process("library", env)
		Expect(err).To(BeNil())
		db, err = gorm.Open(env.SqlDialect, env.SqlUrl)
		Expect(err).To(BeNil())
		middleware.SetJwtSigningKey(env.JwtSigningKey)
		setupMockData()
		adminToken, userToken, err = setupAuthInfo(env)
		Expect(err).To(BeNil())
		dataStore = data_store.DbConnect(env)
		r = router()
	})
	Describe("Check availability", func() {
		It("Should return the availability of specified book", func() {
			req := httptest.NewRequest(http.MethodGet, "/check-availability/{id}", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "bearer "+adminToken)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			resp := rec.Result()
			Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
		})
	})
	Describe("Issue book", func() {
		It("Should issue book to specified user", func() {

		})
	})
	Describe("Handlers test", func() {
		Describe("Get Complete History", func() {
			It("Should return the complete book issue history", func() {
				req := httptest.NewRequest(http.MethodGet, "/complete-history", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "bearer "+adminToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))

			})
		})
	})
})
