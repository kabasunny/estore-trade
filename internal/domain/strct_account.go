// internal/domain/strct_account.go
package domain

import (
	"time"
)

// Account 取引口座
type Account struct {
	ID               string     `json:"id"`                // アカウントID（ユニーク識別子, UUID）
	UserID           string     `json:"user_id"`           // ユーザーID (外部キー)
	AccountType      string     `json:"account_type"`      // 口座種別 (特定口座を想定: "special")
	Balance          float64    `json:"balance"`           // アカウント残高
	AvailableBalance float64    `json:"available_balance"` // 利用可能残高
	Margin           float64    `json:"margin"`            // 証拠金 (今回は使用しない)
	Positions        []Position `json:"positions"`         // ポジションのリスト（取引中または保有中のポジション）
	CreatedAt        time.Time  `json:"created_at"`        // アカウント作成日時
	UpdatedAt        time.Time  `json:"updated_at"`        // アカウント最終更新日時
}
