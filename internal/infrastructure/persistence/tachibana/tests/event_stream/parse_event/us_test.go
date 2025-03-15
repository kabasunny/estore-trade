package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"
	// domain パッケージをインポート

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ３．（８）US 通知例
func TestParseEventUS_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no^B2^Ap_date^B2018.12.03-11:34:51.557^Ap_errno^B0^Ap_err^B^Ap_cmd^BUS^Ap_PV^BMSGSV^Ap_ENO^B5227^Ap_ALT^B0^A" +
		"p_CT^B20181203075545^Ap_MC^B00^Ap_GSCD^B^Ap_SHSB^B^Ap_UC^B01^Ap_UU^B0101^Ap_EDK^B0^Ap_US^B100\n"

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
	assert.Equal(t, "US", event.EventType)
	assert.Equal(t, "5227", event.EventNo) // EventNo (文字列型)
	assert.Equal(t, false, event.IsFirstEvent)
	assert.Equal(t, "MSGSV", event.Provider) // 修正

	// OperationStatus 構造体のフィールドの確認 (event.System が nil でないことを確認してからアクセス)
	if assert.NotNil(t, event.System) {
		assert.Equal(t, "20181203075545", event.System.UpdateDate) //SystemStatusのUpdateDate
		assert.Equal(t, "00", event.System.MarketCode)             // SystemStatus の MarketCode
		assert.Equal(t, "", event.System.UnderlyingAssetCode)      // SystemStatus の UnderlyingAssetCode
		assert.Equal(t, "", event.System.ProductType)              // SystemStatus の ProductType
		assert.Equal(t, "01", event.System.UpdateCategory)         // SystemStatus の UpdateCategory
		assert.Equal(t, "0101", event.System.UpdateUnit)           // SystemStatus の UpdateUnit
		assert.Equal(t, "0", event.System.BusinessDayFlag)         // SystemStatus の BusinessDayFlag
		assert.Equal(t, "100", event.System.UpdateStatus)          // SystemStatus の UpdateStatus
	}
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventUS_Example
