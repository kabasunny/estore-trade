// internal/domain/strct_account.go
package domain

import (
	"time"
)

// Account 取引口座
type Account struct {
	ID               string     // アカウントID（ユニーク識別子, UUID）
	UserID           string     // ユーザーID (外部キー)
	AccountType      string     // 口座種別 (特定口座を想定: "special")
	Balance          float64    // アカウント残高
	AvailableBalance float64    // 利用可能残高
	Margin           float64    // 証拠金 (今回は使用しない)
	Positions        []Position // ポジションのリスト（取引中または保有中のポジション）
	CreatedAt        time.Time  // アカウント作成日時
	UpdatedAt        time.Time  // アカウント最終更新日時
}
