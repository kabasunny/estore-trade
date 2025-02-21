// internal/domain/model.go
package domain

import "time"

// 自動売買システムの中核となるデータ構造

// 取引口座
type Account struct {
	ID        string     // アカウントID（ユニーク識別子）
	Balance   float64    // アカウントの現在の残高
	Positions []Position // ポジションのリスト（取引中または保有中のポジション）
	CreatedAt time.Time  // アカウント作成日時
	UpdatedAt time.Time  // アカウントの最終更新日時
}
