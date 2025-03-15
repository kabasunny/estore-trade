package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_PartialExecution(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_cmd^BEC^Ap_ON^B67890^Ap_IC^B7203^Ap_BBKB^B1^Ap_ODST^B1^Ap_CRPR^B2500.0^Ap_CRSR^B200^Ap_CREXSR^B50^Ap_EXPR^B2500.0^Ap_EXSR^B50"

	// 制御文字をエスケープシーケンスに置換
	replacer := strings.NewReplacer("^A", "\x01", "^B", "\x02")
	replacedMessage := replacer.Replace(originalMessage)
	message := []byte(replacedMessage)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "67890", event.Order.TachibanaOrderID) //修正
	assert.Equal(t, "7203", event.Order.Symbol)            //修正
	assert.Equal(t, "short", event.Order.Side)             //修正
	assert.Equal(t, "1", event.Order.Status)               //修正, 全部執行
	assert.Equal(t, 2500.0, event.Order.Price)             //修正
	assert.Equal(t, 200, event.Order.Quantity)             //修正
	assert.Equal(t, 50, event.Order.FilledQuantity)
	assert.Equal(t, 2500.0, event.Order.AveragePrice)
	assert.Equal(t, 50, event.Order.FilledQuantity) //約定数量のテスト
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_PartialExecution
