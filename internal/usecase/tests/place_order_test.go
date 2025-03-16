// internal/usecase/tests/place_order_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"
	"estore-trade/test/testdata"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestTradingUsecase_PlaceOrder(t *testing.T) {
	logger := zap.NewNop()
	cfg := testdata.NewTestConfig() // testdata パッケージは別途作成

	tests := []struct {
		name          string
		order         *domain.Order
		setupMock     func(*tachibana.MockTachibanaClient, *usecase.MockOrderRepository)
		expectedOrder *domain.Order
		expectedError string // エラーメッセージの部分一致で検証
	}{
		{
			name: "正常系: 現物買い注文",
			order: &domain.Order{
				UUID:       uuid.NewString(), // UUID を生成
				Symbol:     "7974",
				Side:       "long", // long
				OrderType:  "market",
				Quantity:   100,
				MarketCode: "00",
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "1"})
				mockClient.On("GetIssueMaster", mock.Anything, "7974").Return(domain.IssueMaster{TradingUnit: 100}, true)
				mockClient.On("CheckPriceIsValid", mock.Anything, "7974", mock.AnythingOfType("float64"), false).Return(true, nil)
				mockClient.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(&domain.Order{TachibanaOrderID: "12345", Status: "pending"}, nil) // Status も返す
				mockRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)
			},
			expectedOrder: &domain.Order{TachibanaOrderID: "12345", Status: "pending"},
			expectedError: "",
		},
		{
			name: "異常系: システムが閉局状態",
			order: &domain.Order{
				UUID:       uuid.NewString(), // UUID を生成
				Symbol:     "7974",
				Side:       "long", // long
				OrderType:  "market",
				Quantity:   100,
				MarketCode: "00",
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "0"})
			},
			expectedOrder: nil,
			expectedError: "system is not in service",
		},
		{
			name: "異常系: 銘柄マスタ取得失敗",
			order: &domain.Order{
				UUID:      uuid.NewString(), // UUID を生成
				Symbol:    "99999",
				Side:      "long", // long
				OrderType: "market",
				Quantity:  100,
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "1"})
				mockClient.On("GetIssueMaster", mock.Anything, "99999").Return(domain.IssueMaster{}, false)
			},
			expectedOrder: nil,
			expectedError: "invalid issue code",
		},
		{
			name: "異常系: 注文数量が不正",
			order: &domain.Order{
				UUID:       uuid.NewString(), // UUID を生成
				Symbol:     "7974",
				Side:       "long", // long
				OrderType:  "market",
				Quantity:   50, // 売買単位の倍数でない
				MarketCode: "00",
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "1"})
				mockClient.On("GetIssueMaster", mock.Anything, "7974").Return(domain.IssueMaster{TradingUnit: 100}, true) // 売買単位は100
			},
			expectedOrder: nil,
			expectedError: "invalid order quantity",
		},
		{
			name: "異常系: 呼値エラー",
			order: &domain.Order{
				UUID:       uuid.NewString(), // UUID を生成
				Symbol:     "7974",
				Side:       "long", // long
				OrderType:  "limit",
				Quantity:   100,
				Price:      0.1, // 不正な価格
				MarketCode: "00",
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "1"})
				mockClient.On("GetIssueMaster", mock.Anything, "7974").Return(domain.IssueMaster{TradingUnit: 100}, true)
				mockClient.On("CheckPriceIsValid", mock.Anything, "7974", 0.1, false).Return(false, nil) // 呼値エラー
			},
			expectedOrder: nil,
			expectedError: "invalid order price",
		},
		{
			name: "異常系: TachibanaClient.PlaceOrder でエラー",
			order: &domain.Order{
				UUID:       uuid.NewString(), // UUID を生成
				Symbol:     "7974",
				Side:       "long", // long
				OrderType:  "market",
				Quantity:   100,
				MarketCode: "00",
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "1"})
				mockClient.On("GetIssueMaster", mock.Anything, "7974").Return(domain.IssueMaster{TradingUnit: 100}, true)
				mockClient.On("CheckPriceIsValid", mock.Anything, "7974", mock.AnythingOfType("float64"), false).Return(true, nil)
				mockClient.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil, errors.New("API error")) // PlaceOrder がエラーを返す
			},
			expectedOrder: nil,
			expectedError: "API error",
		},
		{
			name: "異常系: OrderRepository.CreateOrder でエラー",
			order: &domain.Order{
				UUID:       uuid.NewString(), // UUID を生成
				Symbol:     "7974",
				Side:       "long", // long
				OrderType:  "market",
				Quantity:   100,
				MarketCode: "00",
			},
			setupMock: func(mockClient *tachibana.MockTachibanaClient, mockRepo *usecase.MockOrderRepository) {
				mockClient.On("GetSystemStatus", mock.Anything).Return(domain.SystemStatus{SystemState: "1"})
				mockClient.On("GetIssueMaster", mock.Anything, "7974").Return(domain.IssueMaster{TradingUnit: 100}, true)
				mockClient.On("CheckPriceIsValid", mock.Anything, "7974", mock.AnythingOfType("float64"), false).Return(true, nil)
				mockClient.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(&domain.Order{TachibanaOrderID: "12345"}, nil)
				mockRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(errors.New("DB error")) // CreateOrder がエラーを返す
			},
			expectedOrder: &domain.Order{TachibanaOrderID: "12345"}, // APIからの戻り値は期待通り
			expectedError: "",                                       // DBエラーはログ出力するが、上位にはエラーを返さない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTachibanaClient := new(tachibana.MockTachibanaClient)
			mockOrderRepo := new(usecase.MockOrderRepository)

			tt.setupMock(mockTachibanaClient, mockOrderRepo)

			tradingUsecase := usecase.NewTradingUsecase(mockTachibanaClient, logger, mockOrderRepo, nil, cfg)
			//mock呼び出しに変更
			//tt.setupMock(mockTachibanaClient, mockOrderRepo)
			placedOrder, err := tradingUsecase.PlaceOrder(context.Background(), tt.order)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, placedOrder) {
					assert.Equal(t, tt.expectedOrder.TachibanaOrderID, placedOrder.TachibanaOrderID)
					if tt.expectedOrder.Status != "" {
						assert.Equal(t, tt.expectedOrder.Status, placedOrder.Status)
					}
				}
			}

			mockTachibanaClient.AssertExpectations(t)
			mockOrderRepo.AssertExpectations(t)
		})
	}
} // go test -v ./internal/usecase/tests/place_order_test.go
