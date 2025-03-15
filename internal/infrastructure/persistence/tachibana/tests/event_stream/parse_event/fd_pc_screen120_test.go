// fd_pc_screen120_test.go
package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ３．（３）FD 通知例 株価ボード・ＰＣ画面１２０
func TestParseEventFD_PCScreen120_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no=1;p_date=2017.05.02-14:13:56.914;p_errno=0;p_err=;p_cmd=FD;p_1001_DPP=6126;p_1001_DPG=0058;p_1001_DYWP=84;p_1001_DV=3899600;" +
		"p_1002_DPP=6533;p_1002_DPG=0058;p_1002_DYWP=93;p_1002_DV=822100;p_1003_DPP=28255;p_1003_DPG=0058;" +
		"p_1003_DYWP=-60;p_1003_DV=1193700;p_1004_DPP=624.9;p_1004_DPG=0058;p_1004_DYWP=1.3;p_1004_DV=8824000;p_1005_DPP=1138.5;" +
		"p_1005_DPG=0058;p_1005_DYWP=5.5;p_1005_DV=4440600;p_1006_DPP=8475;p_1006_DPG=0058;p_1006_DYWP=18;p_1006_DV=2790000;" +
		"p_1007_DPP=206.9;p_1007_DPG=0058;p_1007_DYWP=1.7;p_1007_DV=67564100;p_1008_DPP=716.1;p_1008_DPG=0057;p_1008_DYWP=4.4;" +
		"p_1008_DV=35136000;p_1009_DPP=4155;p_1009_DPG=0058;p_1009_DYWP=32;p_1009_DV=3247600;p_1010_DPP=3746;p_1010_DPG=0058;" +
		"p_1010_DYWP=23;p_1010_DV=3707000;p_1011_DPP=2195.5;p_1011_DPG=0058;p_1011_DYWP=51.0;p_1011_DV=2512100;p_1012_DPP=15040;" +
		"p_1012_DPG=0057;p_1012_DYWP=190;p_1012_DV=5659957;p_1013_DPP=884;p_1013_DPG=0057;p_1013_DYWP=12;p_1013_DV=555200;" +
		"p_1014_DPP=1722.5;p_1014_DPG=0058;p_1014_DYWP=20.0;p_1014_DV=3054900;p_1015_DPP=1618.0;p_1015_DPG=0058;p_1015_DYWP=21.5;" +
		"p_1015_DV=9857100;p_1016_DPP=1075.0;p_1016_DPG=0057;p_1016_DYWP=10.5;p_1016_DV=5967600;p_1017_DPP=3237;p_1017_DPG=0058;" +
		"p_1017_DYWP=-8;p_1017_DV=3944300;p_1018_DPP=2938.5;p_1018_DPG=0058;p_1018_DYWP=-22.5;p_1018_DV=2575000;" +
		"p_1019_DPP=503.1;p_1019_DPG=0058;p_1019_DYWP=8.5;p_1019_DV=8188800;p_1020_DPP=19425;p_1020_DPG=0058;p_1020_DYWP=115;" +
		"p_1020_DV=20446;p_2001_DPP=1677;p_2001_DPG=0057;p_2001_DYWP=9;p_2001_DV=82000;p_2002_DPP=1697;p_2002_DPG=0058;" +
		"p_2002_DYWP=8;p_2002_DV=283400;p_2003_DPP=3905;p_2003_DPG=0057;p_2003_DYWP=25;p_2003_DV=700;p_2004_DPP=602;p_2004_DPG=0058;" +
		"p_2004_DYWP=5;p_2004_DV=68000;p_2005_DPP=361;p_2005_DPG=0000;p_2005_DYWP=0;p_2005_DV=1000;....\n"

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
	// 必要に応じて、Market 構造体の内容をアサーション
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventFD_PCScreen120_Example
