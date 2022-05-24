package links

import "context"

// Processor
type Processor interface {
	ProcessBatch(ctx context.Context, req ProcessBatchRequest) ([]Result, error)
	GetBatch(ctx context.Context, req GetBatchRequest) ([]Result, error)
}
