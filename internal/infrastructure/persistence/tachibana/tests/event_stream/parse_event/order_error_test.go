package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_OrderError(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_cmd^BEC^Ap_ON^B11111^Ap_IC^B2222^Ap_BBKB^B3^Ap_ODST^B2" // 注文エラー (ODST=2)

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
	assert.Equal(t, "11111", event.Order.TachibanaOrderID) //修正
	assert.Equal(t, "2222", event.Order.Symbol)            //修正
	assert.Equal(t, "long", event.Order.Side)              //修正
	assert.Equal(t, "2", event.Order.Status)               //修正 注文エラーステータス
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_OrderError
