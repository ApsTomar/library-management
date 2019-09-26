package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/library/envConfig"
	"github.com/library/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Book-Svc Handler Tests")
}

func setupBookData(env *envConfig.Env) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   9999999,
		"role": models.AdminAccount,
	})
	adminToken, err := token.SignedString([]byte(env.JwtSigningKey))
	if err != nil {
		return "", "", err
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   99999999,
		"role": models.UserAccount,
	})
	userToken, err := token.SignedString([]byte(env.JwtSigningKey))
	if err != nil {
		return "", "", err
	}
	return adminToken, userToken, err
}

func cleanTestData(db *gorm.DB, data *testData) error {
	if err := db.Exec(`delete from author where name = ?`, data.author).Error; err != nil {
		return err
	}
	if err := db.Exec(`delete from subject where name = ?`, data.subject).Error; err != nil {
		return err
	}
	if err := db.Exec(`delete from book where name = ?`, data.book).Error; err != nil {
		return err
	}
	return nil
}
