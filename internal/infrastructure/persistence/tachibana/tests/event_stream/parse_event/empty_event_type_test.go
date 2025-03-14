package parse_event_test

import (
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_EmptyEventType(t *testing.T) {
	message := []byte("p_ON^B12345") // p_cmd がない

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.Error(t, err) // エラーが発生するはず
	assert.Nil(t, event)
	assert.EqualError(t, err, "event type is empty: p_ON^B12345") // エラーメッセージの確認
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_EmptyEventType
