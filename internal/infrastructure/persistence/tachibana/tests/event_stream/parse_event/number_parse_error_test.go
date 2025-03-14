package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_NumberParseError(t *testing.T) {
	message := []byte("p_cmd^BEC^Ap_ON^B12345^Ap_IC^B6758^Ap_BBKB^B3^Ap_ODST^B0^Ap_CRPR^Babc^Ap_CRSR^Bxyz") // 不正な数値

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)
	assert.NoError(t, err) // エラーは発生しない
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, 0.0, tachibana.GetOrderPrice(event.Order))  // パース失敗時はデフォルト値
	assert.Equal(t, 0, tachibana.GetOrderQuantity(event.Order)) // パース失敗時はデフォルト値
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_NumberParseError
