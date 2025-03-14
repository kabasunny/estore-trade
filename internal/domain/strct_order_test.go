// internal/domain/strct_order_test.go
package domain

import (
	"testing"
	"time"
)

func TestOrder_Validate(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name    string // テストケースの名前
		order   Order  // テストする注文の構造体
		wantErr bool   // エラーが発生することを期待するか
	}{
		{
			// 有効な注文のテストケース
			name: "valid order",
			order: Order{
				UUID:             "test-order-id",
				Symbol:           "7203",
				Side:             "buy",
				OrderType:        "market",
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: false,
		},
		{
			// 無効な注文のテストケース - IDが空
			name: "invalid order - empty ID",
			order: Order{
				Symbol:           "7203",
				Side:             "buy",
				OrderType:        "market",
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効な注文のテストケース - Symbolが空
			name: "invalid order - empty Symbol",
			order: Order{
				UUID:             "test-order-id",
				Side:             "buy",
				OrderType:        "market",
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効な注文のテストケース - 無効なSide
			name: "invalid order - invalid Side",
			order: Order{
				UUID:             "test-order-id",
				Symbol:           "7203",
				Side:             "invalid",
				OrderType:        "market",
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効な注文のテストケース - 無効なOrderType
			name: "invalid order - invalid OrderType",
			order: Order{
				UUID:             "test-order-id",
				Symbol:           "7203",
				Side:             "buy",
				OrderType:        "invalid",
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効な注文のテストケース - Quantityが無効
			name: "invalid order - invalid Quantity",
			order: Order{
				UUID:             "test-order-id",
				Symbol:           "7203",
				Side:             "buy",
				OrderType:        "market",
				Quantity:         0,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効な注文のテストケース - 指値注文で価格が0
			name: "invalid order - limit order, price is zero",
			order: Order{
				UUID:             "test-order-id",
				Symbol:           "7203",
				Side:             "buy",
				OrderType:        "limit", // 指値注文
				Price:            0,       // 価格が0
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効な注文のテストケース - 逆指値注文でトリガー価格が0
			name: "invalid order - stop order, TriggerPrice is zero",
			order: Order{
				UUID:             "test-order-id",
				Symbol:           "7203",
				Side:             "buy",
				OrderType:        "stop", // 逆指値注文
				TriggerPrice:     0,      // トリガー価格が0
				Quantity:         100,
				Status:           "pending",
				TachibanaOrderID: "tachibana-order-id",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},

		// 他のテストケース (OrderType, Price, Quantity, Status など)
	}

	// 各テストケースをループで実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注文の検証を実行
			err := tt.order.Validate()
			// エラーの結果が期待値と一致するかを確認
			if (err != nil) != tt.wantErr {
				t.Errorf("Order.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
