package domain_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	pkgmocks "github.com/Lockwarr/codefi/pkg/mocks"
	"github.com/Lockwarr/codefi/pkg/scraper"
	"github.com/Lockwarr/codefi/services/links"
	"github.com/Lockwarr/codefi/services/links/domain"
	"github.com/Lockwarr/codefi/services/links/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type linkProcessorTestSuite struct {
	suite.Suite
	mockScraperClient *pkgmocks.MockScraper
	mockRepo          *mocks.MockRepository
	linkProcessor     links.Processor
}

func (s *linkProcessorTestSuite) SetupTest() {
	s.mockScraperClient = new(pkgmocks.MockScraper)
	s.mockRepo = new(mocks.MockRepository)
	s.linkProcessor = domain.NewLinksProcessor(s.mockRepo, s.mockScraperClient)
}

func (s *linkProcessorTestSuite) AfterTest(suite string, testName string) {
	s.mockRepo.AssertExpectations(s.T())
	s.mockScraperClient.AssertExpectations(s.T())
}

func TestLinkProcessorTestSuite(t *testing.T) {
	suite.Run(t, &linkProcessorTestSuite{})
}

func (s *linkProcessorTestSuite) TestProcessBatch_ThenProcessSucessfully() {
	// Arrange
	urlGenerated, _ := url.Parse("http://google.com")
	expectedScraperResult := scraper.Result{
		PageURL: "test1",
	}

	s.mockScraperClient.On("Scrape", []*url.URL{urlGenerated}).Return([]scraper.Result{expectedScraperResult}, nil)
	s.mockRepo.On("CreateResults", mock.Anything).Return(nil)

	// Act
	res, err := s.linkProcessor.ProcessBatch(context.Background(), links.ProcessBatchRequest{URLs: []*url.URL{urlGenerated}})

	// Assert
	s.Equal(nil, err)
	s.Equal(1, len(res))
}

func (s *linkProcessorTestSuite) TestProcessBatch_WhenCreateResultsFails_ThenFail() {
	// Arrange
	urlGenerated, _ := url.Parse("http://google.com")
	expectedScraperResult := scraper.Result{
		PageURL: "test1",
	}

	s.mockScraperClient.On("Scrape", []*url.URL{urlGenerated}).Return([]scraper.Result{expectedScraperResult}, nil)
	s.mockRepo.On("CreateResults", mock.Anything).Return(errors.New("error"))

	// Act
	res, err := s.linkProcessor.ProcessBatch(context.Background(), links.ProcessBatchRequest{URLs: []*url.URL{urlGenerated}})

	// Assert
	s.Equal("failed to create results error", err.Error())
	s.Equal([]links.Result([]links.Result(nil)), res)
}

func (s *linkProcessorTestSuite) TestGetBatch_ThenSucess() {
	// Arrange
	batchID := "batchID"

	s.mockRepo.On("GetBatchResults", batchID).Return([]links.Result{}, nil)

	// Act
	res, err := s.linkProcessor.GetBatch(context.Background(), links.GetBatchRequest{BatchID: batchID})

	// Assert
	s.Equal(nil, err)
	s.Equal(0, len(res))
}

func (s *linkProcessorTestSuite) TestGetBatch_WhenGetBatchResultsFail_ThenFail() {
	// Arrange
	batchID := "batchID"

	s.mockRepo.On("GetBatchResults", batchID).Return([]links.Result{}, errors.New("error"))

	// Act
	res, err := s.linkProcessor.GetBatch(context.Background(), links.GetBatchRequest{BatchID: batchID})

	// Assert
	s.Equal("failed to get batch error", err.Error())
	s.Equal([]links.Result([]links.Result(nil)), res)
}
