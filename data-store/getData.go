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

func (ds *DataStore) GetBookByID(id uint) (*models.Book, error) {
	books := &models.Book{}
	err := ds.Db.Where("id=?", id).Find(books).Error
	return books, err
}

func (ds *DataStore) GetBooksByName(name string) (*[]models.Book, error) {
	var books []models.Book
	query := `select * from book where name like '%` + name + `%'`
	err := ds.Db.Raw(query).Scan(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByAuthor(authorID uint) (*[]models.Book, error) {
	var books []models.Book
	var bookXauthor []models.BookXAuthor
	err := ds.Db.Where("author_id = ?", authorID).Find(&bookXauthor).Error
	if err != nil {
		return nil, err
	}
	for _, bXa := range bookXauthor {
		b := &models.Book{}
		err := ds.Db.Where("id = ?", bXa.BookID).Find(b).Error
		if err != nil {
			return nil, err
		}
		books = append(books, *b)
	}
	return &books, nil
}

func (ds *DataStore) GetBooksBySubject(subjectID uint) (*[]models.Book, error) {
	var books []models.Book
	var subjectXbook []models.SubjectXBook
	err := ds.Db.Where("subject_id = ?", subjectID).Find(&subjectXbook).Error
	if err != nil {
		return nil, err
	}
	for _, sXb := range subjectXbook {
		b := &models.Book{}
		err := ds.Db.Where("id = ?", sXb.BookID).Find(b).Error
		if err != nil {
			return nil, err
		}
		books = append(books, *b)
	}
	return &books, nil
}

func (ds *DataStore) GetAuthorsByName(name string) (*[]models.Author, error) {
	var authors []models.Author
	query := `select * from author where name like '%` + name + `%'`
	err := ds.Db.Raw(query).Scan(&authors).Error
	return &authors, err
}

func (ds *DataStore) GetAuthorByID(id uint) (*models.Author, error) {
	author := &models.Author{}
	err := ds.Db.Where("id = ?", id).Find(author).Error
	return author, err
}
