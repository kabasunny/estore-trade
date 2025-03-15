package parse_event_test

import (
	"strings"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"
	// domain パッケージをインポート

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ３．（５）NS 通知例
func TestParseEventNS_Example(t *testing.T) {
	// 元のメッセージ (人間が読める形式)
	originalMessage := "p_no^B7^Ap_date^B2020.08.26-12:59:13.598^Ap_errno^B0^Ap_err^B^Ap_cmd^BNS^Ap_PV^BQNSD^Ap_ENO^B166^Ap_ALT^B0^A" +
		"p_ID^B20200826125300_MIO1708^Ap_DT^B20200826^Ap_TM^B125300^Ap_CGN^B1^Ap_CGL^B100^Ap_GRN^B1^Ap_GRL^B3009^A" +
		"p_ISN^B11^Ap_ISL^B4519^C4568^C4661^C6594^C6758^C6861^C7974^C8301^C9437^C9983^C9984^Ap_SKF^B0^A" +
		"p_UPF^B^Ap_HDL^B<NQN>◇東証後場寄り　下げ幅やや拡大、・・・・\n"

	// 制御文字と区切り文字を置換
	replacer := strings.NewReplacer(
		"^A", "\x01",
		"^B", "\x02",
		"^C", "\x03", // 値と値の区切り文字
	)
	replacedMessage := replacer.Replace(originalMessage)
	message := []byte(replacedMessage)

	logger, err := zap.NewDevelopment() //errを捕捉
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	es := tachibana.NewTestEventStream(logger)

	event, err := tachibana.CallParseEvent(es, message)
	//fmt.Printf("%+v\n", event) // デバッグ用: event の内容を出力
	assert.NoError(t, err)
	assert.NotNil(t, event)

	// OrderEvent のフィールドの確認
	assert.Equal(t, "NS", event.EventType)
	assert.Equal(t, "166", event.EventNo) // EventNo (文字列型)
	assert.Equal(t, false, event.IsFirstEvent)
	assert.Equal(t, "QNSD", event.Provider) //修正

	// News 構造体のフィールドの確認 (event.News が nil でないことを確認してからアクセス)
	if assert.NotNil(t, event.News) {
		assert.Equal(t, "20200826125300_MIO1708", event.News.NewsID)
		assert.Equal(t, "20200826", event.News.NewsDate)
		assert.Equal(t, "125300", event.News.NewsTime)
		assert.Equal(t, 1, event.News.NewsCategoryCount)
		assert.Equal(t, []string{"100"}, event.News.CategoryList)
		assert.Equal(t, 11, event.News.RelatedSymbolCount)
		assert.Equal(t, []string{"4519", "4568", "4661", "6594", "6758", "6861", "7974", "8301", "9437", "9983", "9984"}, event.News.Symbols) //修正不要
		assert.Equal(t, "<NQN>◇東証後場寄り　下げ幅やや拡大、・・・・", event.News.Title)
		assert.Equal(t, []string{"3009"}, event.News.GenreList)

	}
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/parse_event -run TestParseEventNS_Example
