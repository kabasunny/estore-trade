package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ３．（３）FD 通知例 株価ボード・ＰＣ画面１
func TestParseEventFD_PCScreen1_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no=1;p_date=2017.05.02-14:10:05.372;p_errno=0;p_err=;p_cmd=FD;x_1_LISS=82509594;t_1_DPP:T=14:10;p_1_QBS=0101;p_1_QAS=0101;" +
		"p_1_DPP=6128;p_1_DPG=0058;p_1_DYWP=86;p_1_DYRP=1.42;p_1_DOP=6057;p_1_DHP=6138;p_1_DLP=6048;p_1_DV=3864200;p_1_AV=4800;" +
		"p_1_QAP=6129;p_1_BV=1900;p_1_QBP=6128;p_1_DHF=0000;p_1_DLF=0000;p_1_DJ=23627822400;x_2_LISS=82509594;t_2_DPP:T=14:08;" +
		"p_2_QBS=0101;p_2_QAS=0101;p_2_DPP=6533;p_2_DPG=0058;p_2_DYWP=93;p_2_DYRP=1.44;p_2_DOP=6470;p_2_DHP=6586;p_2_DLP=6470;" +
		"p_2_DV=815700;p_2_AV=1300;p_2_QAP=6533;p_2_BV=600;p_2_QBP=6532;p_2_DHF=0000;p_2_DLF=0000;p_2_DJ=5336652800;" +
		"x_3_LISS=82509594;t_3_DPP:T=14:09;p_3_QBS=0101;p_3_QAS=0101;p_3_DPP=28245;p_3_DPG=0058;p_3_DYWP=-70;p_3_DYRP=-0.24;" +
		"p_3_DOP=28300;....\n"

	// 制御文字と区切り文字を置換
	replacer := strings.NewReplacer(
		"^A", "\x01",
		"^B", "\x02",
		";", "\x01", // ";" を ^A (\x01) に置換
		"=", "\x02", // "=" を ^B (\x02) に置換
	)
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
	assert.Equal(t, "FD", event.EventType)
	// 他の項目も必要に応じてアサーションする。Exampleのため、一部項目の確認にとどめる
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventFD_PCScreen1_Example
