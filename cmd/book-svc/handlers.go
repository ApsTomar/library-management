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
)

func GetAuthInfoFromContext(ctx context.Context) *models.AuthInfo {
	return ctx.Value(middleware.ContextAuthInfo).(*models.AuthInfo)
}

func addAuthor(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to add authors", http.StatusUnauthorized)
	}
	author := &models.Author{}
	err := json.NewDecoder(r.Body).Decode(author)
	if err != nil {
		glog.Errorf("error while decoding request body: %v", err)
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
	}
	err = dataStore.CreateAuthor(*author)
	if err != nil {
		glog.Errorf("error creating new author: %v", err)
		http.Error(w, "error creating new author", http.StatusInternalServerError)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to add authors", http.StatusUnauthorized)
	}
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		glog.Errorf("error while decoding request body: %v", err)
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
	}
	err = dataStore.CreateBook(*book)
	if err != nil {
		glog.Errorf("error creating new book: %v", err)
		http.Error(w, "error creating new book", http.StatusInternalServerError)
	}
}

func addSubject(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to add authors", http.StatusUnauthorized)
	}
	subject := &models.Subject{}
	err := json.NewDecoder(r.Body).Decode(subject)
	if err != nil {
		glog.Errorf("error while decoding request body: %v", err)
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
	}
	err = dataStore.CreateSubject(*subject)
	if err != nil {
		glog.Errorf("error creating new subject: %v", err)
		http.Error(w, "error creating new subject", http.StatusInternalServerError)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := dataStore.GetBooks()
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
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
	}
	books, err := dataStore.GetBooksByID(uint(bookID))
	if err != nil {
		glog.Errorf("error fetching book: %v", err)
		http.Error(w, "error fetching book", http.StatusInternalServerError)
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
	}
	books, err := dataStore.GetBooksByAuthor(uint(authorID))
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
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
	}
	books, err := dataStore.GetBooksBySubject(uint(subjectID))
	if err != nil {
		glog.Errorf("error fetching books: %v", err)
		http.Error(w, "error fetching books", http.StatusInternalServerError)
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
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}
