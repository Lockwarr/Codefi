// Package helpers contains helper functions to be used elsewhere
package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
)

// GatherUrls - expects multi-line text with a valid url on each line
func GatherUrls(r io.ReadCloser) ([]*url.URL, error) {
	urls := make([]*url.URL, 0, 64)

	scanner := bufio.NewScanner(r)
	line := 1

	for scanner.Scan() {
		validatedUrl, err := validateUrl(scanner, line)
		if err != nil {
			return nil, err
		}
		urls = append(urls, validatedUrl)
		line++
	}

	return urls, nil
}

func validateUrl(scanner *bufio.Scanner, line int) (*url.URL, error) {
	parsedURL, err := url.Parse(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("bad url at line %v: %w", line, err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" { // url.Parse doesn't always return an error so some extra checks are needed
		return nil, fmt.Errorf("bad url at line %v %s %w", line, parsedURL, errors.New("invalid url"))
	}
	return parsedURL, nil
}
