package mocks

import (
	"context"
	"net/url"

	"github.com/Lockwarr/codefi/pkg/scraper"
	"github.com/stretchr/testify/mock"
)

type MockScraper struct {
	mock.Mock
}

func (m *MockScraper) Scrape(ctx context.Context, urls []*url.URL) []scraper.Result {
	args := m.Called(urls)
	return args.Get(0).([]scraper.Result)
}
