package scraper

// Result array of results will be returned after scraping
type Result struct {
	PageURL          string
	InternalLinksNum uint
	ExternalLinksNum uint
	Success          bool
	Error            error
}
