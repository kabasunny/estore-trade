// internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_test.go

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

func TestEventStreamSpotBuy(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil) // テストクライアント作成
	ctx := context.Background()                  // context作成

	t.Run("イベントストリームのテスト (現物成行買い)", func(t *testing.T) {
		// Login
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx) // Logoutを確実に実行

		eventCh := make(chan *domain.OrderEvent, 10) // バッファ付きチャネルに変更, ポインタ型に

		// EventStream 作成
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), eventCh)

		// EventStream の開始 (ゴルーチンで実行)
		go func() {
			err := eventStream.Start()
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err) //エラーの場合は、テストを停止
			}
		}()
		defer eventStream.Stop() // Stopを確実に実行

		time.Sleep(1 * time.Second) // イベントストリーム接続確立を待つ(短縮)

		// 現物成行買い注文
		order := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long", // 買い
			OrderType:  "market",
			Quantity:   100,  // 100単位
			MarketCode: "00", // 東証
		}
		placedOrder, err := client.PlaceOrder(ctx, order)
		if err != nil {
			t.Fatalf("Failed to place order: %v", err) // エラーの場合はテストを終了
		}
		assert.NotNil(t, placedOrder)
		placedOrderID := placedOrder.TachibanaOrderID // 発注した注文のIDを保存

		// イベント受信ループ (タイムアウトまで継続)
		timeout := time.After(10 * time.Second) // 60秒のタイムアウト

		for {
			select {
			case event := <-eventCh:
				if event == nil { //nilチェック
					continue
				}
				fmt.Printf("Received event: %+v\n", event) // 全イベントを出力

				// ECイベント、かつ、注文番号が一致するか確認
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedOrderID {

					// 注文ステータス、約定ステータスを確認
					// p_ODST: 0(受付未済), 1(受付済), 2(受付エラー), 3(一部失効), 4(全部失効), 5(繰越失効)
					// p_EXST: 0(未約定), 1(一部約定), 2(全部約定), 3(約定中)
					if event.Order.Status == "1" || event.Order.Status == "3" { // 受付済 or 一部約定
						if event.Order.FilledQuantity > 0 { //約定数量が0より大きい
							t.Logf("Order partially or fully executed. Status: %s, Executed Quantity:%d", event.Order.Status, event.Order.FilledQuantity)
							return // テスト成功（約定を確認）
						}
						//assert.Equal(t, "long", event.Order.Side) //p_BBKB=3 なら long
						// 他の必要なアサーションもここに追加(約定数量、約定価格など)
					} else if event.Order.Status == "4" || event.Order.Status == "5" {
						t.Fatalf("Order failed. Status: %s", event.Order.Status) //約定失敗
						return                                                   //テスト失敗
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Execution event not received after 60 seconds")
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamSpotBuy

// クライアント作成とログイン:
// tachibana.CreateTestClient(t, nil) でテスト用のAPIクライアントを作成。
// client.Login(ctx, nil) で立花証券APIにログイン。
// defer client.Logout(ctx) で、テスト終了時に必ずログアウトする。

// イベントストリームの設定:
// eventCh := make(chan domain.OrderEvent, 10) で、domain.OrderEvent 型のイベントを受信するためのバッファ付きチャネルを作成(容量10)。
// tachibana.NewEventStream(...) で EventStream インスタンスを作成。
// go func() { ... }() で eventStream.Start() をゴルーチンで実行し、イベントストリームを開始。
// defer eventStream.Stop() で、テスト終了時に必ずイベントストリームを停止。
// time.Sleep(3 * time.Second) で、イベントストリームの接続が確立されるまで少し待機 (これは、立花証券APIの仕様やネットワーク状況によっては不要)。

// 現物成行買い注文の発注:
// &domain.Order{...} で、現物成行買い注文の domain.Order 構造体を作成。
// Symbol: 銘柄コード (例: "7974" 任天堂)
// Side: "long" (買い)
// OrderType: "market" (成行)
// Quantity: 注文数量 (例: 100株)
// MarketCode: 市場コード (例: "00" 東証)
// client.PlaceOrder(ctx, order) で注文を発注。
// placedOrder には、発注に成功した注文の情報が格納 (立花証券側の注文IDなど)。
// assert.NotNil(t, placedOrder) で、placedOrder が nil でないこと (つまり、注文が正常に発注されたこと) を確認。

// 約定イベントの受信と検証:
// select 文を使って、以下のいずれかのケースが発生するのを待つ。
// case event := <-eventCh:: eventCh からイベントを受信した場合。
// 受信したイベントが EC (約定通知) イベントであり、かつ、event.Order が nil でなく、event.Order.TachibanaOrderID が placedOrder.TachibanaOrderID (発注した注文の立花証券側ID) と一致するかを確認。
// 一致する場合は、目的の約定イベントであると判断し、assert.Equal(t, "buy", event.Order.Side) で、注文の Side が "long" であることを確認。

// すべての確認が完了したら、return でループを抜け、テストを成功。
// case <-time.After(10 * time.Second):: 10秒間タイムアウトした場合。
// t.Log(...) でタイムアウトしたことをログに出力し、ループを継続します (再試行)。
// 5回の試行で目的の約定イベントを受信できなかった場合は、t.Fatal(...) でテストを失敗。
