// internal/usecase/tests/place_order_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestTradingUsecase_PlaceOrder(t *testing.T) {

	// テスト用の Config (ロガーは使用しないので、nil を設定)
	cfg := &config.Config{}

	// テスト用のロガーを作成
	testLogger := zaptest.NewLogger(t)

	// テスト用の注文データ
	order := &domain.Order{
		Symbol:    "7203",
		Side:      "buy",
		OrderType: "market",
		Quantity:  100,
	}

	// 正常系のテスト (システム稼働中、有効な注文)
	t.Run("valid order", func(t *testing.T) {
		// モッククライアントとモックリポジトリのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		//mockAccountRepo := new(MockAccountRepository) // 今回は使用しない
		// 期待されるメソッド呼び出しと戻り値を設定
		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})             // システム稼働中
		mockClient.On("GetIssueMaster", "7203").Return(domain.IssueMaster{TradingUnit: 100}, true) // IssueMaster に戻す
		mockClient.On("CheckPriceIsValid", "7203", 0.0, false).Return(true, nil)                   // 成行注文なので価格はチェックしない
		mockClient.On("PlaceOrder", mock.Anything, order).Return(&domain.Order{ID: "order-id", Status: "pending"}, nil)
		mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

		// テスト対象のユースケースを作成
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, cfg)

		// PlaceOrder メソッドを呼び出す
		placedOrder, err := uc.PlaceOrder(context.Background(), order)

		// 結果を検証
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.Equal(t, "order-id", placedOrder.ID)
		assert.Equal(t, "pending", placedOrder.Status)

		// モックが期待通りに呼び出されたことを確認
		mockClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	// 異常系のテストケース
	t.Run("system down", func(t *testing.T) {
		// モッククライアントとモックリポジトリのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		// テスト対象のユースケースを作成
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, cfg)
		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "0"}) // システム停止中

		placedOrder, err := uc.PlaceOrder(context.Background(), order)
		assert.Error(t, err) // エラーが発生することを期待
		assert.Nil(t, placedOrder)
		assert.EqualError(t, err, "system is not in service") // エラーメッセージを比較

		mockClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t) //呼ばれないはず
	})

	t.Run("invalid issue code", func(t *testing.T) {
		// モッククライアントとモックリポジトリのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		// テスト対象のユースケースを作成
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, cfg)
		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})
		mockClient.On("GetIssueMaster", "invalid").Return(domain.IssueMaster{}, false) // 無効な銘柄コード, IssueMaster に戻す

		placedOrder, err := uc.PlaceOrder(context.Background(), &domain.Order{Symbol: "invalid", Side: "buy", OrderType: "market", Quantity: 100})
		assert.Error(t, err)
		assert.Nil(t, placedOrder)
		assert.EqualError(t, err, "invalid issue code: invalid") // エラーメッセージを比較

		mockClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t) //呼ばれないはず
	})

	t.Run("invalid order quantity", func(t *testing.T) {
		// モッククライアントとモックリポジトリのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		// テスト対象のユースケースを作成
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, cfg)
		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})
		mockClient.On("GetIssueMaster", "7203").Return(domain.IssueMaster{TradingUnit: 2}, true) // 売買単位が2, IssueMaster に戻す

		placedOrder, err := uc.PlaceOrder(context.Background(), &domain.Order{Symbol: "7203", Side: "buy", OrderType: "market", Quantity: 1}) // 数量が1
		assert.Error(t, err)
		assert.Nil(t, placedOrder)
		assert.EqualError(t, err, "invalid order quantity. must be multiple of 2")

		mockClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t) //呼ばれないはず
	})
	// 異常系: CheckPriceIsValid が false を返すケース (呼値チェックエラー)
	t.Run("invalid price", func(t *testing.T) {
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		testLogger := zaptest.NewLogger(t)
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, &config.Config{})

		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})
		mockClient.On("GetIssueMaster", "7203").Return(domain.IssueMaster{TradingUnit: 100}, true) //IssueMaster に戻す
		mockClient.On("CheckPriceIsValid", "7203", 1000.0, false).Return(false, nil)               // 無効な価格

		placedOrder, err := uc.PlaceOrder(context.Background(), &domain.Order{Symbol: "7203", Side: "buy", OrderType: "limit", Quantity: 100, Price: 1000})
		assert.Error(t, err)
		assert.Nil(t, placedOrder)
		assert.EqualError(t, err, "invalid order price: 1000.000000")

		mockClient.AssertExpectations(t)

	})

	// 異常系: CheckPriceIsValid がエラーを返すケース
	t.Run("CheckPriceIsValid error", func(t *testing.T) {
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		testLogger := zaptest.NewLogger(t)
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, &config.Config{})
		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})
		mockClient.On("GetIssueMaster", "7203").Return(domain.IssueMaster{TradingUnit: 100}, true) //IssueMaster に戻す
		mockClient.On("CheckPriceIsValid", "7203", 1000.0, false).Return(false, errors.New("price check error"))

		placedOrder, err := uc.PlaceOrder(context.Background(), &domain.Order{Symbol: "7203", Side: "buy", OrderType: "limit", Quantity: 100, Price: 1000})
		assert.Error(t, err)
		assert.Nil(t, placedOrder)
		assert.EqualError(t, err, "error checking price validity: price check error")

		mockClient.AssertExpectations(t)

	})

	// 異常系: PlaceOrder (立花証券API) がエラーを返すケース
	t.Run("PlaceOrder error", func(t *testing.T) {
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		testLogger := zaptest.NewLogger(t)
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, &config.Config{})
		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})
		mockClient.On("GetIssueMaster", "7203").Return(domain.IssueMaster{TradingUnit: 100}, true) //IssueMaster に戻す
		mockClient.On("CheckPriceIsValid", "7203", 0.0, false).Return(true, nil)
		mockClient.On("PlaceOrder", mock.Anything, order).Return(nil, errors.New("tachibana API error")) // PlaceOrder でエラー

		placedOrder, err := uc.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Nil(t, placedOrder)
		assert.EqualError(t, err, "tachibana API error")
		mockClient.AssertExpectations(t)
	})
	// 異常系: CreateOrder (DB) がエラーを返すケース  (DBエラーはエラーにしない)
	t.Run("CreateOrder error", func(t *testing.T) {
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		testLogger := zaptest.NewLogger(t)
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, &config.Config{})

		// 正常な注文が返ってくることを想定
		expectedOrder := &domain.Order{ID: "order-id", Status: "pending"}

		mockClient.On("GetSystemStatus").Return(domain.SystemStatus{SystemState: "1"})
		mockClient.On("GetIssueMaster", "7203").Return(domain.IssueMaster{TradingUnit: 100}, true) //IssueMaster に戻す
		mockClient.On("CheckPriceIsValid", "7203", 0.0, false).Return(true, nil)
		mockClient.On("PlaceOrder", mock.Anything, order).Return(expectedOrder, nil)                                        // PlaceOrder は成功
		mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(errors.New("DB error")) // CreateOrder でエラー

		placedOrder, err := uc.PlaceOrder(context.Background(), order)
		assert.NoError(t, err)                      // エラーは発生しない
		assert.Equal(t, expectedOrder, placedOrder) // PlaceOrderの結果が返る
		mockClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}
