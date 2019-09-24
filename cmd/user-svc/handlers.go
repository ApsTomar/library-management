package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/library/models"
	"github.com/library/password-hash"
	"net/http"
)

func register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := &models.Account{}
		err := json.NewDecoder(r.Body).Decode(account)
		if err != nil {
			glog.Errorf("error while decoding registration request body: %v", err)
			http.Error(w, "error while decoding registration request body", http.StatusInternalServerError)
			return
		}

		account.AccountRole = models.UserAccount
		hashedPwd, err := password_hash.HashPassword(account.Password)
		if err != nil {
			glog.Errorf("error creating password hash: %v", err)
			http.Error(w, "error creating password hash", http.StatusInternalServerError)
			return
		}
		account.PasswordHash = hashedPwd

		err = dataStore.CreateUserAccount(*account)
		if err != nil {
			glog.Errorf("error registering new user: %v", err)
			http.Error(w, "error registering new user", http.StatusInternalServerError)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   account.ID,
			"role": account.AccountRole,
		})
		tokenStr, err := token.SignedString([]byte(env.JwtSigningKey))
		if err != nil {
			glog.Errorf("error signing Jwt key: %v", err)
			http.Error(w, "error signing Jwt key", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: models.UserAccount, Token: tokenStr})
		if err != nil {
			glog.Errorf("error encoding json response: %v", err)
			http.Error(w, "error encoding json response", http.StatusInternalServerError)
		}
	}
}

func login(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details := &models.LoginDetails{}
		err := json.NewDecoder(r.Body).Decode(details)
		details.AccountRole = role
		if err != nil {
			glog.Errorf("error while decoding login request body: %v", err)
			http.Error(w, "error while decoding login request body", http.StatusInternalServerError)
			return
		}
		account, err := dataStore.VerifyUser(*details)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				glog.Errorf("error while logging in: %v", err)
				http.Error(w, "error while logging in", http.StatusUnauthorized)
			} else {
				glog.Errorf("error while logging in: %v", err)
				http.Error(w, "error while logging in", http.StatusInternalServerError)
			}
			return
		}

		ok := password_hash.ValidatePassword(details.Password, account.PasswordHash)
		if !ok{
			glog.Errorf("error while logging in: %v", err)
			http.Error(w, "error while logging in", http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   account.ID,
			"role": account.AccountRole,
		})
		tokenStr, err := token.SignedString([]byte(env.JwtSigningKey))
		if err != nil {
			glog.Errorf("error signing Jwt key: %v", err)
			http.Error(w, "error signing Jwt key", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: models.UserAccount, Token: tokenStr})
		if err != nil {
			glog.Errorf("error encoding json response: %v", err)
			http.Error(w, "error encoding json response", http.StatusInternalServerError)
		}
	}
}
