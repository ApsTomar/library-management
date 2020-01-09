package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/library/models"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestSignupAndLogin(t *testing.T) {
	signupUrl := fmt.Sprintf("%s/register", testServer.URL)
	loginUrl := fmt.Sprintf("%s/login", testServer.URL)
	Convey("POST /register", t, func() {
		Convey("It should register a new user", func() {
			userEmail = "integration@user.com"
			regReq := &models.Account{
				Name:     "IntegrationUser",
				Email:    userEmail,
				Password: "password",
			}
			marshalReq, err := json.Marshal(regReq)
			So(err, ShouldBeNil)
			req, err := http.NewRequest(http.MethodPost, signupUrl, bytes.NewBuffer(marshalReq))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			So(err,ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
		})
	})

	Convey("POST /login", t, func() {
		Convey("Admin should login successfully", func() {
			adminEmail = "integration@admin.com"
			loginReq := &models.LoginDetails{
				Email:       adminEmail,
				Password:    "password",
				AccountRole: "admin",
			}
			marshalReq, err := json.Marshal(loginReq)
			So(err, ShouldBeNil)
			req, err := http.NewRequest(http.MethodPost, loginUrl, bytes.NewBuffer(marshalReq))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			defer resp.Body.Close()
		})
	})

	Convey("POST /login", t, func() {
		Convey("User should login successfully", func() {
			loginReq := &models.LoginDetails{
				Email:       userEmail,
				Password:    "password",
				AccountRole: "user",
			}
			marshalReq, err := json.Marshal(loginReq)
			So(err, ShouldBeNil)
			req, err := http.NewRequest(http.MethodPost, loginUrl, bytes.NewBuffer(marshalReq))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			defer resp.Body.Close()
		})
	})
}
