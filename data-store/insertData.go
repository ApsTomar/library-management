package data_store

import (
	"github.com/library/models"
	"strconv"
	"strings"
)

func (ds *DataStore) CreateUserAccount(acc models.Account) error {
	return ds.Db.Create(&acc).Error
}

func (ds *DataStore) VerifyUser(details models.LoginDetails) (*models.Account, error) {
	account := &models.Account{}
	err := ds.Db.Where("email=? AND account_role=?", details.Email, details.AccountRole).First(account).Error
	return account, err
}

func (ds *DataStore) CreateAuthor(author models.Author) error {
	return ds.Db.Create(&author).Error
}

func (ds *DataStore) CreateBook(book models.Book) error {
	authors := strings.Split(book.AuthorID, ",")
	err := ds.Db.Create(&book).Error
	if err != nil {
		return err
	}
	for _, author := range authors {
		authorId, err := strconv.Atoi(author)
		if err != nil {
			return err
		}
		bookXAuthor := &models.BookXAuthor{
			BookID:   book.ID,
			AuthorID: uint(authorId),
		}
		if err := ds.Db.Create(bookXAuthor).Error; err != nil {
			return err
		}
	}
	subject := &models.Subject{}
	if err := ds.Db.Where("name = ?", book.Subject).First(subject).Error; err != nil {
		return err
	}
	subjectXBook := &models.SubjectXBook{
		SubjectID: subject.ID,
		BookID:    book.ID,
	}
	if err := ds.Db.Create(subjectXBook).Error; err != nil {
		return err
	}
	return err
}

func (ds *DataStore) CreateSubject(subject models.Subject) error {
	return ds.Db.Create(&subject).Error
}
