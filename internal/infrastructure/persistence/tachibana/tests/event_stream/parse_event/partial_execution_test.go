package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_PartialExecution(t *testing.T) {
	message := []byte("p_cmd^BEC^Ap_ON^B67890^Ap_IC^B7203^Ap_BBKB^B1^Ap_ODST^B1^Ap_CRPR^B2500.0^Ap_CRSR^B200^Ap_CREXSR^B50^Ap_EXPR^B2500.0^Ap_EXSR^B50")

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "67890", tachibana.GetOrderTachibanaOrderID(event.Order))
	assert.Equal(t, "7203", tachibana.GetOrderSymbol(event.Order))
	assert.Equal(t, "sell", tachibana.GetOrderSide(event.Order))
	assert.Equal(t, "1", tachibana.GetOrderStatus(event.Order)) //一部執行
	assert.Equal(t, 2500.0, tachibana.GetOrderPrice(event.Order))
	assert.Equal(t, 200, tachibana.GetOrderQuantity(event.Order))
	assert.Equal(t, 50, event.Order.FilledQuantity)
	assert.Equal(t, 2500.0, event.Order.AveragePrice)
	assert.Equal(t, 50, event.Order.FilledQuantity) //約定数量のテスト
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_PartialExecution
