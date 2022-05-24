package handler_test

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/Lockwarr/codefi/services/links"
	"github.com/Lockwarr/codefi/services/links/handler"
	"github.com/Lockwarr/codefi/services/links/mocks"
	"github.com/Lockwarr/codefi/services/links/repository"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	mockLinkProcessor *mocks.MockLinksProcessor
	handler           *handler.Handler
}

func (s *handlerTestSuite) SetupTest() {
	s.mockLinkProcessor = new(mocks.MockLinksProcessor)
	s.handler = handler.NewHandler(s.mockLinkProcessor)
}

func (s *handlerTestSuite) AfterTest(suite string, testName string) {
	s.mockLinkProcessor.AssertExpectations(s.T())
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, &handlerTestSuite{})
}

func (s *handlerTestSuite) TestProcessBatch_DifferentCases_ThenItIsHandledAsExpected() {
	testCases := []struct {
		name           string
		req            *http.Request
		rr             *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name:           "successful results",
			req:            createRequestWithAttachedFile("POST", "/api/v1/links", `testdata\testFile.txt`, false),
			rr:             httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "fail on formFile",
			req:            createRequestWithAttachedFile("POST", "/api/v1/links", `testdata\testFile.txt`, true),
			rr:             httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "fail on gatherUrls",
			req:            createRequestWithAttachedFile("POST", "/api/v1/links", `testdata\badFile.txt`, false),
			rr:             httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "pass emtpy file",
			req:            createRequestWithAttachedFile("POST", "/api/v1/links", `testdata\emptyFile.txt`, false),
			rr:             httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "fail on processor processBatch",
			req:            createRequestWithAttachedFile("POST", "/api/v1/links", `testdata\testFile.txt`, false),
			rr:             httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Arrange
			urlGenerated, _ := url.Parse("https://www.google.com")
			urls := []*url.URL{urlGenerated}
			switch tc.name {
			case "successful results":
				s.mockLinkProcessor.On("ProcessBatch", links.ProcessBatchRequest{URLs: urls}).Return([]links.Result{}, nil)
			case "fail on processor processBatch":
				s.mockLinkProcessor.On("ProcessBatch", links.ProcessBatchRequest{URLs: urls}).Return([]links.Result{}, errors.New("processor fails"))
			default:
			}

			// Act
			s.handler.ProcessBatch(tc.rr, tc.req)

			// Assert
			s.Equal(tc.expectedStatus, tc.rr.Code)
			s.ResetMocks()
		})
	}

}

func (s *handlerTestSuite) ResetMocks() {
	s.mockLinkProcessor = new(mocks.MockLinksProcessor)
	s.handler = handler.NewHandler(s.mockLinkProcessor)
}

func (s *handlerTestSuite) TestGetBatch_WhenRequestIsCorrect_ThenItIsHandled() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	s.mockLinkProcessor.On("GetBatch", links.GetBatchRequest{BatchID: ""}).Return([]links.Result{}, nil)

	// Act
	s.handler.GetBatch(rr, req)

	// Assert
	s.Equal(http.StatusOK, rr.Code)
}

func (s *handlerTestSuite) TestGetBatch_WhenProcessorGetBatchFailsWithInternal_ThenItIsHandled() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	s.mockLinkProcessor.On("GetBatch", links.GetBatchRequest{BatchID: ""}).Return([]links.Result{}, links.ErrInternalServerError)

	// Act
	s.handler.GetBatch(rr, req)

	// Assert
	s.Equal(http.StatusInternalServerError, rr.Code)
}

func (s *handlerTestSuite) TestGetBatch_WhenProcessorGetBatchFails_ThenItIsHandled() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	res := links.Result{
		ID:      "test1",
		BatchID: "testID",
	}

	s.mockLinkProcessor.On("GetBatch", links.GetBatchRequest{BatchID: ""}).Return([]links.Result{res}, repository.ErrBatchNotFound)

	// Act
	s.handler.GetBatch(rr, req)

	// Assert
	s.Equal(http.StatusNotFound, rr.Code)
}

//
func createRequestWithAttachedFile(method, urlPath, filename string, emptyBody bool) *http.Request {
	if emptyBody {
		return httptest.NewRequest(method, urlPath, nil)
	}
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
