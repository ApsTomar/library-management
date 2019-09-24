package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
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
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to add authors", http.StatusUnauthorized)
		return
	}
	author := &models.Author{}
	err := json.NewDecoder(r.Body).Decode(author)
	if err != nil {
		glog.Errorf("error while decoding request body: %v", err)
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
		return
	}
	err = dataStore.CreateAuthor(*author)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			glog.Errorf("duplicate entry: %v", err)
			http.Error(w, "duplicate entry", http.StatusBadRequest)
		}
		glog.Errorf("error creating new author: %v", err)
		http.Error(w, "error creating new author", http.StatusInternalServerError)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to add books", http.StatusUnauthorized)
		return
	}
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		glog.Errorf("error while decoding request body: %v", err)
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
		return
	}
	book.AvailableDate = time.Now()
	book.Available = true
	err = dataStore.CreateBook(*book)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			glog.Errorf("duplicate entry: %v", err)
			http.Error(w, "duplicate entry", http.StatusBadRequest)
		}
		glog.Errorf("error creating new book: %v", err)
		http.Error(w, "error creating new book", http.StatusInternalServerError)
	}
}

func addSubject(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to add subjects", http.StatusUnauthorized)
		return
	}
	subject := &models.Subject{}
	err := json.NewDecoder(r.Body).Decode(subject)
	if err != nil {
		glog.Errorf("error while decoding request body: %v", err)
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
		return
	}
	err = dataStore.CreateSubject(*subject)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			glog.Errorf("duplicate entry: %v", err)
			http.Error(w, "duplicate entry", http.StatusBadRequest)
		}
		glog.Errorf("error creating new subject: %v", err)
		http.Error(w, "error creating new subject", http.StatusInternalServerError)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := dataStore.GetBooks()
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getBooksByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	books, err := dataStore.GetBooksByName(name)
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getBookByBookID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		glog.Errorf("error parsing bookID: %v", err)
		http.Error(w, "error parsing bookID", http.StatusInternalServerError)
		return
	}
	books, err := dataStore.GetBookByID(uint(bookID))
	if err != nil {
		glog.Errorf("error fetching book: %v", err)
		http.Error(w, "error fetching book", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getBooksByAuthorID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		glog.Errorf("error parsing authorID: %v", err)
		http.Error(w, "error parsing authorID", http.StatusInternalServerError)
		return
	}
	books, err := dataStore.GetBooksByAuthor(uint(authorID))
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getBooksBySubjectID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subjectID, err := strconv.Atoi(id)
	if err != nil {
		glog.Errorf("error parsing subjectID: %v", err)
		http.Error(w, "error parsing subjectID", http.StatusInternalServerError)
		return
	}
	books, err := dataStore.GetBooksBySubject(uint(subjectID))
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := dataStore.GetSubjects()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			glog.Errorf("no subjects found: %v", err)
			http.Error(w, "no subjects found", http.StatusNoContent)
		} else {
			glog.Errorf("error fetching all subjects: %v", err)
			http.Error(w, "error fetching all subjects", http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(subjects)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := dataStore.GetAuthors()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			glog.Errorf("no authors found: %v", err)
			http.Error(w, "no authors found", http.StatusNoContent)
		} else {
			glog.Errorf("error fetching all authors: %v", err)
			http.Error(w, "error fetching all authors", http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getAuthorByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	authors, err := dataStore.GetAuthorsByName(name)
	if err != nil {
		glog.Errorf("error fetching authors by name: %v", err)
		http.Error(w, "error fetching authors by name", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getAuthorByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		glog.Errorf("error parsing authorID: %v", err)
		http.Error(w, "error parsing authorID", http.StatusInternalServerError)
		return
	}
	author, err := dataStore.GetAuthorByID(uint(authorID))
	if err != nil {
		glog.Errorf("error fetching author by id: %v", err)
		http.Error(w, "error fetching author by id", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(author)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}
