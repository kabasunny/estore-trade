// internal/domain/model/account.go
package model

import (
	"context"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model             // ID, CreatedAt, UpdatedAt, DeletedAt を追加
	UserID          string `gorm:"type:varchar(255);unique;not null"`
	TachibanaUserID string `gorm:"type:varchar(255);not null"`
	Password        string `gorm:"type:varchar(255);not null"`
	SecondPassword  string `gorm:"type:varchar(255);not null"`
	AccountType     string `gorm:"type:varchar(50);not null"`
}

// TableName overrides the table name used by User to `profiles`
func (Account) TableName() string {
	return "accounts"
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, id int) (*Account, error)
	GetAccountByUserID(ctx context.Context, userID string) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	GetPositions(ctx context.Context, accountID int) ([]Position, error)
}
