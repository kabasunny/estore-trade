// internal/usecase/tests/trading_usecase_integration_test.go
package usecase_test

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order" // 変更
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"
	"estore-trade/test/docker" // Docker テストヘルパー

	// UUID 生成用

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

func TestTradingUsecase_Integration(t *testing.T) {
	// テスト用の Logger を作成
	logger := zaptest.NewLogger(t)

	// Docker を使用してテスト用 DB を起動
	db, cleanup, err := docker.SetupTestDatabase()
	require.NoError(t, err)
	defer cleanup()

	// テストケースごとにDBを初期化
	cleanupDB := setUp(db)
	defer cleanupDB()

	// テスト用の MasterData を作成
	masterData := &domain.MasterData{
		IssueMap: map[string]domain.IssueMaster{
			"7974": {IssueCode: "7974", IssueName: "任天堂", TradingUnit: 100},
		},
		// IssueMarketMap を追加
		IssueMarketMap: map[string]map[string]domain.IssueMarketMaster{
			"7974": {
				"00": {
					IssueCode:  "7974",
					MarketCode: "00",
					// 他の必要なフィールドも設定
				},
			},
		},
	}

	// TachibanaClient を作成 (テスト用の設定と MasterData を使用)
	client := tachibana.CreateTestClient(t, masterData) // masterData を渡す

	// ★★★ ここで SetTargetIssues を呼び出す ★★★
	err = client.SetTargetIssues(context.Background(), []string{"7974"})
	assert.NoError(t, err)

	// OrderRepository を作成
	orderRepo := order.NewOrderRepository(db)

	// TradingUsecase を作成
	tradingUsecase := usecase.NewTradingUsecase(client, logger, orderRepo, nil, client.GetConfig())

	// テストケース
	t.Run("PlaceOrder and GetOrderStatus", func(t *testing.T) {
		ctx := context.Background()

		// ログイン
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// テスト用の注文データを作成
		order := &domain.Order{
			// UUID:       uuid.NewString(), // UUID を生成
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: "00", // 東証
		}

		// 注文を発注
		placedOrder, err := tradingUsecase.PlaceOrder(ctx, order)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.TachibanaOrderID)

		// UpdateOrderByEventで更新
		// event　を定義し、orderを更新
		event := &domain.OrderEvent{
			EventType: "EC",
			Order:     placedOrder,
		}
		err = tradingUsecase.UpdateOrderByEvent(ctx, event)
		assert.NoError(t, err)

		// GetOrderStatusで確認 (リトライ処理を追加)
		var retrievedOrder *domain.Order
		maxRetries := 5
		retryInterval := 1 * time.Second
		expectedStatus := "" //期待するステータス
		//取引時間内なら、filledを期待する
		if usecase.IsWithinTradingHours() {
			expectedStatus = "全部約定" // 取引時間内なら "filled" を期待
			println("取引時間内")
		} else {
			println("取引時間外")
		}
		for i := 0; i < maxRetries; i++ {
			retrievedOrder, err = tradingUsecase.GetOrderStatus(ctx, placedOrder.TachibanaOrderID, time.Now().Format("20060102"))
			assert.NoError(t, err)
			assert.NotNil(t, retrievedOrder)
			if retrievedOrder.Status == expectedStatus {
				break
			}
			t.Logf("GetOrderStatus retry: %d", i+1)
			time.Sleep(retryInterval)
		}

		// 取得した注文情報が、発注した注文と一致することを確認 (一部)
		assert.Equal(t, placedOrder.TachibanaOrderID, retrievedOrder.TachibanaOrderID)
		if expectedStatus != "" { //期待するステータスがある場合
			assert.Equal(t, expectedStatus, retrievedOrder.Status) // Status が期待値と一致することを確認
		}

	})
	t.Run("CancelOrder", func(t *testing.T) {
		ctx := context.Background()

		// ログイン (CancelOrder はログイン状態に依存するため)
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// 1. 注文を作成 (キャンセル対象の注文)
		// --- 逆指値買い注文 ---
		orderToCancel := &domain.Order{
			Symbol:       "7974",
			Side:         "long",
			OrderType:    "stop", // 逆指値
			Quantity:     100,
			MarketCode:   "00",
			TriggerPrice: 12000, // 適当なトリガー価格を設定（現在価格より高め）
			Price:        0,     //トリガー後、成行き
		}
		placedOrderToCancel, err := tradingUsecase.PlaceOrder(ctx, orderToCancel)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrderToCancel)
		assert.NotEmpty(t, placedOrderToCancel.TachibanaOrderID)

		// UpdateOrderByEventで更新
		// event　を定義し、orderを更新
		event := &domain.OrderEvent{
			EventType: "EC",
			Order:     placedOrderToCancel, //修正
		}
		err = tradingUsecase.UpdateOrderByEvent(ctx, event)
		assert.NoError(t, err)

		// 2. 注文をキャンセル
		err = tradingUsecase.CancelOrder(ctx, placedOrderToCancel.TachibanaOrderID) //修正
		assert.NoError(t, err)

	})
}
func setUp(db *gorm.DB) func() {
	// ここで初期化処理を行う（テーブルのクリーンアップなど）
	db.Exec("DELETE FROM orders") // 例：orders テーブルを空にする

	// クリーンアップ関数を返す
	return func() {
		// ここでクリーンアップ処理を行う（必要に応じて）
		db.Exec("DELETE FROM orders")
	}
}

// go test -v ./internal/usecase/tests/trading_usecase_integration_test.go
