package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"
	// domain パッケージをインポート (必要に応じて)
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ４．使用例 （２）応答
func TestParseEventResponse_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no^B1^Ap_date^B2020.06.18-07:30:34.810^Ap_errno^B0^Ap_err^B^Ap_cmd^BSS^Ap_PV^BMSGSV^Ap_ENO^B3^Ap_ALT^B0^Ap_CT^B20200618052959^Ap_LK^B1^Ap_SS^B1\n" +
		"p_no^B2^Ap_date^B2020.06.18-07:30:34.820^Ap_errno^B0^Ap_err^B^Ap_cmd^BUS^Ap_PV^BMSGSV^Ap_ENO^B7^Ap_ALT^B0^Ap_CT^B20200618064527^Ap_MC^B01^A" +
		"p_GSCD^B101^Ap_SHSB^B03^Ap_UC^B02^Ap_UU^B0201^Ap_EDK^B0^Ap_US^B050\n" +
		"p_no^B3^Ap_date^B2020.06.18-07:30:34.831^Ap_errno^B0^Ap_err^B^Ap_cmd^BUS^Ap_PV^BMSGSV^Ap_ENO^B9^Ap_ALT^B0^Ap_CT^B20200618064531^Ap_MC^B01^A" +
		"p_GSCD^B101^Ap_SHSB^B04^Ap_UC^B02^Ap_UU^B0202^Ap_EDK^B0^Ap_US^B050\n" +
		"p_no^B4^Ap_date^B2020.06.18-07:30:39.842^Ap_errno^B0^Ap_err^B^Ap_cmd^BKP\n" +
		"p_no^B5^Ap_date^B2020.06.18-07:30:44.842^Ap_errno^B0^Ap_err^B^Ap_cmd^BKP\n" +
		"p_no^B1^Ap_date^B2020.06.18-07:30:45.533^Ap_errno^B2^Ap_err^Bsession inactive.^Ap_cmd^BST\n"

	// 制御文字をエスケープシーケンスに置換
	replacer := strings.NewReplacer("^A", "\x01", "^B", "\x02")
	replacedMessage := replacer.Replace(originalMessage)

	// バイト列に変換 (不要、後で削除)
	//message := []byte(replacedMessage)

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	messages := strings.Split(replacedMessage, "\n") // 修正: replacedMessage を使う

	for _, msg := range messages {
		if len(msg) == 0 {
			continue
		}
		event, err := tachibana.CallParseEvent(es, []byte(msg)) // 修正: []byte(msg) を使う
		assert.NoError(t, err)
		assert.NotNil(t, event)

		// 各イベントタイプに応じた追加のアサーション (ここでは簡略化)
		switch event.EventType {
		case "SS":
			assert.Equal(t, "3", event.EventNo)
			if assert.NotNil(t, event.System) {
				assert.Equal(t, "1", event.System.SystemState) //SystemStatusのSystemState
			}
		case "US":
			if event.EventNo == "7" {
				if assert.NotNil(t, event.System) {
					assert.Equal(t, "0201", event.System.UpdateUnit) //SystemStatusのUpdateUnit
				}
			} else if event.EventNo == "9" {
				if assert.NotNil(t, event.System) {
					assert.Equal(t, "0202", event.System.UpdateUnit) //SystemStatusのUpdateUnit
				}
			}
		case "KP":
			// KP イベントのテスト
			assert.Equal(t, "KP", event.EventType) // event.EventType が "KP" であることを確認
		case "ST":
			assert.Equal(t, "ST", event.EventType)
			if assert.NotNil(t, event.System) {
				assert.Equal(t, "2", event.System.ErrNo)                  //SystemStatusのErrNo
				assert.Equal(t, "session inactive.", event.System.ErrMsg) //SystemStatusのErrMsg
			}
		}
	}
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventResponse_Example
