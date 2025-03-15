package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_NumberParseError(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_cmd^BEC^Ap_ON^B12345^Ap_IC^B6758^Ap_BBKB^B3^Ap_ODST^B0^Ap_CRPR^Babc^Ap_CRSR^Bxyz" // 不正な数値

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
	assert.NoError(t, err) // エラーは発生しない
	assert.NotNil(t, event)
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, 0.0, event.Order.Price)  // 修正: パース失敗時はデフォルト値
	assert.Equal(t, 0, event.Order.Quantity) // 修正: パース失敗時はデフォルト値
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_NumberParseError
