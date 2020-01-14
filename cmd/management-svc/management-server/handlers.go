package management_server

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
)

func GetAuthInfoFromContext(ctx context.Context) *models.AuthInfo {
	return ctx.Value(middleware.ContextAuthInfo).(*models.AuthInfo)
}

func (srv *Server) getCompleteHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_complete_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	history, err := srv.DB.GetCompleteHistory()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_complete_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_complete_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_complete_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_history", err, http.StatusInternalServerError)
		return
	}
	history, err := srv.DB.GetHistory(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) checkAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "check_availability", err, http.StatusInternalServerError)
		return
	}
	avail, err := srv.DB.CheckAvailability(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "check_availability", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "check_availability", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(avail)
	if err != nil {
		handleError(w, ctx, srv, "check_availability", err, http.StatusInternalServerError)
	}
}

func (srv *Server) issueBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "issue_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	book := r.FormValue("bookId")
	user := r.FormValue("userId")
	bookID, err := strconv.Atoi(book)
	if err != nil {
		handleError(w, ctx, srv, "issue_book", err, http.StatusInternalServerError)
		return
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		handleError(w, ctx, srv, "issue_book", err, http.StatusInternalServerError)
		return
	}

	err = srv.DB.IssueBook(uint(bookID), uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "issue_book", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "issue_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book issued successfully!")
	if err != nil {
		handleError(w, ctx, srv, "issue_book", err, http.StatusInternalServerError)
	}
}

func (srv *Server) returnBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "return_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
		return
	}
	err = srv.DB.ReturnBook(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "return_book", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book return processed successfully!")
	if err != nil {
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, ctx context.Context, srv *Server, task string, err error, statusCode int) {
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
