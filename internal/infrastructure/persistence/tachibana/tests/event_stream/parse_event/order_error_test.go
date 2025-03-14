package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_OrderError(t *testing.T) {
	message := []byte("p_cmd^BEC^Ap_ON^B11111^Ap_IC^B2222^Ap_BBKB^B3^Ap_ODST^B2") // 注文エラー (ODST=2)

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.NoError(t, err) // エラーは発生しない
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "11111", tachibana.GetOrderTachibanaOrderID(event.Order))
	assert.Equal(t, "2222", tachibana.GetOrderSymbol(event.Order))
	assert.Equal(t, "buy", tachibana.GetOrderSide(event.Order))
	assert.Equal(t, "2", tachibana.GetOrderStatus(event.Order)) // 注文エラーステータス
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_OrderError
