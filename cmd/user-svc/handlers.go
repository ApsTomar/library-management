package main

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/library/efk"
	"github.com/library/models"
	"github.com/library/password-hash"
	"net/http"
)

func register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := &models.Account{}
		err := json.NewDecoder(r.Body).Decode(account)
		if err != nil {
			handleError(w, "registration", err, http.StatusInternalServerError)
			return
		}
		account.AccountRole = models.UserAccount
		hashedPwd, err := password_hash.HashPassword(account.Password)
		if err != nil {
			handleError(w, "registration", err, http.StatusInternalServerError)
			return
		}
		account.PasswordHash = hashedPwd
		err = dataStore.CreateUserAccount(*account)
		if err != nil {
			handleError(w, "registration", err, http.StatusInternalServerError)
			return
		}

		acc, err := dataStore.VerifyUser(*&models.LoginDetails{
			Email:       account.Email,
			Password:    account.Password,
			AccountRole: account.AccountRole,
		})
		if err != nil {
			handleError(w, "registration", err, http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   acc.ID,
			"role": acc.AccountRole,
		})

		tokenStr, err := token.SignedString([]byte(env.JwtSigningKey))
		if err != nil {
			handleError(w, "registration", err, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: models.UserAccount, Token: tokenStr})
		if err != nil {
			handleError(w, "registration", err, http.StatusInternalServerError)
			return
		}
	}
}

func login(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details := &models.LoginDetails{}
		err := json.NewDecoder(r.Body).Decode(details)
		details.AccountRole = role
		if err != nil {
			handleError(w, "login", err, http.StatusInternalServerError)
			return
		}
		account, err := dataStore.VerifyUser(*details)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				handleError(w, "login", errors.New("no user found"), http.StatusInternalServerError)
			} else {
				handleError(w, "login", err, http.StatusInternalServerError)
			}
			return
		}

		ok := password_hash.ValidatePassword(details.Password, account.PasswordHash)
		if !ok {
			handleError(w, "login", errors.New("invalid password"), http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   account.ID,
			"role": account.AccountRole,
		})
		tokenStr, err := token.SignedString([]byte(env.JwtSigningKey))
		if err != nil {
			handleError(w, "login", err, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: models.UserAccount, Token: tokenStr})
		if err != nil {
			handleError(w, "login", err, http.StatusInternalServerError)
			return
		}
	}
}

func handleError(w http.ResponseWriter, task string, err error, statusCode int) {
	efk.LogError(logger, efkTag, task, err, statusCode)
	http.Error(w, err.Error(), statusCode)
	glog.Error(err)
}
