//go:build component
// +build component

package component_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"

	"github.com/Lockwarr/codefi/pkg/scraper"
	"github.com/Lockwarr/codefi/services/links"
	"github.com/Lockwarr/codefi/services/links/domain"
	"github.com/Lockwarr/codefi/services/links/handler"
	"github.com/Lockwarr/codefi/services/links/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
	scraper   *scraper.Scraper
	processor links.Processor
	repo      links.Repository
	listener  net.Listener
	router    *chi.Mux
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	// Arrange
	var err error
	s.router = chi.NewRouter()
	s.listener, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		s.T().Fatal(err)
	}
	srv := &http.Server{Handler: s.router}

	s.repo = repository.NewInMemoryDB()
	s.scraper = scraper.NewScraper()
	s.processor = domain.NewLinksProcessor(s.repo, s.scraper)
	h := handler.NewHandler(s.processor)

	s.router.Route("/api/v1/", func(r chi.Router) {
		r.Post("/links", h.ProcessBatch)
		r.Get("/links/{batchID}", h.GetBatch)
	})

	go func() {
		srv.Serve(s.listener)
	}()
}

func (s *e2eTestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	p.Signal(syscall.SIGINT)
}

func (s *e2eTestSuite) SetupTest() {
}

func (s *e2eTestSuite) TearDownTest() {
}

func (s *e2eTestSuite) Test_EndToEnd_SuccessfulExtraction() {
	// Arrange
	req := createRequestWithAttachedFile("POST", "http://"+s.listener.Addr().String()+"/api/v1/links", `testdata\goodUrls.txt`)

	// Act
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	// Assert
	byteBody, err := ioutil.ReadAll(resp.Body)
	data := links.Response{}
	err = json.Unmarshal(byteBody, &data)
	s.NoError(err)

	jsonString, _ := json.Marshal(data.Data)
	results := links.ProcessBatchResponse{}
	err = json.Unmarshal(jsonString, &results)
	s.NoError(err)

	s.Equal(2, len(results.Results))
	s.Equal(true, results.Results[0].Success)
	s.Equal(true, results.Results[1].Success)
}

func (s *e2eTestSuite) Test_EndToEnd_BadUrls_ThenFail() {
	// Arrange
	req := createRequestWithAttachedFile("POST", "http://"+s.listener.Addr().String()+"/api/v1/links", `testdata\badUrls.txt`)

	// Act
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	// Assert
	byteBody, err := ioutil.ReadAll(resp.Body)
	data := links.Response{}
	err = json.Unmarshal(byteBody, &data)
	s.NoError(err)

	jsonString, _ := json.Marshal(data.Data)
	results := links.ProcessBatchResponse{}
	err = json.Unmarshal(jsonString, &results)
	s.NoError(err)

	s.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	s.Equal(0, len(results.Results))
}

func (s *e2eTestSuite) Test_EndToEnd_SuccessfulBatchRetrieval() {
	// Arrange
	expectedResults := []links.Result{
		{
			ID:      "test123",
			BatchID: "test1235",
			Success: true,
		},
	}
	err := s.repo.CreateResults(context.Background(), expectedResults)
	s.NoError(err)

	req, err := http.NewRequest("GET", "http://"+s.listener.Addr().String()+"/api/v1/links/"+expectedResults[0].BatchID, nil)
	s.NoError(err)

	// Act
	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	// Assert
	byteBody, err := ioutil.ReadAll(resp.Body)
	data := links.Response{}
	err = json.Unmarshal(byteBody, &data)
	s.NoError(err)

	jsonString, _ := json.Marshal(data.Data)
	results := links.GetBatchResponse{}
	err = json.Unmarshal(jsonString, &results)
	s.NoError(err)

	s.Equal(1, len(results.Results))
	s.Equal(true, results.Results[0].Success)
	s.Equal("test123", results.Results[0].ID)
	s.Equal("test1235", results.Results[0].BatchID)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *e2eTestSuite) Test_EndToEnd_NotExistingBatchRetrieval() {
	// Arrange
	req, err := http.NewRequest("GET", "http://"+s.listener.Addr().String()+"/api/v1/links/noBatch", nil)
	s.NoError(err)

	// Act
	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	// Assert
	byteBody, err := ioutil.ReadAll(resp.Body)
	data := links.Response{}
	err = json.Unmarshal(byteBody, &data)
	s.NoError(err)

	jsonString, _ := json.Marshal(data.Data)
	results := links.GetBatchResponse{}
	err = json.Unmarshal(jsonString, &results)
	s.NoError(err)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.Equal(0, len(results.Results))
}

func createRequestWithAttachedFile(method, urlPath, filename string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("urlsFile", "testFile.txt")
	if err != nil {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return nil
	}
	writer.Close()
	req := httptest.NewRequest(method, urlPath, bytes.NewReader(body.Bytes()))
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req
}
