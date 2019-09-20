package data_store

import "github.com/library/models"

func (ds *DataStore) CreateUserAccount(acc models.Account) error {
	return ds.Db.Create(acc).Error
}

func (ds *DataStore) VerifyUser(details models.LoginDetails) (*models.Account, error) {
	account := &models.Account{}
	err := ds.Db.First(account).Where("email=? AND role=?", details.Email, details.AccountRole).Error
	return account, err
}

func (ds *DataStore) CreateAuthor(author models.Author) error {
	return ds.Db.Create(author).Error
}

func (ds *DataStore) CreateBook(book models.Book) error {
	return ds.Db.Create(book).Error
}

func (ds *DataStore) CreateSubject(subject models.Subject) error {
	return ds.Db.Create(subject).Error
}

