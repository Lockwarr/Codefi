package scraper

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/html"
)

var concurencyLimit = 1000

// ScraperService ...
type ScraperService interface {
	Scrape(ctx context.Context, urls []*url.URL) []Result
}

type Scraper struct {
	httpClient *http.Client
	wg         *sync.WaitGroup
}

func NewScraper() *Scraper {
	return &Scraper{httpClient: cleanhttp.DefaultClient(), wg: &sync.WaitGroup{}}
}

// Scrape - adds to waitgroup and starts a goroutine for each url that we have passed
// each goroutine decrements waitgroup when finished. It is finished when our scraping worker
// finishes and sends to the channel.
// After all goroutines finish, we iterate through the channel to gather the results
func (s *Scraper) Scrape(ctx context.Context, urls []*url.URL) []Result {
	resultsChan := make(chan Result)
	results := make([]Result, 0, len(urls))

	for _, url := range urls {
		s.wg.Add(1)
		url := url // https://go.dev/doc/faq#closures_and_goroutines

		go func() {
			defer s.wg.Done()
			select {
			case resultsChan <- s.startScrapingWorker(ctx, url, resultsChan):
			}
		}()
	}

	go func() {
		s.wg.Wait()
		close(resultsChan)
	}()

	// wait for results
	for res := range resultsChan {
		results = append(results, res)
	}
	log.Println("Successfully fetched all internal and external links.")
	return results
}

func (s *Scraper) startScrapingWorker(ctx context.Context, url *url.URL, resultsChan chan Result) Result {
	result := Result{PageURL: url.String(), Success: false}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		result.Error = err
		return result
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		result.Error = errors.New("bad status code")
		return result
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		result.Error = err
		return result
	}

	external, internal, err := CountLinks(url, document)
	if err != nil {
		result.Error = err
		result.ExternalLinksNum = external
		result.InternalLinksNum = internal
		return result
	}

	result.ExternalLinksNum = external
	result.InternalLinksNum = internal
	result.Success = true

	return result
}
