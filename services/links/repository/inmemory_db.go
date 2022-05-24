package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Lockwarr/codefi/services/links"
)

var (
	ErrBatchNotFound = errors.New("batch of results not found")
)

type inMemoryDB struct {
	results map[string][]links.Result
	rw      *sync.RWMutex
}

// NewInMemoryDB ..
func NewInMemoryDB() links.Repository {
	return &inMemoryDB{results: map[string][]links.Result{}, rw: &sync.RWMutex{}}
}

// ListResults - lists all results that we have so far
func (mem *inMemoryDB) ListResults(ctx context.Context) map[string][]links.Result {
	return mem.results
}

// CreateResults - assign batch id to the processed urls and save results
func (mem *inMemoryDB) CreateResults(ctx context.Context, results []links.Result) error {
	mem.rw.Lock()
	defer mem.rw.Unlock()

	if len(results) == 0 {
		return errors.New("no results were passed")
	}

	// we only call this function with urls processed in the same batch
	mem.results[results[0].BatchID] = results

	return nil
}

// GetBatchResults - get batch of processed urls by batch id
// if it doesn't exists an error is returned
func (r *inMemoryDB) GetBatchResults(ctx context.Context, batchID string) ([]links.Result, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	results, ok := r.results[batchID]
	if !ok {
		return nil, ErrBatchNotFound
	}

	return results, nil
}
