package scraper

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestCountLinks(t *testing.T) {
	tests := []struct {
		name           string
		wantedExternal uint
		wantedInternal uint
		wantedErr      error
		url            string
		setup          func(t *testing.T, url string) (*url.URL, *html.Node)
	}{
		{
			name:           "successfully count links",
			url:            "http://localhost.com/",
			wantedExternal: 3,
			wantedInternal: 2,
			wantedErr:      nil,
			setup: func(t *testing.T, baseURL string) (*url.URL, *html.Node) {
				parsedBaseURL, err := url.Parse(baseURL)
				if err != nil {
					t.Error(err)
				}

				f, err := os.Open("testdata/links_success.html")
				if err != nil {
					t.Error(err)
				}

				defer f.Close()
				document, err := html.Parse(f)
				if err != nil {
					t.Error(err)
				}

				return parsedBaseURL, document
			},
		},
		{
			name:           "successfully count links when some links are malformed",
			url:            "http://localhost/",
			wantedExternal: 2,
			wantedInternal: 1,
			wantedErr:      nil,
			setup: func(t *testing.T, baseURL string) (*url.URL, *html.Node) {
				parsedBaseURL, err := url.Parse(baseURL)
				if err != nil {
					t.Error(err)
				}
				f, err := os.Open("testdata/links_malformed.html")
				if err != nil {
					t.Error(err)
				}
				defer f.Close()
				document, err := html.Parse(f)
				if err != nil {
					t.Error(err)
				}

				return parsedBaseURL, document
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL, document := tt.setup(t, tt.url)
			actualExternal, actualInternal, err := CountLinks(testURL, document)
			assert.Equal(t, tt.wantedErr, err)
			assert.Equal(t, tt.wantedExternal, actualExternal)
			assert.Equal(t, tt.wantedInternal, actualInternal)
		})
	}
}
