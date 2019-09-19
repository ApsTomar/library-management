package data_store

import (
	"github.com/jinzhu/gorm"
	"github.com/library/envConfig"
	"github.com/library/migrations"
	"github.com/library/models"
	"log"
)

type DataStore struct {
	Db *gorm.DB
}

type DbUtil interface {
	GetAccountRoleID(string) (uint, error)
	CreateUserAccount(models.Account) error
}

func DbConnect(dbConfig *envConfig.Env) DbUtil {
	db, err := gorm.Open(dbConfig.SqlDialect, dbConfig.SqlUrl)
	if err != nil {
		log.Fatalf("DB connection not established due to: %v", err)
	}
	err = migrations.InitMySQL(db)
	if err != nil {
		log.Fatalf("error running migrations: %v", err)
	}
	return &DataStore{Db: db}
}

func (ds *DataStore) GetAccountRoleID(role string) (uint, error) {
	accType := &models.AccountType{}
	err := ds.Db.First(accType).Where("role=?", role).Error
	if err != nil {
		return -1, nil
	}
	return accType.ID, nil
}

func (ds *DataStore) CreateUserAccount(acc models.Account) error {
	return ds.Db.Create(acc).Error
}
