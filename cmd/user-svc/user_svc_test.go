package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/cmd/user-svc/user-server"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/library/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("User-Service", func() {
	var (
		r          *chi.Mux
		adminEmail string
		userEmail  string
		err        error
	)
	BeforeSuite(func() {
		env = &envConfig.Env{}
		err = envconfig.Process("library", env)
		Expect(err).To(BeNil())
		testRun = true
		dataStore = data_store.DbConnect(env, true)
		srv = user_server.NewServer(env, dataStore, nil)
		srv.TestRun = true
		r = user_server.SetupRouter(srv)
		middleware.SetJwtSigningKey(srv.Env.JwtSigningKey)
		err = setupUserData(dataStore.Db)
		Expect(err).To(BeNil())
	})
	Describe("Handlers Test", func() {
		Describe("Registration Test", func() {
			It("Should register a new user", func() {
				userEmail = "unit@user.com"
				regReq := &models.Account{
					Name:     "testUser",
					Email:    "unit@user.com",
					Password: "password",
				}
				marshalReq, err := json.Marshal(regReq)
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalReq))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				resp.Body.Close()
			})
		})

		Describe("Login Test", func() {
			Context("Admin Login", func() {
				It("Should provide valid JWT token", func() {
					adminEmail = "unit@admin.com"
					loginReq := &models.LoginDetails{
						Email:       adminEmail,
						Password:    "password",
						AccountRole: "admin",
					}
					marshalReq, err := json.Marshal(loginReq)
					Expect(err).To(BeNil())
					req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalReq))
					req.Header.Set("Content-Type", "application/json")
					rec := httptest.NewRecorder()
					r.ServeHTTP(rec, req)
					resp := rec.Result()
					Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
					defer resp.Body.Close()
				})
			})
			Context("User Login", func() {
				It("Should provide valid JWT token", func() {
					loginReq := &models.LoginDetails{
						Email:       userEmail,
						Password:    "password",
						AccountRole: "user",
					}
					marshalReq, err := json.Marshal(loginReq)
					Expect(err).To(BeNil())
					req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalReq))
					req.Header.Set("Content-Type", "application/json")
					rec := httptest.NewRecorder()
					r.ServeHTTP(rec, req)
					resp := rec.Result()
					Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
					defer resp.Body.Close()
				})
			})
			Context("Login with Invalid Credentials", func() {
				It("Should return Status Unauthorized", func() {
					loginReq := &models.LoginDetails{
						Email:       userEmail,
						Password:    "invalidPassword",
						AccountRole: "user",
					}
					marshalReq, err := json.Marshal(loginReq)
					Expect(err).To(BeNil())
					req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalReq))
					req.Header.Set("Content-Type", "application/json")
					rec := httptest.NewRecorder()
					r.ServeHTTP(rec, req)
					resp := rec.Result()
					Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusUnauthorized))
					defer resp.Body.Close()
				})
			})
		})
	})
	AfterSuite(func() {
		err = cleanTestData(dataStore.Db, adminEmail, userEmail)
		Expect(err).To(BeNil())
	})
})
