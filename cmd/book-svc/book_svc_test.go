package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	book_server "github.com/library/cmd/book-svc/book-server"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/library/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strconv"
)

type testData struct {
	author  string
	book    string
	subject string
}

var _ = Describe("Book-Service", func() {
	var (
		r            *chi.Mux
		adminToken   string
		userToken    string
		data         *testData
		authorCount  int
		subjectCount int
		bookCount    int
		err          error
	)

	BeforeSuite(func() {
		env = &envConfig.Env{}
		err = envconfig.Process("library", env)
		Expect(err).To(BeNil())
		testRun = true
		dataStore = data_store.DbConnect(env, true)
		srv = book_server.NewServer(env, dataStore, nil)
		srv.TestRun = true
		r = book_server.SetupRouter(srv)
		middleware.SetJwtSigningKey(srv.Env.JwtSigningKey)
		adminToken, userToken, err = setupAuthInfo(env)
		Expect(err).To(BeNil())
		r = book_server.SetupRouter(srv)
		data = &testData{}
	})
	Describe("Handlers Test", func() {
		Describe("Add Author", func() {
			It("Should create a new author in DB", func() {
				authorReq := &models.Author{
					Name:        "testAuthor",
					DateOfBirth: "29 February 1600",
				}
				marshalReq, err := json.Marshal(authorReq)
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodPost, "/admin/add/author", bytes.NewBuffer(marshalReq))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+adminToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				data.author = authorReq.Name
				err = dataStore.Db.Table("author").Count(&authorCount).Error
				Expect(err).To(BeNil())
			})
		})
		Describe("Add Subject", func() {
			It("Should create a new subject in DB", func() {
				subjectReq := &models.Subject{
					Name: "testSubject",
				}
				marshalReq, err := json.Marshal(subjectReq)
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodPost, "/admin/add/subject", bytes.NewBuffer(marshalReq))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+adminToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				data.subject = subjectReq.Name
				err = dataStore.Db.Table("subject").Count(&subjectCount).Error
				Expect(err).To(BeNil())
			})
		})
		Describe("Add Book", func() {
			It("Should create a new book in DB", func() {
				author := &models.Author{}
				err = dataStore.Db.Where("name = 'testAuthor'").First(author).Error
				Expect(err).To(BeNil())
				bookReq := &models.Book{
					Name:     "testBook",
					Subject:  "testSubject",
					AuthorID: strconv.Itoa(int(author.ID)),
				}
				marshalReq, err := json.Marshal(bookReq)
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodPost, "/admin/add/book", bytes.NewBuffer(marshalReq))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+adminToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				data.book = bookReq.Name
				err = dataStore.Db.Table("book").Count(&bookCount).Error
				Expect(err).To(BeNil())
			})
		})
		Describe("Get All Authors", func() {
			It("Should return all the authors", func() {
				req := httptest.NewRequest(http.MethodGet, "/get/authors", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var authors []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&authors)
				Expect(len(authors)).To(BeEquivalentTo(authorCount))
			})
		})
		Describe("Get All Books", func() {
			It("Should return all the books", func() {
				req := httptest.NewRequest(http.MethodGet, "/get/books", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var books []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&books)
				Expect(len(books)).To(BeEquivalentTo(bookCount))
			})
		})
		Describe("Get All Subjects", func() {
			It("Should return all the subjects", func() {
				req := httptest.NewRequest(http.MethodGet, "/get/subjects", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var subjects []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&subjects)
				Expect(len(subjects)).To(BeEquivalentTo(subjectCount))
			})
		})
		Describe("Get Books By Name", func() {
			It("Should return the books matching the name", func() {
				req := httptest.NewRequest(http.MethodGet, "/get/books-by-name/testBook", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var books []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&books)
				Expect(books[0]["name"].(string)).To(BeEquivalentTo("testBook"))
			})
		})
		Describe("Get Book By ID", func() {
			It("Should return the book matching the id", func() {
				book := &models.Book{}
				err = dataStore.Db.Where("name = 'testBook'").First(book).Error
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get/book-by-id/%v", book.ID), nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var books map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&books)
				Expect(books["name"].(string)).To(BeEquivalentTo("testBook"))

			})
		})
		Describe("Get Books By Author", func() {
			It("Should return books by specific author", func() {
				author := &models.Author{}
				err = dataStore.Db.Where("name = 'testAuthor'").First(author).Error
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get/books-by-author/%v", author.ID), nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var books []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&books)
				Expect(books[0]["name"].(string)).To(BeEquivalentTo("testBook"))

			})
		})
		Describe("Get Books By Subject", func() {
			It("Should return books on the specific subject", func() {
				subject := &models.Subject{}
				err = dataStore.Db.Where("name = 'testSubject'").First(subject).Error
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get/books-by-subject/%v", subject.ID), nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var books []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&books)
				Expect(books[0]["name"].(string)).To(BeEquivalentTo("testBook"))

			})
		})
		Describe("Get Author By Name", func() {
			It("Should return author matching the name", func() {
				req := httptest.NewRequest(http.MethodGet, "/get/author-by-name/testAuthor", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var authors []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&authors)
				Expect(authors[0]["name"].(string)).To(BeEquivalentTo("testAuthor"))

			})
		})
		Describe("Get Author By ID", func() {
			It("Should return author matching the id", func() {
				author := &models.Author{}
				err = dataStore.Db.Where("name = 'testAuthor'").First(author).Error
				Expect(err).To(BeNil())
				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get/author-by-id/%v", author.ID), nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				resp := rec.Result()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
				var authors map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&authors)
				Expect(authors["name"].(string)).To(BeEquivalentTo("testAuthor"))

			})
		})
	})
	AfterSuite(func() {
		err = cleanTestData(dataStore.Db, data)
		Expect(err).To(BeNil())
	})
})
