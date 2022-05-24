package mocks

import (
	"context"

	"github.com/Lockwarr/codefi/services/links"
	"github.com/stretchr/testify/mock"
)

type MockLinksProcessor struct {
	mock.Mock
}

func (m *MockLinksProcessor) ProcessBatch(ctx context.Context, req links.ProcessBatchRequest) ([]links.Result, error) {
	args := m.Called(req)
	return args.Get(0).([]links.Result), args.Error(1)
}

func (m *MockLinksProcessor) GetBatch(ctx context.Context, req links.GetBatchRequest) ([]links.Result, error) {
	args := m.Called(req)
	return args.Get(0).([]links.Result), args.Error(1)
}
