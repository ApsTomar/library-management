package data_store

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/library/envConfig"
	"github.com/library/migrations"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
)

type DataStore struct {
	Db *gorm.DB
}

type DbUtil interface {
	InsertData
	GetData
	BookIssue
	VerifyUser(models.LoginDetails) (*models.Account, error)
	ClearUserSvcData(string, string) error
	ClearBookSvcData(string, string, string) error
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

func DbConnect(dbConfig *envConfig.Env, testing bool) DbUtil {
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
		}).Fatal("DB connection not established")
	}
	err = migrations.InitMySQL(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error running migrations")
	}
	return &DataStore{Db: db}
}

func (ds *DataStore) ClearUserSvcData(adminEmail, userEmail string) error {
	if err := ds.Db.Exec(`delete from account where email = ?`, adminEmail).Error; err != nil {
		return err
	}
	if err := ds.Db.Exec(`delete from account where email = ?`, userEmail).Error; err != nil {
		return err
	}
	return nil
}

func (ds *DataStore) ClearBookSvcData(authorName, SubjectName, BookName string) error {
	if err := ds.Db.Exec(`delete from author where name = ?`, authorName).Error; err != nil {
		return err
	}
	if err := ds.Db.Exec(`delete from subject where name = ?`, SubjectName).Error; err != nil {
		return err
	}
	if err := ds.Db.Exec(`delete from book where name = ?`, BookName).Error; err != nil {
		return err
	}
	return nil
}
