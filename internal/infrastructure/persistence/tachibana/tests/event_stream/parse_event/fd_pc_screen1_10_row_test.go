package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ３．（３）FD 通知例 株価ボード・ＰＣ画面１～１０、行指定板[1]
func TestParseEventFD_PCScreen1_10_Row_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no=1;p_date=2016.11.30-10:18:32.380;p_errno=0;p_err=;p_cmd=FD;p_1_DPP=18340;t_1_DPP:T=10:18;p_1_DPG=0058;p_1_DYWP=50;" +
		"p_1_DYRP=0.27;p_1_DOP=18350;p_1_DHP=18450;p_1_DLP=18290;p_1_DV=31907;p_1_QAS=0101;p_1_QBS=0101;p_1_GAV1=157;" +
		"p_1_GAP1=18350;p_1_GAV2=187;p_1_GAP2=18360;p_1_GAV3=255;p_1_GAP3=18370;p_1_GAV4=268;p_1_GAP4=18380;p_1_GAV5=196;" +
		"p_1_GAP5=18390;p_1_GAV6=219;p_1_GAP6=18400;p_1_GAV7=191;p_1_GAP7=18410;p_1_GAV8=235;p_1_GAP8=18420;p_1_GBV1=153;" +
		"p_1_GBP1=18340;p_1_GBV2=259;p_1_GBP2=18330;p_1_GBV3=232;p_1_GBP3=18320;p_1_GBV4=288;p_1_GBP4=18310;p_1_GBV5=294;" +
		"p_1_GBP5=18300;p_1_GBV6=414;p_1_GBP6=18290;p_1_GBV7=285;p_1_GBP7=18280;p_1_GBV8=282;p_1_GBP8=18270\n"

	// 制御文字と区切り文字を置換
	replacer := strings.NewReplacer(
		"^A", "\x01",
		"^B", "\x02",
		";", "\x01", // ";" を ^A (\x01) に置換
		"=", "\x02", // "=" を ^B (\x02) に置換
	)
	replacedMessage := replacer.Replace(originalMessage)

	message := []byte(replacedMessage)

	logger, _ := zap.NewDevelopment()
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "FD", event.EventType)
	// 他の項目も必要に応じてアサーションする。Exampleのため、一部項目の確認にとどめる
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventFD_PCScreen1_10_Row_Example
