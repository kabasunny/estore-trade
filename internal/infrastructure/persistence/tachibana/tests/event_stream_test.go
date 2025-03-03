// internal/infrastructure/persistence/tachibana/tests/event_stream_test.go
package tachibana_test

import (
	"errors"
	"net/http"
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

	// イベント受信用チャネルの作成（使用しないが、EventStreamの作成に必要）
	eventCh := make(chan domain.OrderEvent, 3)

	// Start メソッドのテスト
	t.Run("Start and Stop successfully", func(t *testing.T) {
		// SetupTestClientを使ってモッククライアントと設定を作成
		client, cfg := tachibana.SetupTestClient(t)

		// EventStream の作成
		es := tachibana.NewEventStream(client, cfg, logger, eventCh)

		// Start メソッドをゴルーチンで実行
		startErrCh := make(chan error)
		go func() {
			startErrCh <- es.Start()
		}()

		// 少し待ってから Stop() を呼び出す (Start() がある程度実行される時間を確保)
		time.Sleep(100 * time.Millisecond)

		// EventStream を停止
		err := es.Stop()
		assert.NoError(t, err)

		// Start() がエラーを返していないことを確認
		select {
		case err := <-startErrCh:
			assert.NoError(t, err) // Start() がエラーを返したらテスト失敗
		default:
			// Start() がまだ終了していなければOK (Stop() で停止されたはず)
		}
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

		case <-time.After(15 * time.Second): // タイムアウト時間を短縮
			t.Fatal("Timeout: Start method did not return an error")
		}
	})
}

// TestParseEvent は parseEvent メソッドの単体テスト
func TestParseEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment() // テスト用ロガー

	tests := []struct {
		name        string
		message     []byte
		expected    *domain.OrderEvent
		expectedErr error
	}{
		{
			name:    "Valid EC event",
			message: []byte("p_no^B1^Ap_date^B2024.03.08-15:30:00.000^Ap_cmd^BEC^Ap_ENO^B123^Ap_ON^B456^Ap_IC^B7974^Ap_BBKB^B3^Ap_ODST^B1^Ap_CRPR^B7000^Ap_CRSR^B1"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 8, 15, 30, 0, 0, time.UTC),
				EventType: "EC",
				EventNo:   123,
				Order: &domain.Order{
					ID:       "456",
					Symbol:   "7974",
					Side:     "buy",
					Status:   "1", // Tachibana Securitiesのステータスコード
					Price:    7000,
					Quantity: 1,
				},
			},
			expectedErr: nil,
		},
		{
			name:    "Valid SS event",
			message: []byte("p_no^B1^Ap_date^B2024.03.08-16:00:00.000^Ap_cmd^BSS^Ap_ENO^B456"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 8, 16, 0, 0, 0, time.UTC),
				EventType: "SS",
				EventNo:   456,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid FD event",
			message: []byte("p_no^B1^Ap_date^B2024.03.09-10:00:00.000^Ap_cmd^BFD^Ap_ENO^B789^Ap_ZBI^B7050^Ap_DK^B10000"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 9, 10, 0, 0, 0, time.UTC),
				EventType: "FD",
				EventNo:   789,
				// Order フィールドは nil (FD イベントには Order 情報はない)
			},
			expectedErr: nil,
		},
		{
			name:    "Valid NS event",
			message: []byte("p_no^B1^Ap_date^B2024.03.09-11:00:00.000^Ap_cmd^BNS^Ap_ENO^B1011^Ap_NC^B12345^Ap_ND^B20240309110000"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 9, 11, 0, 0, 0, time.UTC),
				EventType: "NS",
				EventNo:   1011,
				// Order フィールドは nil (NS イベントには Order 情報はない)
			},
			expectedErr: nil,
		},
		{
			name:    "Valid ST event",
			message: []byte("p_no^B1^Ap_date^B2024.03.09-12:00:00.000^Ap_cmd^BST^Ap_ENO^B1314"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 9, 12, 0, 0, 0, time.UTC), //仮
				EventType: "ST",
				EventNo:   1314,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid KP event",
			message: []byte("p_no^B1^Ap_date^B2024.03.09-13:00:00.000^Ap_cmd^BKP^Ap_ENO^B1516"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 9, 13, 0, 0, 0, time.UTC), //仮
				EventType: "KP",
				EventNo:   1516,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid US event",
			message: []byte("p_no^B1^Ap_date^B2024.03.09-14:00:00.000^Ap_cmd^BUS^Ap_ENO^B1718"),
			expected: &domain.OrderEvent{
				Timestamp: time.Date(2024, 3, 9, 14, 0, 0, 0, time.UTC), //仮
				EventType: "US",
				EventNo:   1718,
			},
			expectedErr: nil,
		},
		{
			name:        "Invalid event (missing p_cmd)",
			message:     []byte("p_no^B1^Ap_date^B2024.03.03-10:00:00.000"),
			expected:    nil,
			expectedErr: errors.New("event type is empty"),
		},
		// 他のテストケース (異なるイベントタイプ、エラーケースなど) を追加
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// EventStream インスタンスを作成
			es := &tachibana.EventStream{} // Logger は直接設定しない
			// Logger を設定
			tachibana.SetEventStreamLogger(es, logger)

			// ヘルパー関数を使用して parseEvent を呼び出す
			event, err := tachibana.ParseEventForTest(es, tt.message)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				// 期待されるイベントと、実際にパースされたイベントを比較 (DeepEqual を使用)
				assert.Equal(t, tt.expected, event)
			}
		})
	}
}

