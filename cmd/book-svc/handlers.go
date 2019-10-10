package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/library/efk"
	"github.com/library/middleware"
	"github.com/library/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAuthInfoFromContext(ctx context.Context) *models.AuthInfo {
	return ctx.Value(middleware.ContextAuthInfo).(*models.AuthInfo)
}

func addAuthor(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		handleError(w, "add_author", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	author := &models.Author{}
	err := json.NewDecoder(r.Body).Decode(author)
	if err != nil {
		handleError(w, "add_author", err, http.StatusInternalServerError)
		return
	}
	err = dataStore.CreateAuthor(*author)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, "add_author", errors.New("duplicate entry"), http.StatusBadRequest)
			return
		}
		handleError(w, "add_author", err, http.StatusInternalServerError)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		handleError(w, "add_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		handleError(w, "add_book", err, http.StatusInternalServerError)
		return
	}
	book.AvailableDate = time.Now()
	book.Available = true
	err = dataStore.CreateBook(*book)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, "add_book", err, http.StatusBadRequest)
			return
		}
		handleError(w, "add_book", err, http.StatusInternalServerError)
	}
}

func addSubject(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		handleError(w, "add_subject", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	subject := &models.Subject{}
	err := json.NewDecoder(r.Body).Decode(subject)
	if err != nil {
		handleError(w, "add_subject", err, http.StatusInternalServerError)
		return
	}
	err = dataStore.CreateSubject(*subject)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, "add_subject", errors.New("duplicate entry"), http.StatusBadRequest)
			return
		}
		handleError(w, "add_subject", err, http.StatusInternalServerError)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := dataStore.GetBooks()
	if err != nil {
		if err == gorm.ErrRecordNotFound || books == nil {
			handleError(w, "get_books", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_books", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, "get_books", err, http.StatusInternalServerError)
	}
}

func getBooksByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	books, err := dataStore.GetBooksByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, "get_books_by_name", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_books_by_name", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, "get_books_by_name", err, http.StatusInternalServerError)
	}
}

func getBookByBookID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, "get_book_by_id", err, http.StatusInternalServerError)
		return
	}
	book, err := dataStore.GetBookByID(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, "get_book_by_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_book_by_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		handleError(w, "get_book_by_id", err, http.StatusInternalServerError)
	}
}

func getBooksByAuthorID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, "get_books_by_author_id", err, http.StatusInternalServerError)
		return
	}
	books, err := dataStore.GetBooksByAuthor(uint(authorID))
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, "get_books_by_author_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_books_by_author_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, "get_books_by_author_id", err, http.StatusInternalServerError)
	}
}

func getBooksBySubjectID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subjectID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, "get_books_by_subject_id", err, http.StatusInternalServerError)
		return
	}
	books, err := dataStore.GetBooksBySubject(uint(subjectID))
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, "get_books_by_subject_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_books_by_subject_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, "get_books_by_subject_id", err, http.StatusInternalServerError)
	}
}

func getSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := dataStore.GetSubjects()
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*subjects) == 0 {
			handleError(w, "get_subjects", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, "get_subjects", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(subjects)
	if err != nil {
		handleError(w, "get_subjects", err, http.StatusInternalServerError)
	}
}

func getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := dataStore.GetAuthors()
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*authors) == 0 {
			handleError(w, "get_authors", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, "get_authors", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		handleError(w, "get_authors", err, http.StatusInternalServerError)
	}
}

func getAuthorByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	authors, err := dataStore.GetAuthorsByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*authors) == 0 {
			handleError(w, "get_author_by_name", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_author_by_name", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		handleError(w, "get_author_by_name", err, http.StatusInternalServerError)
	}
}

func getAuthorByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, "get_author_by_id", err, http.StatusInternalServerError)
		return
	}
	author, err := dataStore.GetAuthorByID(uint(authorID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, "get_author_by_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, "get_author_by_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(author)
	if err != nil {
		handleError(w, "get_author_by_id", err, http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, task string, err error, statusCode int) {
	efk.LogError(logger, efkTag, task, err, statusCode)
	http.Error(w, err.Error(), statusCode)
	glog.Error(err)
}