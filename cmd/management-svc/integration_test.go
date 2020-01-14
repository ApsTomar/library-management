package main

import (
	"encoding/json"
	"fmt"
	"github.com/library/models"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestManagementService(t *testing.T) {
	Convey("GET /check-availability", t, func() {
		getUrl := fmt.Sprintf("%s/user/check-availability/101010", testServer.URL)
		Convey("It should return the availability of the specified book", func() {
			req, err := http.NewRequest(http.MethodGet, getUrl, nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+userToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var avail bool
			err = json.NewDecoder(resp.Body).Decode(&avail)
			So(err, ShouldBeNil)
			So(avail, ShouldEqual, true)
		})
	})

	Convey("POST /issue-book", t, func() {
		postUrl := fmt.Sprintf("%s/admin/issue-book", testServer.URL)
		formData := url.Values{
			"userId": {"101010"},
			"bookId": {"101010"},
		}
		Convey("It should issue the book to the specified user", func() {
			req, err := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(formData.Encode()))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var book models.Book
			err = dataStore.Db.Where("id = ?", "101010").Find(&book).Error
			So(err, ShouldBeNil)
			So(book.Available, ShouldEqual, false)
		})
	})

	Convey("GET /history", t, func() {
		getUrl := fmt.Sprintf("%s/admin/get-history/101010", testServer.URL)
		Convey("It should retrieve the issue history of the specified book", func() {
			req, err := http.NewRequest(http.MethodGet, getUrl, nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var history []map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&history)
			So(err, ShouldBeNil)
			if len(history) > 0 {
				So(history[0]["bookId"], ShouldEqual, 101010)
			}
		})
	})

	Convey("GET /return-book", t, func() {
		getUrl := fmt.Sprintf("%s/admin/return-book/101010", testServer.URL)
		Convey("It should change the availability status of that book", func() {
			req, err := http.NewRequest(http.MethodGet, getUrl, nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var book models.Book
			err = dataStore.Db.Where("id = ?", "101010").Find(&book).Error
			So(err, ShouldBeNil)
			So(book.Available, ShouldEqual, true)
		})
	})
}
