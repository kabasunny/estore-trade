package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_EmptyEventType(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_ON^B12345" // p_cmd がない

	// 制御文字をエスケープシーケンスに置換
	replacer := strings.NewReplacer("^B", "\x02")
	replacedMessage := replacer.Replace(originalMessage)
	message := []byte(replacedMessage)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err) //logger作成のエラーキャッチ
	}
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.Error(t, err) // エラーが発生するはず
	assert.Nil(t, event)
	assert.EqualError(t, err, "event type is empty: p_ON\x0212345") // エラーメッセージの確認(修正)
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_EmptyEventType
