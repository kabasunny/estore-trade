package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_OrderAccepted(t *testing.T) {
	message := []byte("p_cmd^BEC^Ap_ON^B98765^Ap_IC^B4567^Ap_BBKB^B3^Ap_ODST^B1^Ap_CRPR^B100.0^Ap_CRSR^B50") //注文受付

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "98765", tachibana.GetOrderTachibanaOrderID(event.Order))
	assert.Equal(t, "4567", tachibana.GetOrderSymbol(event.Order))
	assert.Equal(t, "buy", tachibana.GetOrderSide(event.Order))
	assert.Equal(t, "1", tachibana.GetOrderStatus(event.Order))  //注文受付
	assert.Equal(t, 100.0, tachibana.GetOrderPrice(event.Order)) //注文価格
	assert.Equal(t, 50, tachibana.GetOrderQuantity(event.Order)) //注文数
	assert.Equal(t, 0, event.Order.FilledQuantity)               //約定数は0
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_OrderAccepted
