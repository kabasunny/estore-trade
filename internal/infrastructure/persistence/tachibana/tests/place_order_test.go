// internal/infrastructure/persistence/tachibana/tests/place_order_test.go
package tachibana_test

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	// 正常系のテスト (成行注文)
	// t.Run("正常系: 成行注文が成功すること", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil)
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())

	// 	order := &domain.Order{
	// 		Symbol:    "7974", // 任天堂
	// 		Side:      "buy",
	// 		OrderType: "market",
	// 		Quantity:  100,
	// 	}

	// 	placedOrder, err := client.PlaceOrder(context.Background(), order)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, placedOrder)
	// 	assert.NotEmpty(t, placedOrder.ID)             // 注文IDが設定されている
	// 	assert.Equal(t, "pending", placedOrder.Status) // ステータスが"pending"

	// 	// 1秒待機
	// 	time.Sleep(1 * time.Second)
	// })

	// 正常系のテスト (指値注文)
	// t.Run("正常系: 指値注文が成功すること", func(t *testing.T) {

	// 	// 1秒待機
	// 	// time.Sleep(2 * time.Second)
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil)
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())

	// 	order := &domain.Order{
	// 		Symbol:    "7974",
	// 		Side:      "buy",
	// 		OrderType: "limit",
	// 		Quantity:  100,
	// 		Price:     10000.0, // 指値価格
	// 	}

	// 	placedOrder, err := client.PlaceOrder(context.Background(), order)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, placedOrder)
	// 	assert.NotEmpty(t, placedOrder.ID)
	// 	assert.Equal(t, "pending", placedOrder.Status)
	// 	// 1秒待機
	// 	time.Sleep(1 * time.Second)
	// })

	// 正常系のテスト（逆指値）
	t.Run("正常系: 逆指値注文が成功すること", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "buy",
			OrderType:    "stop",
			Quantity:     100,
			TriggerPrice: 4000.0,  // 逆指値条件価格
			Price:        10500.0, // 逆指値がトリガーされた後の注文価格（指値）
		}

		placedOrder, err := client.PlaceOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.ID)
		assert.Equal(t, "pending", placedOrder.Status)
		// 1秒待機
		time.Sleep(1 * time.Second)
	})

	// 	// 正常系のテスト（逆指値+指値）
	// 	t.Run("正常系: 逆指値+指値注文が成功すること", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, nil)
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:       "7974", // 例: 任天堂
	// 			Side:         "buy",
	// 			OrderType:    "stop_limit",
	// 			Quantity:     100,
	// 			TriggerPrice: 4000.0, // 逆指値条件価格
	// 			Price:        3500.0, // 逆指値がトリガーされた後の注文価格（指値）
	// 		}

	// 		placedOrder, err := client.PlaceOrder(context.Background(), order)
	// 		assert.NoError(t, err)
	// 		assert.NotNil(t, placedOrder)
	// 		assert.NotEmpty(t, placedOrder.ID)
	// 		assert.Equal(t, "pending", placedOrder.Status)
	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: 無効な注文 (Side)
	// 	t.Run("異常系: 無効な Side でエラー", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, nil)
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:    "7974",
	// 			Side:      "invalid_side", // 無効な値
	// 			OrderType: "market",
	// 			Quantity:  100,
	// 		}

	// 		_, err = client.PlaceOrder(context.Background(), order)
	// 		assert.Error(t, err) // エラーが発生するはず
	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: 無効な注文 (OrderType)
	// 	t.Run("異常系: 無効な OrderType でエラー", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, nil)
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:    "7974",
	// 			Side:      "buy",
	// 			OrderType: "invalid_type", // 無効な値
	// 			Quantity:  100,
	// 		}

	// 		_, err = client.PlaceOrder(context.Background(), order)
	// 		assert.Error(t, err)
	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: 数量0
	// 	t.Run("異常系: 数量 0 でエラー", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, nil)
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:    "7974",
	// 			Side:      "buy",
	// 			OrderType: "market",
	// 			Quantity:  0, // 無効な値
	// 		}

	// 		_, err = client.PlaceOrder(context.Background(), order)
	// 		assert.Error(t, err)
	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: contextキャンセル
	// 	t.Run("異常系: context キャンセルでエラー", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, nil)
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:    "7974",
	// 			Side:      "buy",
	// 			OrderType: "market",
	// 			Quantity:  100,
	// 		}

	// 		ctx, cancel := context.WithCancel(context.Background())
	// 		cancel() // キャンセル

	// 		_, err = client.PlaceOrder(ctx, order)
	// 		assert.Error(t, err) // キャンセルによるエラー
	// 		assert.True(t, errors.Is(err, context.Canceled))
	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: APIエラー (無効な銘柄コード)
	// 	t.Run("異常系: 無効な銘柄コードでAPIエラー", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, nil)
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:    "invalid_code", // 無効な値
	// 			Side:      "buy",
	// 			OrderType: "market",
	// 			Quantity:  100,
	// 		}

	// 		_, err = client.PlaceOrder(context.Background(), order)
	// 		assert.Error(t, err)
	// 		// "API returned an error" のようなエラーメッセージが含まれていることを確認
	// 		assert.Contains(t, err.Error(), "API returned an error")
	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: sOrderNumberがない
	// 	t.Run("異常系: レスポンスにsOrderNumberがない", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, &domain.MasterData{}) // クライアントを再作成
	// 		err := client.Login(context.Background(), nil)
	// 		assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		// masterURLを書き換えて存在しないURLにする
	// 		originalMasterURL := client.GetMasterURLForTest() // 正しいURLを保持
	// 		client.SetRequestURLForTest("https://invalid.example.com/")

	// 		order := &domain.Order{
	// 			Symbol:    "7974",
	// 			Side:      "buy",
	// 			OrderType: "market",
	// 			Quantity:  100,
	// 		}

	// 		_, err = client.PlaceOrder(context.Background(), order)
	// 		assert.Error(t, err)

	// 		// requestURLを戻す
	// 		client.SetRequestURLForTest(originalMasterURL) // URLを戻す

	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	// 異常系: APIエラー (ログインしないで注文)
	// 	t.Run("異常系: ログインしないで注文", func(t *testing.T) {
	// 		client := tachibana.CreateTestClient(t, &domain.MasterData{}) // クライアントを再作成
	// 		//err := client.Login(context.Background(), nil) Loginしない
	// 		//assert.NoError(t, err)
	// 		defer client.Logout(context.Background())

	// 		order := &domain.Order{
	// 			Symbol:    "7974",
	// 			Side:      "buy",
	// 			OrderType: "market",
	// 			Quantity:  100,
	// 		}

	// 		_, err := client.PlaceOrder(context.Background(), order)
	// 		fmt.Printf("err: %v \n", err)
	// 		assert.Error(t, err)

	// 		// 1秒待機
	// 		time.Sleep(1 * time.Second)
	// 	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order_test.go
