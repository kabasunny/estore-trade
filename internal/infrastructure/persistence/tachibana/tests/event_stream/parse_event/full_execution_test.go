package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseEventEC_FullExecution(t *testing.T) {
	// テスト用のメッセージ (全部約定の例) - 人間が読める形式
	originalMessage := "p_no^B1^Ap_cmd^BEC^Ap_ON^B12345^Ap_IC^B6758^Ap_BBKB^B3^Ap_ODST^B2^Ap_CRPR^B1500.0^Ap_CRSR^B100^Ap_CREXSR^B100^Ap_EXPR^B1500.0^Ap_EXSR^B100\n"

	// 制御文字をエスケープシーケンスに置換  <-- strings.NewReplacer を使う
	replacer := strings.NewReplacer("^A", "\x01", "^B", "\x02")
	replacedMessage := replacer.Replace(originalMessage)
	message := []byte(replacedMessage)

	// EventStream インスタンスの作成 (Logger はモック)
	logger, err := zap.NewDevelopment() // err を捕捉
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	es := tachibana.NewTestEventStream(logger)

	// parseEvent の呼び出し
	event, err := tachibana.CallParseEvent(es, message)

	// アサーション
	assert.NoError(t, err)  //<-- エラーは発生しない、に修正
	assert.NotNil(t, event) //<-- eventはNotNil、に修正
	assert.Equal(t, "EC", event.EventType)
	assert.NotNil(t, event.Order)
	assert.Equal(t, "12345", event.Order.TachibanaOrderID) //修正
	assert.Equal(t, "6758", event.Order.Symbol)            //修正
	assert.Equal(t, "long", event.Order.Side)              //修正 "long" を期待
	assert.Equal(t, "2", event.Order.Status)               //修正
	assert.Equal(t, 1500.0, event.Order.Price)             //修正
	assert.Equal(t, 100, event.Order.Quantity)             //修正
	assert.Equal(t, 100, event.Order.FilledQuantity)       //約定数
	assert.Equal(t, 1500.0, event.Order.AveragePrice)      //平均約定価格
	assert.Equal(t, 100, event.Order.FilledQuantity)       //約定数量のテスト
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventEC_FullExecution
