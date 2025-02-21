package domain

import (
	"context"
)

// 取引アカウント（Account）データの永続化操作を抽象化
type AccountRepository interface {
	GetAccount(ctx context.Context, id string) (*Account, error) // 指定されたIDのアカウントを取得
	UpdateAccount(ctx context.Context, account *Account) error   // 指定されたアカウントのデータを更新
}
