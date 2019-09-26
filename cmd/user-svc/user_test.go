package main

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/library/data-store"
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
		db         *gorm.DB
		adminEmail string
		userEmail  string
		err        error
	)
	BeforeSuite(func() {
		env = &envConfig.Env{}
		err = envconfig.Process("library", env)
		Expect(err).To(BeNil())
		db, err = gorm.Open(env.SqlDialect, env.SqlUrl)
		Expect(err).To(BeNil())
		err = setupUserData(db)
		Expect(err).To(BeNil())
		dataStore = data_store.DbConnect(env)
		middleware.SetJwtSigningKey(env.JwtSigningKey)
	})
	Describe("Handlers Test", func() {
		Describe("Registration Test", func() {
			It("Should register a new user", func() {
				regReq := &models.Account{
					Name:     "testUser",
					Email:    "test@user.com",
					Password: "password",
				}
				marshalReq, err := json.Marshal(regReq)
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(marshalReq))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				register().ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				resp.Close = true
			})
		})

		Describe("Login Test", func() {
			Context("Admin Login", func() {
				It("Should provide valid JWT token", func() {
					adminEmail = "test@admin.com"
					loginReq := &models.LoginDetails{
						Email:    adminEmail,
						Password: "password",
					}
					marshalReq, err := json.Marshal(loginReq)
					Expect(err).To(BeNil())
					req := httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewBuffer(marshalReq))
					req.Header.Set("Content-Type", "application/json")
					rec := httptest.NewRecorder()
					login(models.AdminAccount).ServeHTTP(rec, req)
					resp := rec.Result()
					Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				})
			})
			Context("User Login", func() {
				It("Should provide valid JWT token", func() {
					userEmail = "test@user.com"
					loginReq := &models.LoginDetails{
						Email:    userEmail,
						Password: "password",
					}
					marshalReq, err := json.Marshal(loginReq)
					Expect(err).To(BeNil())
					req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(marshalReq))
					req.Header.Set("Content-Type", "application/json")
					rec := httptest.NewRecorder()
					login(models.UserAccount).ServeHTTP(rec, req)
					resp := rec.Result()
					Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				})
			})
			Context("Login with Invalid Credentials", func() {
				It("Should return Status Unauthorized", func() {
					loginReq := &models.LoginDetails{
						Email:    "invalid@user.com",
						Password: "invalidPassword",
					}
					marshalReq, err := json.Marshal(loginReq)
					Expect(err).To(BeNil())
					req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(marshalReq))
					req.Header.Set("Content-Type", "application/json")
					rec := httptest.NewRecorder()
					login(models.UserAccount).ServeHTTP(rec, req)
					resp := rec.Result()
					Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusUnauthorized))
				})
			})
		})
	})
	AfterSuite(func() {
		err = cleanTestData(db, adminEmail, userEmail)
		Expect(err).To(BeNil())
		err = db.Close()
		Expect(err).To(BeNil())
	})
})