// EventStream.Start() を経由したイベント受信とパースのテスト
func TestEventStream_Start_ReceiveAndParse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	logger, _ := zap.NewDevelopment()
	eventCh := make(chan domain.OrderEvent, 10) // バッファ付きチャネル

	// SetupTestClient を使ってモッククライアントと設定を作成
	client, cfg := tachibana.SetupTestClient(t)
	// EventStream の作成
	es := tachibana.NewEventStream(client, cfg, logger, eventCh)

	// 個別のモックレスポンスを設定
	mockResponseEC1 := `p_no^B1^Ap_date^B2024.03.08-15:30:00.000^Ap_cmd^BEC^Ap_ENO^B1^Ap_ON^B1001^Ap_IC^B7974^Ap_BBKB^B3^Ap_ODST^B1^Ap_CRPR^B7000^Ap_CRSR^B1`
	mockResponseSS := `p_no^B2^Ap_date^B2024.03.08-15:31:00.000^Ap_cmd^BSS^Ap_ENO^B2`
	mockResponseEC2 := `p_no^B3^Ap_date^B2024.03.08-15:32:00.000^Ap_cmd^BEC^Ap_ENO^B3^Ap_ON^B1002^Ap_IC^B9984^Ap_BBKB^B1^Ap_ODST^B2^Ap_CRPR^B2500^Ap_CRSR^B2`

	// カスタムレスポンダ関数
	requestCount := 0
	responder := func(req *http.Request) (*http.Response, error) {
		requestCount++
		switch requestCount {
		case 1:
			return httpmock.NewStringResponse(200, mockResponseEC1), nil
		case 2:
			return httpmock.NewStringResponse(200, mockResponseSS), nil
		case 3:
			return httpmock.NewStringResponse(200, mockResponseEC2), nil
		default:
			return httpmock.NewStringResponse(404, "Not Found"), nil // 4回目以降は404
		}
	}

	//eventURL取得をモック
	tachibana.SetEventURLForTest(client, "https://example.com/mocked_event") //eventURLをモック

	// レスポンダを登録
	httpmock.RegisterResponder("GET", "https://example.com/mocked_event", responder)

	// EventStream を開始 (ゴルーチンで実行)
	go es.Start()

	// eventCh からイベントを受信し、内容を確認
	// 期待されるイベント (EC イベント)
	expectedEC1 := domain.OrderEvent{
		Timestamp: time.Date(2024, 3, 8, 15, 30, 0, 0, time.UTC),
		EventType: "EC",
		EventNo:   1,
		Order: &domain.Order{
			ID:       "1001",
			Symbol:   "7974",
			Side:     "buy",
			Status:   "1",
			Price:    7000,
			Quantity: 1,
		},
	}
	expectedSS := domain.OrderEvent{
		Timestamp: time.Date(2024, 3, 8, 15, 31, 0, 0, time.UTC),
		EventType: "SS",
		EventNo:   2,
	}
	expectedEC2 := domain.OrderEvent{
		Timestamp: time.Date(2024, 3, 8, 15, 32, 0, 0, time.UTC),
		EventType: "EC",
		EventNo:   3,
		Order: &domain.Order{
			ID:       "1002",
			Symbol:   "9984",
			Side:     "sell",
			Status:   "2",
			Price:    2500,
			Quantity: 2,
		},
	}

	// イベントを3つ受信するまで待つ (最大5秒)
	receivedEvents := make(map[int]domain.OrderEvent)
	timeout := time.After(5 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case event := <-eventCh:
			receivedEvents[event.EventNo] = event // EventNo をキーとしてマップに格納
		case <-timeout:
			t.Fatalf("Timeout: Expected 3 events, but received only %d", len(receivedEvents))
			return
		}
	}

	// 受信したイベントを検証
	assert.Equal(t, expectedEC1, receivedEvents[1], "Event 1 (EC1) does not match")
	assert.Equal(t, expectedSS, receivedEvents[2], "Event 2 (SS) does not match")
	assert.Equal(t, expectedEC2, receivedEvents[3], "Event 3 (EC2) does not match")

	// EventStream を停止
	es.Stop()

	// リクエスト回数を検証 (オプション)
	info := httpmock.GetCallCountInfo()
	// カウント情報を取得: URLとメソッドを指定して、そのリクエストが何回呼ばれたか
	count := info["GET https://example.com/mocked_event"]
	assert.Equal(t, 3, count, "Expected 3 HTTP requests, but got %d", count) // 3回のつもりが増えてしまったりしていないかなど

}
