package tachibana_test

import (
	"errors"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEventStream_StartStop(t *testing.T) {
	// httpmockを有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// ロガーの作成 (テスト用のロガー)
	logger, _ := zap.NewDevelopment()

	// イベント受信用チャネルの作成
	eventCh := make(chan domain.OrderEvent, 3)

	// 共通のhttpmock設定 (成功ケース用)
	httpmock.RegisterResponder("GET", "https://example.com/event",
		httpmock.NewStringResponder(200, "p_no^B1^Ap_date^B2024.03.03-10:00:00.000^Ap_cmd^BSS"), //仮のデータ
	)

	// Start メソッドのテスト
	t.Run("Start successfully", func(t *testing.T) {
		// SetupTestClientを使ってモッククライアントと設定を作成
		client, cfg := tachibana.SetupTestClient(t)

		// EventStream の作成
		es := tachibana.NewEventStream(client, cfg, logger, eventCh)

		go func() {
			// EventStream を開始
			es.Start()
		}()

		select {
		case event := <-eventCh:
			t.Log("Received an event:", event)
		case <-time.After(3 * time.Second):
			t.Fatal("Timeout: Event not received")
		}

		// EventStream を停止
		err := es.Stop()
		assert.NoError(t, err)
	})

	// Start メソッドが失敗するケースのテスト (httpmock を活用)
	t.Run("Start failure", func(t *testing.T) {
		// SetupTestClient を使ってモッククライアントと設定を作成
		client, cfg := tachibana.SetupTestClient(t)

		// EventStream の作成
		es := tachibana.NewEventStream(client, cfg, logger, eventCh)

		// モックを設定: ネットワークエラーをシミュレート
		// eventURL取得をモック
		tachibana.SetEventURLForTest(client, "https://example.com/mocked_event") //eventURLをモック
		httpmock.RegisterResponder("GET", "https://example.com/mocked_event",    //モック用URL
			httpmock.NewErrorResponder(errors.New("network error")), // ネットワークエラーを返す
		)

		// Start メソッドをゴルーチンで実行
		startErrCh := make(chan error)
		go func() {
			startErrCh <- es.Start()
		}()

		select {
		case err := <-startErrCh:
			// エラーが発生したことを確認し、エラーメッセージをチェック
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "network error", "Expected a network error")

		case <-time.After(15 * time.Second):
			t.Fatal("Timeout: Start method did not return an error")
		}
	})
}
