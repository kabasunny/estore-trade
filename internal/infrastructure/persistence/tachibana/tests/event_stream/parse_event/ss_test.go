package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"
	// domain パッケージをインポート

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ３．（７）SS 通知例
func TestParseEventSS_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no^B1^Ap_date^B2020.06.18-07:30:34.810^Ap_errno^B0^Ap_err^B^Ap_cmd^BSS^Ap_PV^BMSGSV^Ap_ENO^B3^Ap_ALT^B0^A" +
		"p_CT^B20200618052959^Ap_LK^B1^Ap_SS^B1\n"

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
	assert.Equal(t, "SS", event.EventType)
	assert.Equal(t, "3", event.EventNo) // EventNo (文字列型)
	assert.Equal(t, false, event.IsFirstEvent)
	assert.Equal(t, "MSGSV", event.Provider) //修正

	// SystemStatus 構造体のフィールドの確認 (event.System が nil でないことを確認してからアクセス)
	if assert.NotNil(t, event.System) {
		assert.Equal(t, "20200618052959", event.System.UpdateDate) //SystemStatusのUpdateDate
		assert.Equal(t, "1", event.System.LoginPermission)         // SystemStatus の LoginPermission
		assert.Equal(t, "1", event.System.SystemState)             // SystemStatus の SystemState
	}
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventSS_Example
