package mocks

import (
	"context"

	"github.com/Lockwarr/codefi/services/links"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateResults(ctx context.Context, results []links.Result) error {
	args := m.Called(results)
	return args.Error(0)
}

func (m *MockRepository) GetBatchResults(ctx context.Context, batchID string) ([]links.Result, error) {
	args := m.Called(batchID)
	return args.Get(0).([]links.Result), args.Error(1)
}

func (m *MockRepository) ListResults(ctx context.Context) map[string][]links.Result {
	args := m.Called()
	return args.Get(0).(map[string][]links.Result)
}
