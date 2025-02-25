// internal/domain/iface_account_repository.go
package domain

import (
	"context"
)

// AccountRepository は取引アカウント（Account）データの永続化操作を抽象化
type AccountRepository interface {
	GetAccount(ctx context.Context, id string) (*Account, error)             // 指定されたIDのアカウントを取得
	GetAccountByUserID(ctx context.Context, userID string) (*Account, error) // ユーザーIDでアカウントを取得
	UpdateAccount(ctx context.Context, account *Account) error               // 指定されたアカウントのデータを更新
	CreateAccount(ctx context.Context, account *Account) error               // アカウント作成
	// ListAccounts(ctx context.Context) ([]*Account, error) // 全アカウント取得は不要と判断
	// DeleteAccount(ctx context.Context, id string) error // アカウント削除は不要と判断
}
