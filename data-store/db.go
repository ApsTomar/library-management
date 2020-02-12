package data_store

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/library/envConfig"
	"github.com/library/migrations"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
	"time"
)

type DataStore struct {
	Db *gorm.DB
}

type DbUtil interface {
	InsertData
	GetData
	BookIssue
	VerifyUser(models.LoginDetails) (*models.Account, error)
}

type InsertData interface {
	CreateUserAccount(models.Account) error
	CreateAuthor(author models.Author) error
	CreateBook(book models.Book) error
	CreateSubject(book models.Subject) error
}

type GetData interface {
	GetSubjects() (*[]models.Subject, error)
	GetAuthors() (*[]models.Author, error)
	GetBooks() (*[]models.Book, error)
	GetBooksByName(string) (*[]models.Book, error)
	GetBookByID(uint) (*models.Book, error)
	GetBooksByAuthor(uint) (*[]models.Book, error)
	GetBooksBySubject(uint) (*[]models.Book, error)
	GetAuthorsByName(string) (*[]models.Author, error)
	GetAuthorByID(uint) (*models.Author, error)
}

type BookIssue interface {
	GetHistory(uint) (*[]models.BookHistory, error)
	GetCompleteHistory() (*[]models.BookHistory, error)
	CheckAvailability(uint) (bool, error)
	IssueBook(uint, uint) error
	ReturnBook(uint) error
}

var retryAttempts = 0

func DbConnect(dbConfig *envConfig.Env, testing bool) *DataStore {
	var sqlUrl string
	if testing {
		sqlUrl = dbConfig.TestSqlUrl
	} else {
		sqlUrl = dbConfig.SqlUrl
	}
	db, err := gorm.Open(dbConfig.SqlDialect, sqlUrl)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Info("DB connection not established, retrying ...")
		time.Sleep(time.Second * 5)
		retryAttempts++
		if retryAttempts > 5 {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("DB connection not established")
		}
		return DbConnect(dbConfig, testing)
	} else {
		err = migrations.InitMySQL(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error running migrations")
		}
		return &DataStore{Db: db}
	}
}