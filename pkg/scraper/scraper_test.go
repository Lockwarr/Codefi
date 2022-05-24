package scraper_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/Lockwarr/codefi/pkg/scraper"
	"github.com/stretchr/testify/suite"
)

type scraperTestSuite struct {
	suite.Suite
	scraper scraper.ScraperService
}

func (s *scraperTestSuite) SetupTest() {
	s.scraper = scraper.NewScraper()
}

func (s *scraperTestSuite) AfterTest(suite string, testName string) {
}

func TestScraperTestSuite(t *testing.T) {
	suite.Run(t, &scraperTestSuite{})
}

func (s *scraperTestSuite) TestCreateResults_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	urlGenerated, _ := url.Parse("http://google.com")

	//Act
	actualResults := s.scraper.Scrape(ctx, []*url.URL{urlGenerated})

	// Assert
	s.Equal(1, len(actualResults))
}
