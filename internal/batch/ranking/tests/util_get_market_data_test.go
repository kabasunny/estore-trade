package ranking_test

import (
	"context"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMarketData is a simple test to ensure GetMarketData retrieves market data.
func TestGetMarketData(t *testing.T) {
	// Setup a mock TachibanaClient
	mockClient := &tachibana.MockTachibanaClient{}

	// Define a slice of issue codes for the test
	issueCodes := []string{"1234", "5678"}

	// Call GetMarketData with the mock client and issue codes
	marketData, err := ranking.GetMarketData(context.Background(), mockClient, issueCodes)

	fmt.Println(marketData)

	// Assert that no error occurred
	assert.NoError(t, err, "GetMarketData should not return an error")

	// Assert that marketData is not nil
	assert.NotNil(t, marketData, "Market data should not be nil")

	// Assert that marketData contains exactly the number of issue codes requested
	assert.Len(t, marketData, len(issueCodes), "Length of market data should be 0")

	// Assert that the issue codes in marketData match the requested ones

}
