package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var jwtSigningKey string

const ContextAuthInfo = "auth-info"

func SetJwtSigningKey(key string) {
	jwtSigningKey = key
}

func CheckAuth() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acc := &models.AuthInfo{}
			var token string
			t, ok := r.Header["Authorization"]
			if ok && len(t) >= 1 {
				token = t[0]
				token = strings.TrimPrefix(token, "Bearer ")
			}
			if token == "" {
				logrus.WithFields(logrus.Fields{
					"statusCode": http.StatusUnauthorized,
				}).Error("empty token")
				http.Error(w, "empty token", http.StatusUnauthorized)
				return
			}
			if err := ValidateToken(acc, token, jwtSigningKey); err != nil {
				logrus.WithFields(logrus.Fields{
					"statusCode": http.StatusUnauthorized,
					"error":      err,
				}).Error("invalid token")
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextAuthInfo, acc)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateToken(claims jwt.Claims, token, jwtSigningKey string) error {
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("unexpected signing Method: %v", token.Header["alg"])
			return nil, msg
		}
		return []byte(jwtSigningKey), nil
	})
	if err != nil {
		return err
	}
	if parsedToken == nil || !parsedToken.Valid {
		return err
	}
	return nil
}
