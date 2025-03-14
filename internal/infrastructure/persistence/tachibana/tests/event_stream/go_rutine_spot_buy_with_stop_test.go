// internal/infrastructure/persistence/tachibana/tests/event_stream/go_rutine_spot_buy_with_stop_separate_test.go
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

func TestEventStreamSpotBuyWithStopSeparateGoRutine(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (現物成行買い + 逆指値売り - 分割注文)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		eventCh := make(chan *domain.OrderEvent, 100) // バッファを増やす
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), eventCh)

		// チャネルの定義を Run の中に移動
		buyOrderOKCh := make(chan bool, 1)
		stopOrderOKCh := make(chan bool, 1)
		doneCh := make(chan bool)
		errCh := make(chan error, 1)

		// 注文ID と 約定イベント を保持する変数を Run スコープで定義
		var buyOrderID string
		var stopLossOrderID string
		var buyEvent *domain.OrderEvent

		go func() {
			err := eventStream.Start()
			if err != nil {
				errCh <- fmt.Errorf("EventStream Start returned error: %v", err)
				return //goroutineを終了
			}
		}()
		defer eventStream.Stop()

		// イベント処理 (goroutine 2)
		go func() {
			defer func() {
				doneCh <- true // 終了を通知
			}()
			for event := range eventCh {
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)

				if event.EventType == "EC" && event.Order != nil {
					// 注文IDでイベントを振り分け
					if event.Order.TachibanaOrderID == buyOrderID {
						if (event.Order.Status == "1" || event.Order.Status == "3") && event.Order.FilledQuantity > 0 {
							t.Logf("Buy order executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
							buyEvent = event     // 約定情報を保存
							buyOrderOKCh <- true // 約定を通知
							return               // 買い注文の約定イベント受信後、処理を終了
						}
					} else if event.Order.TachibanaOrderID == stopLossOrderID {
						if event.Order.NotificationType == "100" && event.Order.Status == "1" {
							t.Logf("Stop-loss order accepted. Status: %s, NotificationType: %s", event.Order.Status, event.Order.NotificationType)
							stopOrderOKCh <- true // 受付完了を通知
							return                // 逆指値注文の受付イベント受信後、処理を終了
						}
					}
				}
			}
		}()

		time.Sleep(3 * time.Second) // イベントストリーム接続確立を待つ

		// --- 1. 現物成行買い注文 ---
		buyOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: "00", // 東証
		}

		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)
		buyOrderID = placedBuyOrder.TachibanaOrderID // 注文ID を設定

		// --- 2. 買い注文の約定確認 ---
		select {
		case <-buyOrderOKCh:
			t.Log("Buy order confirmed.")
		case err := <-errCh:
			t.Fatalf("Error in event processing: %v", err)
			return
		case <-time.After(60 * time.Second):
			t.Fatal("Timeout: Buy order execution event not received")
			return
		}

		// --- 3. 逆指値売り注文 (ストップロス) ---
		stopLossOrder := &domain.Order{
			Symbol:       buyEvent.Order.Symbol,         // 買い注文と同じ銘柄
			Side:         "short",                       // 売り
			OrderType:    "stop",                        // 逆指値
			Quantity:     buyEvent.Order.FilledQuantity, // 買い注文の約定数量
			MarketCode:   "00",                          // 買い注文と同じ市場
			TriggerPrice: 9000.0,                        // ストップロスのトリガー価格
			Price:        0,                             // トリガー後成行
		}

		placedStopLossOrder, err := client.PlaceOrder(ctx, stopLossOrder)
		if err != nil {
			t.Fatalf("Failed to place stop-loss order: %v", err)
		}
		assert.NotNil(t, placedStopLossOrder)
		stopLossOrderID = placedStopLossOrder.TachibanaOrderID // 注文ID を設定

		// --- 4. 逆指値売り注文の受付確認 ---
		select {
		case <-stopOrderOKCh:
			t.Log("Stop-loss order accepted.")
		case err := <-errCh:
			t.Fatalf("Error in event processing: %v", err)
			return
		case <-time.After(60 * time.Second):
			t.Fatal("Timeout: Stop-loss order was not accepted.")
			return
		}

		// --- 5. すべて完了 ---
		select {
		case <-doneCh:
			t.Log("Event processing goroutine finished.")
		case <-time.After(5 * time.Second): // 5秒待っても終わらなければエラー
			t.Fatal("Timeout: Event processing goroutine did not finish.")
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamSpotBuyWithStopSeparateGoRutine

// 登場人物:
// メインルーチン: テストの全体的な流れを制御。注文の発注、結果の検証、後処理（ログアウトなど）を行う。
// イベントストリーム開始ルーチン (goroutine 1): イベントストリームの開始と停止を担当。
// イベント処理ルーチン (goroutine 2): イベントストリームから受信したイベントを処理し、メインルーチンに結果を通知。

// 処理フロー:
// 準備:
// eventCh: イベントストリームからのイベントを受信するためのチャネル (容量大きめ)。
// buyOrderOKCh: 成行買い注文の約定完了をメインルーチンに通知するためのチャネル (容量1)。
// stopOrderOKCh: 逆指値売り注文の受付完了をメインルーチンに通知するためのチャネル (容量1)。
// doneCh: イベント処理ルーチンの終了を待つためのチャネル (容量1)。
// errCh: エラーをメインルーチンに通知するためのチャネル。
// イベントストリーム開始 (goroutine 1):
// tachibana.NewEventStream() で EventStream インスタンスを作成。
// eventStream.Start() を呼び出して、イベントストリームを開始。
// エラーが発生した場合は、errCh にエラーを送信。
// ゴルーチン終了時に eventStream.Stop() を呼び出して、イベントストリームを停止。
// イベント処理開始 (goroutine 2):
// eventCh からのイベント受信を待ち続ける無限ループ。
// 受信したイベントの種類 (EventType)、注文ID (TachibanaOrderID)、ステータス (Status)、通知種別 (NotificationType) などをチェック。
// 成行買い注文の約定イベント (EC, Status が "1" or "3", FilledQuantity > 0): buyOrderOKCh に true を送信 (または、約定情報を送信)。
// 逆指値売り注文の受付イベント (EC, NotificationType が "100", Status が "1"): stopOrderOKCh に true を送信。
// エラーが発生した場合: errCh にエラーを送信。
// ループの最後に doneCh に true を送信（ゴルーチンが終了したことを通知）。
// メインルーチン:
// goroutine 1, 2 を起動。
// 成行買い注文を発注。
// buyOrderOKCh からの受信を待つ (買い注文の約定完了を待つ)。
// 逆指値売り注文を発注。
// stopOrderOKCh からの受信を待つ (逆指値注文の受付完了を待つ)。
// errCh からのエラー受信を監視。
// doneChからの受信を待つ
// 後処理 (ログアウトなど)。
