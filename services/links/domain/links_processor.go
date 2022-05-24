package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/Lockwarr/codefi/pkg/scraper"
	"github.com/Lockwarr/codefi/services/links"
	"github.com/google/uuid"
)

type linkProcessor struct {
	scraperClient scraper.ScraperService
	repo          links.Repository
}

// NewLinksProcessor ..
func NewLinksProcessor(repo links.Repository, scraperClient scraper.ScraperService) links.Processor {
	return &linkProcessor{scraperClient: scraperClient, repo: repo}
}

// ProcessBatch - process batch of urls to find external and internal links
func (p *linkProcessor) ProcessBatch(ctx context.Context, req links.ProcessBatchRequest) ([]links.Result, error) {
	batchResults := []links.Result{}
	batchID := uuid.NewString()

	results := p.scraperClient.Scrape(context.Background(), req.URLs)

	for _, result := range results {
		batchResults = append(batchResults, links.Result{
			ID:               uuid.NewString(),
			BatchID:          batchID,
			PageURL:          result.PageURL,
			InternalLinksNum: result.InternalLinksNum,
			ExternalLinksNum: result.ExternalLinksNum,
			Success:          result.Success,
			Error:            result.Error,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		})
	}

	err := p.repo.CreateResults(ctx, batchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to create results %w", err)
	}

	return batchResults, err
}

// GetBatch - get batch of urls results
func (s *linkProcessor) GetBatch(ctx context.Context, req links.GetBatchRequest) ([]links.Result, error) {
	results, err := s.repo.GetBatchResults(ctx, req.BatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch %w", err)
	}

	return results, nil
}
