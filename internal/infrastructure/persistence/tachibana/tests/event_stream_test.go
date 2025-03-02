// internal/infrastructure/persistence/tachibana/tests/event_stream_test.go
package tachibana_test

import (
	"testing"
	"time"

	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// internal/infrastructure/persistence/tachibana/tests/event_stream_test.go
func TestEventStream_StartStop(t *testing.T) {
	// ... (httpmock の設定など、他の部分は変更なし) ...
	// ロガーの作成 (テスト用のロガー)
	logger, _ := zap.NewDevelopment()

	// イベント受信用チャネルの作成
	//var eventCh chan domain.OrderEvent = make(chan domain.OrderEvent, 10) // 型を明示
	eventCh := make(chan domain.OrderEvent, 3)

	//共通のhttpmock設定
	httpmock.RegisterResponder("GET", "https://example.com/event",
		httpmock.NewStringResponder(200, "p_no^B1^Ap_date^B2024.03.03-10:00:00.000^Ap_cmd^BSS"), //仮のデータ
	)

	// Start メソッドのテスト
	t.Run("Start successfully", func(t *testing.T) {
		// SetupTestClientを使ってモッククライアントと設定を作成
		client, _ := tachibana.SetupTestClient(t) //configは使わない

		// テスト用の設定 (テストに必要な最小限の設定)
		cfg := &config.Config{
			EventRid:          "test_rid",
			EventBoardNo:      "test_board_no",
			EventEvtCmd:       "EC,SS",                //SSも必要
			TachibanaBaseURL:  "https://example.com/", // モックするのでダミーでOK
			TachibanaUserID:   "testuser",
			TachibanaPassword: "testpassword",
		}
		// EventStream の作成, http.DefaultClientを渡す
		es := tachibana.NewEventStream(client, cfg, logger, eventCh)

		// Start() エラー用のチャネル
		// startErrCh := make(chan error)

		// EventStream をゴルーチンで開始
		go func() {
			es.Start()
			// err := es.Start()
			// if err != nil {
			// 	startErrCh <- err // Start() が失敗したらエラーをチャネルに送信
			// 	return
			// }
			// startErrCh <- nil // Start() が成功したら nil を送信 (完了を通知)

		}()

		// go func() {
		// 	select {
		// 	case startErr := <-startErrCh:
		// 		if startErr != nil {
		// 			t.Errorf("EventStream Start が失敗しました: %v", startErr) // t.Error または t.Errorf を使用
		// 			return                                              // エラーをログ出力したらゴルーチンを終了
		// 		}
		// 	case <-time.After(5 * time.Second):
		// 		t.Errorf("EventStream Start の完了を待機中にタイムアウト") // t.Error または t.Errorf を使用
		// 		return                                       // タイムアウトをログ出力したらゴルーチンを終了
		// 	}
		// 	// エラーがなければ、テスト成功を通知するチャネルに送信するなど、必要に応じて処理を追加
		// }()

		select { // 元の select 文 (イベント受信を待つ)
		case event := <-eventCh:
			t.Log("Received an event:", event)
		case <-time.After(3 * time.Second):
			t.Fatal("Timeout: Event not received")
		}

		// EventStream を停止
		err := es.Stop()
		assert.NoError(t, err)
	})

	// Login に失敗するケース
	// t.Run("Login fails", func(t *testing.T) {
	// 	// SetupTestClientを使ってモッククライアントと設定を作成(Login失敗させる)
	// 	client, _ := tachibana.SetupTestClient(t)
	// 	// テスト用の設定 (テストに必要な最小限の設定)
	// 	cfg := &config.Config{
	// 		EventRid:          "test_rid",
	// 		EventBoardNo:      "test_board_no",
	// 		EventEvtCmd:       "EC,SS",
	// 		TachibanaBaseURL:  "https://example.com/", // モックするのでダミーでOK
	// 		TachibanaUserID:   "wronguser",            // ★誤ったユーザーID★
	// 		TachibanaPassword: "testpassword",
	// 	}
	// 	// EventStream の作成, http.DefaultClientを渡す
	// 	es := tachibana.NewEventStream(client, cfg, logger, eventCh)

	// 	err := es.Start()
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "failed to login for event stream")

	// 	// Stop を呼び出してもエラーにならないことを確認
	// 	err = es.Stop() //Loginに失敗しているので、Stopは呼ばれない
	// 	assert.NoError(t, err)
	// })
}

// func TestEventStream_parseEvent(t *testing.T) {
// 	// ... (TestEventStream_parseEvent の内容は変更なし) ...
// 	logger := zaptest.NewLogger(t) // テスト用のロガー
// 	cfg := &config.Config{}
// 	mockClient := &tachibana.MockTachibanaClient{} //MockClientを使う
// 	eventCh := make(chan domain.OrderEvent)
// 	es := tachibana.NewEventStream(mockClient, cfg, logger, eventCh)

// 	t.Run("Valid EC event", func(t *testing.T) {
// 		message := []byte("p_cmd^BEC^Ap_date^B2023.11.15-10:00:00.000^Ap_ON^B12345^Ap_IC^B7203^Ap_BBKB^B3^Ap_ODST^B1^Ap_CRPR^B1500^Ap_CRSR^B100")
// 		event, err := tachibana.ParseEventForTest(es, message)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, event)
// 		assert.Equal(t, "EC", event.EventType)
// 		assert.Equal(t, "12345", event.Order.ID)
// 		assert.Equal(t, "7203", event.Order.Symbol)
// 		assert.Equal(t, "buy", event.Order.Side)
// 		assert.Equal(t, "1", event.Order.Status)
// 		assert.Equal(t, 1500.0, event.Order.Price)
// 		assert.Equal(t, 100, event.Order.Quantity)
// 	})

