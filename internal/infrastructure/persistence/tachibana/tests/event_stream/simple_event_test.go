package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestEventStreamSimple(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil) //tachibanaClient -> *Client
	ctx := context.Background()

	t.Run("シンプルなイベントストリーム接続テスト", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		eventCh := make(chan *domain.OrderEvent, 10) //型を変更
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), eventCh)

		go func() {
			// err := eventStream.Start() //Start()ではなく、StartSample()
			err := eventStream.Start(ctx) //contextを渡す
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		// KPメッセージを2回受信するまで待機 (最大30秒)
		timeout := time.After(60 * time.Second)
		kpCount := 0
		for {
			select {
			case event := <-eventCh:
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "KP" {
					kpCount++
					if kpCount >= 2 {
						return // 2回KPを受信したらテスト成功
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: KP event not received twice within 30 seconds. Received %d KP events", kpCount)
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamSimple

// ログイン: client.Login で立花証券APIにログイン

// イベントストリーム開始:
// tachibana.NewEventStream で EventStream インスタンスを作成。
// go func() { ... }() 内で eventStream.Start() を呼び出し、イベントストリームを開始 (ゴルーチンで実行)。
// Start メソッドは、client.GetEventURL() で取得したベースURLと、設定値からクエリパラメータを組み立てて、イベントストリーム用のURLを作成。
// HTTP GET リクエストを送信し、レスポンスステータスコードが200 (OK) であることを確認。
// processResponseBody メソッドでレスポンスボディを処理。

// イベント受信 (processResponseBody):
// processResponseBody は、レスポンスボディを bufio.NewReader で読み込み、ReadBytes('\n') で行ごとに処理。
// 受信したメッセージをログに出力 (16進数ダンプ、Shift-JIS、UTF-8)。
// transform.Bytes で Shift-JIS から UTF-8 に変換。
// es.parseEvent でメッセージを解析し、domain.OrderEvent 構造体を作成。
// es.eventCh <- event で、解析されたイベントをチャネルに送信 (ポインタを送信)。

// テストコード (受信側):
// eventCh := make(chan *domain.OrderEvent, 10) で、domain.OrderEvent のポインタ型を受信するチャネルを作成。
// select 文を使って、eventCh からのイベント受信、またはタイムアウトを待つ。
// 受信したイベントの EventType が "KP" であれば、kpCount をインクリメント。
// kpCount が 2 になったら (KPメッセージを2回受信したら)、テスト成功として return。
// 30秒以内にKPメッセージが2回受信されなければ、テスト失敗。

// イベントストリーム停止/ログアウト:
// defer eventStream.Stop() で、テスト終了時に EventStream.Stop() を呼び出し、イベントストリームを停止 (チャネルをクローズ)。
// defer client.Logout(ctx) で、立花証券APIからログアウト。
