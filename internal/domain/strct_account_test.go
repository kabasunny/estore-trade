// internal/domain/strct_account_test.go
package domain

import (
	"testing"
	"time"
)

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		wantErr bool // エラーが発生することを期待するか
	}{
		{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