// 	t.Run("Valid NS event", func(t *testing.T) {
// 		message := []byte("p_cmd^BNS^Ap_date^B2024.03.02-13:04:05.000")
// 		event, err := tachibana.ParseEventForTest(es, message)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, event)
// 		assert.Equal(t, "NS", event.EventType)
// 		assert.Nil(t, event.Order)
// 	})

// 	t.Run("Empty event type", func(t *testing.T) {
// 		message := []byte("p_date^B2023.11.15-10:00:00.000")
// 		_, err := tachibana.ParseEventForTest(es, message)
// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "event type is empty: p_date^B2023.11.15-10:00:00.000")
// 	})

// 	t.Run("Invalid date format", func(t *testing.T) {
// 		message := []byte("p_cmd^BEC^Ap_date^Binvalid-date")
// 		_, err := tachibana.ParseEventForTest(es, message)
// 		assert.Error(t, err) // エラーは発生するが、処理は続行する
// 	})

// 	t.Run("Multiple market codes", func(t *testing.T) {
// 		message := []byte("p_cmd^BEC^Ap_date^B2023.11.15-10:00:00.000^Ap_ON^B12345^Ap_IC^B7203^Ap_MC^B00^C01^Ap_BBKB^B3")
// 		event, err := tachibana.ParseEventForTest(es, message)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, event)
// 		assert.Equal(t, "EC", event.EventType)
// 		assert.Equal(t, "12345", event.Order.ID)
// 		assert.Equal(t, "7203", event.Order.Symbol)
// 		assert.Equal(t, "00", event.Order.MarketCode) //現状のコードでは最初の値だけ
// 	})

// 	t.Run("Malformed message", func(t *testing.T) {
// 		message := []byte("p_cmd^BEC^Ap_date^B2023.11.15-10:00:00.000^A^Ap_ON^B12345") // 不正な形式
// 		event, err := tachibana.ParseEventForTest(es, message)
// 		assert.NoError(t, err)  //エラーは返ってこない
// 		assert.NotNil(t, event) // 必要な値が設定されていなくても、とりあえずeventは作成される
// 	})

// 	t.Run("Long message", func(t *testing.T) {
// 		longMessage := "p_cmd^BEC^Ap_date^B2024.03.04-14:30:00.000^Ap_ON^B11111^Ap_IC^B7203^Ap_BBKB^B3^Ap_ODST^B0^Ap_CRPR^B6500"
// 		for i := 0; i < 100; i++ { // 100個のダミーフィールドを追加
// 			longMessage += fmt.Sprintf("^Ap_dummy%d^Bvalue%d", i, i)
// 		}
// 		event, err := tachibana.ParseEventForTest(es, []byte(longMessage))
// 		assert.NoError(t, err)
// 		assert.NotNil(t, event)
// 		assert.Equal(t, "EC", event.EventType)
// 		assert.Equal(t, "11111", event.Order.ID)    // 注文番号
// 		assert.Equal(t, "7203", event.Order.Symbol) // 銘柄コード
// 		assert.Equal(t, "buy", event.Order.Side)    // 売買区分
// 		assert.Equal(t, "0", event.Order.Status)    // 注文ステータス
// 		assert.Equal(t, 6500.0, event.Order.Price)  // 価格
// 	})
// }

// func TestEventStream_sendEvent(t *testing.T) {
// 	logger := zaptest.NewLogger(t) // テスト用のロガー
// 	cfg := &config.Config{}
// 	mockClient := &tachibana.MockTachibanaClient{}
// 	eventCh := make(chan domain.OrderEvent)

// 	es := tachibana.NewEventStream(mockClient, cfg, logger, eventCh)

// 	t.Run("Send event successfully", func(t *testing.T) {
// 		event := domain.OrderEvent{EventType: "test"}

// 		// ゴルーチンでイベントを受信
// 		go func() {
// 			receivedEvent := <-eventCh
// 			assert.Equal(t, event, receivedEvent)
// 		}()

// 		tachibana.SendEventForTest(es, &event) //ヘルパー関数を使用
// 	})

// 	t.Run("Event channel full", func(t *testing.T) {
// 		// eventCh はバッファなしチャネルなので、すぐにフルになる
// 		event1 := domain.OrderEvent{EventType: "test1"}
// 		event2 := domain.OrderEvent{EventType: "test2"}

// 		tachibana.SendEventForTest(es, &event1) // ヘルパー関数

// 		// ゴルーチンで sendEvent を呼び出し、ブロックされることを確認
// 		done := make(chan bool)
// 		go func() {
// 			tachibana.SendEventForTest(es, &event2) // ヘルパー関数
// 			close(done)
// 		}()

// 		// 少し待って、まだブロックされていることを確認
// 		time.Sleep(100 * time.Millisecond)
// 		select {
// 		case <-done:
// 			t.Fatal("sendEvent should have been blocked")
// 		default:
// 			// ブロックされている (期待通り)
// 		}
// 	})
// 	t.Run("Stop signal received", func(t *testing.T) {
// 		event := domain.OrderEvent{EventType: "test"}
// 		es.Stop()                              // EventStream を停止
// 		tachibana.SendEventForTest(es, &event) // ヘルパー関数 //イベント送信、停止しているので、送信されない
// 	})
// }
