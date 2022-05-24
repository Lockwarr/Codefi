package helpers_test

import (
	"errors"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Lockwarr/codefi/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGatherUrls(t *testing.T) {
	// Arrange
	tests := []struct {
		name                 string
		testData             io.ReadCloser
		expectedGatheredUrls []string
		err                  error
	}{
		{
			name: "successfully gather urls",
			testData: io.NopCloser(strings.NewReader(`https://www.google.com
https://www.facebook.com`)),
			expectedGatheredUrls: []string{"www.google.com", "www.facebook.com"},
			err:                  nil,
		},
		{
			name:                 "gather bad urls",
			testData:             io.NopCloser(strings.NewReader(`not an url`)),
			expectedGatheredUrls: nil,
			err:                  errors.New(`bad url at line 1 not%20an%20url invalid url`),
		},
		{
			name: "gather urls with invalid characters",
			testData: io.NopCloser(strings.NewReader(`https://ww�w.asd.com`)),
			expectedGatheredUrls: nil,
			err:                  errors.New(`should throw error bad url at inner parse`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			gatheredUrls, err := helpers.GatherUrls(tt.testData)

			// Assert
			switch tt.err {
			case nil:
				assert.Equal(t, tt.err, err)
				assert.Equal(t, tt.expectedGatheredUrls[0], gatheredUrls[0].Host)
				assert.Equal(t, tt.expectedGatheredUrls[1], gatheredUrls[1].Host)
			default:
				assert.Error(t, err)
				assert.Equal(t, []*url.URL([]*url.URL(nil)), gatheredUrls)
			}
		})
	}

}

func TestGatherUrlsFromFile(t *testing.T) {
	// Arrange
	file, err := os.Open(`testdata\urls.txt`)
	assert.NoError(t, err)

	// Act
	gatheredUrls, err := helpers.GatherUrls(file)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "www.google.com", gatheredUrls[0].Host)
	assert.Equal(t, "www.facebook.com", gatheredUrls[1].Host)
}
