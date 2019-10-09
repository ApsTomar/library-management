package data_store

import (
	"errors"
	"github.com/library/models"
	"time"
)

func (ds *DataStore) GetCompleteHistory() (*[]models.BookHistory, error) {
	var history []models.BookHistory
	err := ds.Db.Find(&history).Error
	return &history, err
}

func (ds *DataStore) GetHistory(id uint) (*[]models.BookHistory, error) {
	var history []models.BookHistory
	query := `select * from book_history where book_id = ?`
	err := ds.Db.Raw(query, id).Scan(&history).Error
	return &history, err
}

func (ds *DataStore) CheckAvailability(id uint) (bool, error) {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", id).First(book).Error
	if err != nil {
		return false, err
	}
	return book.Available, nil
}

func (ds *DataStore) IssueBook(bookID, userID uint) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	if book.Available == false {
		return errors.New("book unavailable")
	}
	currentTime := time.Now()
	returnDate := currentTime.AddDate(0, 0, 15)
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"available":      false,
		"available_date": returnDate,
	}).Error
	if err != nil {
		return err
	}
	history := &models.BookHistory{
		UserID:     userID,
		BookID:     bookID,
		IssueDate:  &currentTime,
		ReturnDate: &returnDate,
		Returned:   false,
	}
	return ds.Db.Create(history).Error
}

func (ds *DataStore) ReturnBook(bookID uint) error {
	book := &models.Book{}
	err := ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"available":      true,
		"available_date": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	history := &models.BookHistory{}
	return ds.Db.Model(history).Where("book_id = ?", bookID).Updates(map[string]interface{}{
		"return_date": time.Now(),
		"returned":    true,
	}).Error
}
