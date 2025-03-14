// internal/infrastructure/persistence/tachibana/tests/event_stream/simple_credit_sell_test.go

package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestPlaceSimpleCreditSellOrder(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	err := client.Login(ctx, nil)
	assert.NoError(t, err)
	defer client.Logout(ctx)

	// 1. エントリー注文 (信用成行売り)
	entryOrder := &domain.Order{
		Symbol:     "7974", // 例: 任天堂
		Side:       "short",
		OrderType:  "market",
		TradeType:  "credit_open", // 信用新規
		Quantity:   100,
		MarketCode: "00", // 東証
	}
	placedEntryOrder, err := client.PlaceOrder(ctx, entryOrder)
	assert.NoError(t, err)
	assert.NotEmpty(t, placedEntryOrder.TachibanaOrderID)

}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestPlaceSimpleCreditSellOrder
