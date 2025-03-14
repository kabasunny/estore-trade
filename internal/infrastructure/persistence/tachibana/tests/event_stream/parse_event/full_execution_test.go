package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_FullExecution(t *testing.T) {
	// テスト用のメッセージ (全部約定の例)
	message := []byte("p_cmd^BEC^Ap_ON^B12345^Ap_IC^B6758^Ap_BBKB^B3^Ap_ODST^B2^Ap_CRPR^B1500.0^Ap_CRSR^B100^Ap_CREXSR^B100^Ap_EXPR^B1500.0^Ap_EXSR^B100")

	// EventStream インスタンスの作成 (Logger はモック)
	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	// parseEvent の呼び出し
	event, err := tachibana.CallParseEvent(es, message)

	// アサーション
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "12345", tachibana.GetOrderTachibanaOrderID(event.Order))
	assert.Equal(t, "6758", tachibana.GetOrderSymbol(event.Order))
	assert.Equal(t, "buy", tachibana.GetOrderSide(event.Order))
	assert.Equal(t, "2", tachibana.GetOrderStatus(event.Order))   //全部執行
	assert.Equal(t, 1500.0, tachibana.GetOrderPrice(event.Order)) //注文価格
	assert.Equal(t, 100, tachibana.GetOrderQuantity(event.Order)) //注文数
	assert.Equal(t, 100, event.Order.FilledQuantity)              //約定数
	assert.Equal(t, 1500.0, event.Order.AveragePrice)             //平均約定価格
	assert.Equal(t, 100, event.Order.FilledQuantity)              //約定数量のテスト
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_FullExecution
