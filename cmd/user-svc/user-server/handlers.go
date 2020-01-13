package user_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/library/efk"
	"github.com/library/middleware"
	"github.com/library/models"
	"github.com/library/password-hash"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (srv *Server) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		account := &models.Account{}
		err := json.NewDecoder(r.Body).Decode(account)
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		account.AccountRole = models.UserAccount
		hashedPwd, err := password_hash.HashPassword(account.Password)
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		account.PasswordHash = hashedPwd
		err = srv.DB.CreateUserAccount(*account)
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				handleError(w, ctx, srv, "registration", err, http.StatusBadRequest)
				return
			}
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		// get the created user account
		acc, err := srv.DB.VerifyUser(*&models.LoginDetails{
			Email:       account.Email,
			Password:    account.Password,
			AccountRole: account.AccountRole,
		})
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   acc.ID,
			"role": acc.AccountRole,
		})

		tokenStr, err := token.SignedString([]byte(srv.Env.JwtSigningKey))
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"statusCode": http.StatusOK,
		}).Info(fmt.Sprintf("new user registered with email: %v", account.Email))

		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: models.UserAccount, Token: tokenStr})
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
	}
}

func (srv *Server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		details := &models.LoginDetails{}
		err := json.NewDecoder(r.Body).Decode(details)
		if err != nil {
			handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			return
		}
		account, err := srv.DB.VerifyUser(*details)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				handleError(w, ctx, srv, "login", errors.New(fmt.Sprintf("no such %v found", details.AccountRole)), http.StatusBadRequest)
			} else {
				handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			}
			return
		}

		ok := password_hash.ValidatePassword(details.Password, account.PasswordHash)
		if !ok {
			handleError(w, ctx, srv, "login", errors.New("invalid password"), http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   account.ID,
			"role": account.AccountRole,
		})
		tokenStr, err := token.SignedString([]byte(srv.Env.JwtSigningKey))
		if err != nil {
			handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"statusCode": http.StatusOK,
		}).Info(fmt.Sprintf("user login with email: %v", account.Email))
		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: details.AccountRole, Token: tokenStr})
		if err != nil {
			handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			return
		}
	}
}

func handleError(w http.ResponseWriter, ctx context.Context, srv *Server, task string, err error, statusCode int) {
	if !srv.TestRun {
		srv.TracingID = ctx.Value(middleware.RequestTracingID).(string)
		efk.LogError(srv.EfkLogger, srv.EfkTag, srv.TracingID, task, err, statusCode)
	}
	http.Error(w, err.Error(), statusCode)

	logrus.WithFields(logrus.Fields{
		"tracingID":  srv.TracingID,
		"statusCode": statusCode,
		"error":      err,
	}).Error(task)
}
