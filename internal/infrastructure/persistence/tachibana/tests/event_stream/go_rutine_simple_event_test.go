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

func TestEventStreamSimpleGoRutine(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("シンプルなイベントストリーム接続テスト", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		eventCh := make(chan *domain.OrderEvent, 10)
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), eventCh)

		// イベントストリーム開始 (goroutine 1)
		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		// イベント処理ルーチン (goroutine 2)
		doneCh := make(chan struct{}) // 終了通知用チャネル
		errCh := make(chan error, 1)  // エラー通知チャネル
		kpCount := 0
		go func() {
			defer close(doneCh) // goroutine終了時にチャネルを閉じる
			for {
				select {
				case event, ok := <-eventCh:
					if !ok { // eventChが閉じられた場合
						return
					}
					fmt.Printf("Received event: %+v\n", event)
					if event.EventType == "KP" {
						kpCount++
						if kpCount >= 2 {
							return // 2回KPを受信したら終了
						}
					}
				case <-time.After(30 * time.Second): // 30秒のタイムアウト
					errCh <- fmt.Errorf("Timeout: KP event not received twice within 30 seconds. Received %d KP events", kpCount)
					return
				}
			}
		}()

		// メインルーチン: イベント処理ルーチンの完了を待つ
		select {
		case <-doneCh: // イベント処理ルーチンが正常終了
			// 成功
		case err := <-errCh: // イベント処理ルーチンでエラー発生
			t.Fatal(err)
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamSimpleGoRutine
