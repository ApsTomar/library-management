package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	AdminAccount = "admin"
	UserAccount  = "user"
)

type Account struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccountRole  string `json:"accountRole"`
	Password     string `gorm:"-" json:"password"`
	PasswordHash string `json:"-"`
}

func (Account) TableName() string {
	return "account"
}

type Author struct {
	gorm.Model
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
}

func (Author) TableName() string {
	return "author"
}

type Book struct {
	gorm.Model
	Name          string    `json:"name"`
	Subject       string    `json:"subject"`
	AuthorID      string    `json:"authorId"`
	Available     bool      `json:"available"`
	AvailableDate time.Time `json:"availableDate"`
}

func (Book) TableName() string {
	return "book"
}

type BookXAuthor struct {
	BookID   uint `json:"bookId"`
	AuthorID uint `json:"authorId"`
}

func (BookXAuthor) TableName() string {
	return "book_x_author"
}

type Subject struct {
	gorm.Model
	Name string `json:"name"`
}

func (Subject) TableName() string {
	return "subject"
}

type SubjectXBook struct {
	SubjectID uint `json:"subjectId"`
	BookID    uint `json:"bookId"`
}

func (SubjectXBook) TableName() string {
	return "subject_x_book"
}

type BookHistory struct {
	UserID     uint       `json:"userId"`
	BookID     uint       `json:"bookId"`
	IssueDate  *time.Time `json:"issueDate"`
	ReturnDate *time.Time `json:"returnDate"`
	Returned   bool       `json:"returned"`
}

func (BookHistory) TableName() string {
	return "book_history"
}

type Response struct {
	AccountRole string `json:"accountRole"`
	Token       string `json:"token"`
}

type LoginDetails struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccountRole string `json:"accountRole"`
}

type AuthInfo struct {
	Role string
	jwt.StandardClaims
}
