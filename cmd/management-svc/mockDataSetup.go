package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/library/envConfig"
	"github.com/library/models"
	"time"
)

func setupAuthToken(env *envConfig.Env, db *gorm.DB) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   101010,
		"role": models.AdminAccount,
	})
	adminToken, err := token.SignedString([]byte(env.JwtSigningKey))
	if err != nil {
		return "", "", err
	}

	user := &models.Account{
		BaseModel:   *&models.BaseModel{ID: 101010},
		Name:        "testUser",
		Email:       "integration@user.com",
		AccountRole: "user",
	}
	err = db.Create(user).Error
	if err != nil {
		return "", "", err
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   101010,
		"role": models.UserAccount,
	})
	userToken, err := token.SignedString([]byte(env.JwtSigningKey))
	if err != nil {
		return "", "", err
	}

	return adminToken, userToken, err
}

func setupMockData(db *gorm.DB) error {
	author := models.Author{
		BaseModel:   *&models.BaseModel{ID: 101010},
		Name:        "intTestAuthor",
		DateOfBirth: "29 February 1600",
	}
	err := db.Create(&author).Error
	if err != nil {
		return err
	}

	subject := models.Subject{
		BaseModel: *&models.BaseModel{ID: 101010},
		Name:      "intTestSubject",
	}
	err = db.Create(&subject).Error
	if err != nil {
		return err
	}
	book := models.Book{
		BaseModel:     *&models.BaseModel{ID: 101010},
		Name:          "intTestBook",
		Subject:       "intTestSubject",
		AuthorID:      "101010",
		AuthorName:    "intTestAuthor",
		Available:     true,
		AvailableDate: time.Now(),
	}
	err = db.Create(&book).Error
	if err != nil {
		return err
	}
	return nil
}

func cleanMockData(db *gorm.DB) error {
	if err := db.Exec(`delete from account where id = ?`, "101010").Error; err != nil {
		return err
	}
	if err := db.Exec(`delete from author where id = ?`, "101010").Error; err != nil {
		return err
	}
	if err := db.Exec(`delete from subject where id = ?`, "101010").Error; err != nil {
		return err
	}
	if err := db.Exec(`delete from book where id = ?`, "101010").Error; err != nil {
		return err
	}
	return nil
}
