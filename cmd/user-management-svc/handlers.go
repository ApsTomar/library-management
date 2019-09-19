package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
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
		}

		roleId, err := dataStore.GetAccountRoleID(models.UserAccount)
		if err != nil {
			glog.Errorf("error fetching account_role_id: %v", err)
			http.Error(w, "error fetching account_role_id", http.StatusInternalServerError)
		}
		account.AccountRoleId = roleId
		hashedPwd, err := password_hash.HashPassword(account.Password)
		if err != nil {
			glog.Errorf("error creating password hash: %v", err)
			http.Error(w, "error creating password hash", http.StatusInternalServerError)
		}
		account.PasswordHash = hashedPwd

		err = dataStore.CreateUserAccount(*account)
		if err != nil {
			glog.Errorf("error registering new user: %v", err)
			http.Error(w, "error registering new user", http.StatusInternalServerError)
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   account.ID,
			"role": account.AccountRoleId,
		})
		tokenStr,err := token.SignedString([]byte(env.JwtSigningKey))
		if err!= nil{
			glog.Errorf("error signing Jwt key: %v", err)
			http.Error(w, "error signing Jwt key", http.StatusInternalServerError)
		}
		err:= json.NewEncoder(w).Encode()
	}
}

func login(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
