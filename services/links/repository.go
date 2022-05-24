package links

import "context"

// Repository
type Repository interface {
	CreateResults(ctx context.Context, results []Result) error
	GetBatchResults(ctx context.Context, batchID string) ([]Result, error)
	ListResults(ctx context.Context) map[string][]Result
	// GetResult(ctx context.Context, resultID string) (Result, error) TODO with real db
}
