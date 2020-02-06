package book_server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/library/efk"
	"github.com/library/middleware"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAuthInfoFromContext(ctx context.Context) *models.AuthInfo {
	return ctx.Value(middleware.ContextAuthInfo).(*models.AuthInfo)
}

func (srv *Server) addAuthor(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "add_author", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	author := &models.Author{}
	err := json.NewDecoder(r.Body).Decode(author)
	if err != nil {
		handleError(w, ctx, srv, "add_author", err, http.StatusInternalServerError)
		return
	}
	err = srv.DB.CreateAuthor(*author)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, ctx, srv, "add_author", errors.New("duplicate entry"), http.StatusBadRequest)
			return
		}
		handleError(w, ctx, srv, "add_author", err, http.StatusInternalServerError)
	}
}

func (srv *Server) addBook(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "add_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		handleError(w, ctx, srv, "add_book", err, http.StatusInternalServerError)
		return
	}
	book.AvailableDate = time.Now()
	book.Available = true
	err = srv.DB.CreateBook(*book)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, ctx, srv, "add_book", err, http.StatusBadRequest)
			return
		}
		handleError(w, ctx, srv, "add_book", err, http.StatusInternalServerError)
	}
}

func (srv *Server) addSubject(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "add_subject", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	subject := &models.Subject{}
	err := json.NewDecoder(r.Body).Decode(subject)
	if err != nil {
		handleError(w, ctx, srv, "add_subject", err, http.StatusInternalServerError)
		return
	}
	err = srv.DB.CreateSubject(*subject)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, ctx, srv, "add_subject", errors.New("duplicate entry"), http.StatusBadRequest)
			return
		}
		handleError(w, ctx, srv, "add_subject", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooks(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	books, err := srv.DB.GetBooks()
	if err != nil {
		if err == gorm.ErrRecordNotFound || books == nil {
			handleError(w, ctx, srv, "get_books", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByName(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	name := chi.URLParam(r, "name")
	books, err := srv.DB.GetBooksByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_name", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_name", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_name", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBookByBookID(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_book_by_id", err, http.StatusInternalServerError)
		return
	}
	book, err := srv.DB.GetBookByID(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_book_by_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_book_by_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		handleError(w, ctx, srv, "get_book_by_id", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByAuthorID(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_author_id", err, http.StatusInternalServerError)
		return
	}
	books, err := srv.DB.GetBooksByAuthor(uint(authorID))
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_author_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_author_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_author_id", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksBySubjectID(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	subjectID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_subject_id", err, http.StatusInternalServerError)
		return
	}
	books, err := srv.DB.GetBooksBySubject(uint(subjectID))
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_subject_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_subject_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_subject_id", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getSubjects(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	subjects, err := srv.DB.GetSubjects()
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*subjects) == 0 {
			handleError(w, ctx, srv, "get_subjects", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_subjects", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(subjects)
	if err != nil {
		handleError(w, ctx, srv, "get_subjects", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getAuthors(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authors, err := srv.DB.GetAuthors()
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*authors) == 0 {
			handleError(w, ctx, srv, "get_authors", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_authors", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		handleError(w, ctx, srv, "get_authors", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getAuthorByName(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	name := chi.URLParam(r, "name")
	authors, err := srv.DB.GetAuthorsByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*authors) == 0 {
			handleError(w, ctx, srv, "get_author_by_name", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_author_by_name", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		handleError(w, ctx, srv, "get_author_by_name", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getAuthorByID(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_author_by_id", err, http.StatusInternalServerError)
		return
	}
	author, err := srv.DB.GetAuthorByID(uint(authorID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_author_by_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_author_by_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(author)
	if err != nil {
		handleError(w, ctx, srv, "get_author_by_id", err, http.StatusInternalServerError)
	}
}

func (srv *Server) health() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func handleError(w *middleware.LogResponseWriter, ctx context.Context, srv *Server, task string, err error, statusCode int) {
	if !srv.TestRun {
		srv.TracingID = ctx.Value(middleware.RequestTracingID).(string)
		efk.LogError(srv.EfkLogger, srv.EfkTag, srv.TracingID, task, err, statusCode)
	}
	http.Error(w, err.Error(), statusCode)

	logrus.WithFields(logrus.Fields{
		"tracingID":  srv.TracingID,
		"statusCode": statusCode,
		"error":      err,
	}).Error(task)
}
