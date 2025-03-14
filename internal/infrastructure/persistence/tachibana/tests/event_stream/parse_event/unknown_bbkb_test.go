package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_UnknownBBKB(t *testing.T) {
	message := []byte("p_cmd^BEC^Ap_ON^B12345^Ap_IC^B6758^Ap_BBKB^B9^Ap_ODST^B0") // 未知のBBKB

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.NoError(t, err) // エラーは発生しない
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "", tachibana.GetOrderSide(event.Order)) // 不明なBBKBの場合は空文字列 (またはデフォルト値)

}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_UnknownBBKB
