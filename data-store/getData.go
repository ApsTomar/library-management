package data_store

import (
	"github.com/library/models"
)

func (ds *DataStore) GetSubjects() (*[]models.Subject, error) {
	var subjects []models.Subject
	err := ds.Db.Find(&subjects).Error
	return &subjects, err
}

func (ds *DataStore) GetAuthors() (*[]models.Author, error) {
	var authors []models.Author
	err := ds.Db.Find(&authors).Error
	return &authors, err
}

func (ds *DataStore) GetBooks() (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByID(id uint) (*models.Book, error) {
	books := &models.Book{}
	err := ds.Db.Find(books).Where("id=?", id).Error
	return books, err
}

func (ds *DataStore) GetBooksByAuthor(authorID uint) (*[]models.Book, error) {
	var books []models.Book
	query := `select * from books where id = (select book_id from book_x_author where author_id = ?)`
	if err := ds.Db.Raw(query, authorID).Scan(&books).Error; err != nil {
		return nil, err
	}
	return &books, nil
}

func (ds *DataStore) GetBooksBySubject(subjectID uint) (*[]models.Book, error) {
	var books []models.Book
	query := `select * from books where id = (select book_id from subject_x_book where subject_id = ?)`
	if err := ds.Db.Raw(query, subjectID).Scan(&books).Error; err != nil {
		return nil, err
	}
	return &books, nil
}
