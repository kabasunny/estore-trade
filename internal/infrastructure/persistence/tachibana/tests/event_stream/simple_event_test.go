// internal/infrastructure/persistence/tachibana/tests/event_stream/simple_event_test.go
package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher パッケージをインポート
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestEventStreamSimple(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("シンプルなイベントストリーム接続テスト", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// OrderEventDispatcher の作成
		dispatcher := dispatcher.NewOrderEventDispatcher(client.GetLogger())

		// EventStream 作成 (dispatcher を渡す)
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), dispatcher)

		go func() {
			err := eventStream.Start(ctx) //contextを渡す
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		// イベント受信用チャネルを作成
		eventCh := make(chan *domain.OrderEvent, 10)
		// dispatcherに登録
		dispatcher.Subscribe("system", eventCh) // "system" 購読者に変更
		defer dispatcher.Unsubscribe("system", eventCh)

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
				t.Fatalf("Timeout: KP event not received twice within 60 seconds. Received %d KP events", kpCount)
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/simple_event_test.go
