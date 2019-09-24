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

func getCompleteHistory(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to manage books", http.StatusUnauthorized)
		return
	}
	history, err := dataStore.GetCompleteHistory()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			glog.Errorf("no history found: %v", err)
			http.Error(w, "no history found", http.StatusInternalServerError)
		} else {
			glog.Errorf("error fetching book history: %v", err)
			http.Error(w, "error fetching book history: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to manage books", http.StatusUnauthorized)
		return
	}
	name := chi.URLParam(r, "name")
	history, err := dataStore.GetHistory(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			glog.Errorf("no history found: %v", err)
			http.Error(w, "no history found", http.StatusInternalServerError)
		} else {
			glog.Errorf("error fetching book history: %v", err)
			http.Error(w, "error fetching book history: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func checkAvailability(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		glog.Errorf("error parsing bookID: %v", err)
		http.Error(w, "error parsing bookID", http.StatusInternalServerError)
		return
	}
	avail, err := dataStore.CheckAvailability(uint(bookID))
	if err != nil {
		glog.Errorf("error checking availability of book: %v", err)
		http.Error(w, "error checking availability of book: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(avail)
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func issueBook(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to manage books", http.StatusUnauthorized)
		return
	}
	book := r.FormValue("bookId")
	user := r.FormValue("userId")
	bookID, err := strconv.Atoi(book)
	if err != nil {
		glog.Errorf("error parsing bookID: %v", err)
		http.Error(w, "error parsing bookID", http.StatusInternalServerError)
		return
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		glog.Errorf("error parsing userID: %v", err)
		http.Error(w, "error parsing userID", http.StatusInternalServerError)
		return
	}

	err = dataStore.IssueBook(uint(bookID), uint(userID))
	if err != nil {
		glog.Errorf("error issuing book: %v", err)
		http.Error(w, "error issuing book: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book issued successfully!")
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}

func returnBook(w http.ResponseWriter, r *http.Request) {
	authInfo := GetAuthInfoFromContext(r.Context())
	if authInfo.Role != models.AdminAccount {
		glog.Errorf("permission denied")
		http.Error(w, "Only Admin is authorized to manage books", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		glog.Errorf("error parsing bookID: %v", err)
		return
	}
	err = dataStore.ReturnBook(uint(bookID))
	if err != nil {
		glog.Errorf("error returning book: %v", err)
		http.Error(w, "error returning book: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book return processed successfully!")
	if err != nil {
		glog.Errorf("error encoding json response: %v", err)
		http.Error(w, "error encoding json response", http.StatusInternalServerError)
	}
}
