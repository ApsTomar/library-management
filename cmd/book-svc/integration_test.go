package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/library/models"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strconv"
	"testing"
)

func TestBookService(t *testing.T) {
	Convey("POST /add-author", t, func() {
		url := fmt.Sprintf("%s/admin/add/author", testServer.URL)
		Convey("It should create a new author", func() {
			authorReq := &models.Author{
				Name:        "testAuthor",
				DateOfBirth: "29 February 1600",
			}
			marshalReq, err := json.Marshal(authorReq)
			So(err, ShouldBeNil)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalReq))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
		})
	})

	Convey("POST /add-subject", t, func() {
		url := fmt.Sprintf("%s/admin/add/subject", testServer.URL)
		Convey("It should create a new subject", func() {
			subjectReq := &models.Subject{
				Name: "testSubject",
			}
			marshalReq, err := json.Marshal(subjectReq)
			So(err, ShouldBeNil)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalReq))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
		})
	})

	Convey("POST /add-book", t, func() {
		url := fmt.Sprintf("%s/admin/add/book", testServer.URL)
		Convey("It should create a new book", func() {
			authors, err := dataStore.GetAuthorsByName("testAuthor")
			So(err, ShouldBeNil)
			testAuthorID = strconv.Itoa(int((*authors)[0].ID))
			bookReq := &models.Book{
				Name:     "testBook",
				Subject:  "testSubject",
				AuthorID: testAuthorID,
			}
			marshalReq, err := json.Marshal(bookReq)
			So(err, ShouldBeNil)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalReq))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
		})
	})

	Convey("POST /get-book-by-name", t, func() {
		url := fmt.Sprintf("%s/get/books-by-name/testBook", testServer.URL)
		Convey("It should retrieve book by name", func() {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+userToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var books []map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&books)
			So(books[0]["name"].(string), ShouldEqual, "testBook")
		})
	})

	Convey("POST /get-author-by-name", t, func() {
		url := fmt.Sprintf("%s/get/author-by-name/testAuthor", testServer.URL)
		Convey("It should retrieve author by name", func() {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+userToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var authors []map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&authors)
			So(authors[0]["name"].(string), ShouldEqual, "testAuthor")
		})
	})

	Convey("POST /get-book-by-author", t, func() {
		url := fmt.Sprintf("%s/get/books-by-author/%s", testServer.URL, testAuthorID)
		Convey("It should retrieve books by author", func() {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+userToken)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			defer resp.Body.Close()
			var books []map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&books)
			So(books[0]["name"].(string), ShouldEqual, "testBook")
		})
	})

}
