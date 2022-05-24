package repository_test

import (
	"context"
	"testing"

	"github.com/Lockwarr/codefi/services/links"
	"github.com/Lockwarr/codefi/services/links/repository"
	"github.com/stretchr/testify/suite"
)

type inmemoryDBTestSuite struct {
	suite.Suite
	inMemoryDB links.Repository
}

func (s *inmemoryDBTestSuite) SetupTest() {
	s.inMemoryDB = repository.NewInMemoryDB()
}

func (s *inmemoryDBTestSuite) AfterTest(suite string, testName string) {
}

func TestInmemoryDBTestSuite(t *testing.T) {
	suite.Run(t, &inmemoryDBTestSuite{})
}

func (s *inmemoryDBTestSuite) TestCreateResults_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	results := links.Result{
		ID:      "testID",
		BatchID: "testBatchID",
	}
	expectedResults := map[string][]links.Result{"testBatchID": {results}}
	// Act
	err := s.inMemoryDB.CreateResults(ctx, []links.Result{results})
	actualResults := s.inMemoryDB.ListResults(ctx)

	// Assert
	s.Equal(nil, err)
	s.Equal(expectedResults, actualResults)
}

func (s *inmemoryDBTestSuite) TestCreateResults_WhenZeroResults_ThenFail() {
	// Arrange
	ctx := context.Background()

	// Act
	err := s.inMemoryDB.CreateResults(ctx, []links.Result{})

	// Assert
	s.Equal("no results were passed", err.Error())
}

func (s *inmemoryDBTestSuite) TestGetBatchResults_ThenSuccess() {
	// Arrange
	batchID := "testBatchID"
	ctx := context.Background()
	results := links.Result{
		ID:      "testID",
		BatchID: batchID,
	}
	expectedResults := []links.Result{results}
	// Act
	_ = s.inMemoryDB.CreateResults(ctx, expectedResults)
	actualResults, err := s.inMemoryDB.GetBatchResults(ctx, batchID)

	// Assert
	s.Equal(nil, err)
	s.Equal(expectedResults[0].BatchID, actualResults[0].BatchID)
	s.Equal(len(expectedResults), len(actualResults))
}

func (s *inmemoryDBTestSuite) TestGetBatchResults_WhenNotExistingBatchIDPassed_ThenFail() {
	// Arrange
	batchID := "testBatchID"
	ctx := context.Background()

	// Act
	actualResults, err := s.inMemoryDB.GetBatchResults(ctx, batchID)

	// Assert
	s.Equal(repository.ErrBatchNotFound, err)
	s.Equal(0, len(actualResults))
}
