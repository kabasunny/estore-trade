// internal/domain/strct_order_test.go
package domain

import (
	"testing"
	"time"
)

func TestOrder_Validate(t *testing.T) {
	tests := []struct {
		name    string
		order   Order
		wantErr bool
	}{
		{
			name: "valid order",
			order: Order{
				ID:               "test-order-id",
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
			name: "invalid order - empty Symbol",
			order: Order{
				ID:               "test-order-id",
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
			name: "invalid order - invalid Side",
			order: Order{
				ID:               "test-order-id",
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
			name: "invalid order - invalid OrderType",
			order: Order{
				ID:               "test-order-id",
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
			name: "invalid order - invalid Quantity",
			order: Order{
				ID:               "test-order-id",
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
			name: "invalid order - limit order, price is zero",
			order: Order{
				ID:               "test-order-id",
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
			name: "invalid order - stop order, TriggerPrice is zero",
			order: Order{
				ID:               "test-order-id",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Order.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
