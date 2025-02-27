// internal/domain/strct_position_test.go
package domain

import (
	"testing"
	"time"
)

func TestPosition_Validate(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name     string   // テストケースの名前
		position Position // テストするポジションの構造体
		wantErr  bool     // エラーが発生することを期待するか
	}{
		{
			// 有効なポジションのテストケース
			name: "valid position",
			position: Position{
				Symbol:   "7203",
				Side:     "long",
				Quantity: 100,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: false,
		},
		{
			// 無効なポジションのテストケース - Symbolが空
			name: "invalid position - empty Symbol",
			position: Position{
				Side:     "long",
				Quantity: 100,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効なポジションのテストケース - 無効なSide
			name: "invalid position - invalid Side",
			position: Position{
				Symbol:   "7203",
				Side:     "invalid",
				Quantity: 100,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効なポジションのテストケース - Quantityが0
			name: "invalid position - zero Quantity",
			position: Position{
				Symbol:   "7203",
				Side:     "long",
				Quantity: 0,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効なポジションのテストケース - Priceが0
			name: "invalid position - zero Price",
			position: Position{
				Symbol:   "7203",
				Side:     "long",
				Quantity: 100,
				Price:    0,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
	}

	// 各テストケースをループで実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ポジションの検証を実行
			err := tt.position.Validate()
			// エラーの結果が期待値と一致するかを確認
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
