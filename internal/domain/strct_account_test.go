// internal/domain/strct_account_test.go
package domain

import (
	"testing"
	"time"
)

func TestAccount_Validate(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name    string  // テストケースの名前
		account Account // テストするアカウントの構造体
		wantErr bool    // エラーが発生することを期待するか
	}{
		{
			// 有効なアカウントのテストケース
			name: "valid account",
			account: Account{
				ID:               "test-account-id",
				UserID:           "test-user-id",
				AccountType:      "special",
				Balance:          100000,
				AvailableBalance: 100000,
				Margin:           0,
				Positions:        []Position{},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: false,
		},
		{
			// 無効なアカウントのテストケース - IDが空
			name: "invalid account - empty ID",
			account: Account{
				UserID:           "test-user-id",
				AccountType:      "special",
				Balance:          10000,
				AvailableBalance: 100000,
				Positions:        []Position{},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効なアカウントのテストケース - UserIDが空
			name: "invalid account - empty UserID",
			account: Account{
				ID:               "test-account-id",
				AccountType:      "special",
				Balance:          10000,
				AvailableBalance: 100000,
				Positions:        []Position{},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効なアカウントのテストケース - 無効なAccountType
			name: "invalid account - invalid AccountType",
			account: Account{
				ID:               "test-account-id",
				UserID:           "test-user-id",
				AccountType:      "invalid",
				Balance:          10000,
				AvailableBalance: 100000,
				Positions:        []Position{},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			// 無効なアカウントのテストケース - 負のBalance
			name: "invalid account - negative Balance",
			account: Account{
				ID:               "test-account-id",
				UserID:           "test-user-id",
				AccountType:      "special",
				Balance:          -100,
				AvailableBalance: 100000,
				Positions:        []Position{},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		// 他のテストケースを追加 (AvailableBalance, Margin など)
	}

	// 各テストケースをループで実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// アカウントの検証を実行
			err := tt.account.Validate()
			// エラーの結果が期待値と一致するかを確認
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
