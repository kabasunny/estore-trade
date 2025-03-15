package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_OrderAccepted(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_cmd^BEC^Ap_ON^B98765^Ap_IC^B4567^Ap_BBKB^B3^Ap_ODST^B1^Ap_CRPR^B100.0^Ap_CRSR^B50" //注文受付

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
	assert.Equal(t, "98765", event.Order.TachibanaOrderID) //修正
	assert.Equal(t, "4567", event.Order.Symbol)            //修正
	assert.Equal(t, "long", event.Order.Side)              //修正
	assert.Equal(t, "1", event.Order.Status)               //修正  注文受付
	assert.Equal(t, 100.0, event.Order.Price)              //修正
	assert.Equal(t, 50, event.Order.Quantity)              //修正
	assert.Equal(t, 0, event.Order.FilledQuantity)         //約定数は0
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_OrderAccepted
