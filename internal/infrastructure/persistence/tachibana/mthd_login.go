// internal/infrastructure/persistence/tachibana/mthd_login.go
package tachibana

import (
	"context"
	"time"
	//"estore-trade/internal/config" // 不要
)

// Login は API にログインし、仮想URLを返す。有効期限内ならキャッシュされたURLを返す
func (tc *TachibanaClientImple) Login(ctx context.Context, cfg interface{}) error {
	// 不要
	//userID := cfg.(*config.Config).TachibanaUserID // 型アサーション
	//password := cfg.(*config.Config).TachibanaPassword

	// Read Lock: キャッシュされたURLが有効ならそれを返す
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.requestURL != "" && tc.masterURL != "" && tc.priceURL != "" && tc.eventURL != "" {
		tc.mu.RUnlock()
		return nil
	}
	tc.mu.RUnlock()

	// Write Lock: 新しいURLを取得
	tc.mu.Lock()
	defer tc.mu.Unlock()

	//loggined, err := login(ctx, tc, userID, password) // configから取得していた箇所を変更
	loggined, err := login(ctx, tc, tc.userID, tc.password) //自身のフィールドを参照
	tc.loggined = loggined

	return err
}
