package tachibana

import (
	"context"
	"time"

	"estore-trade/internal/config"
)

// Login は API にログインし、仮想URLを返す。有効期限内ならキャッシュされたURLを返す
func (tc *TachibanaClientImple) Login(ctx context.Context, cfg interface{}) error {
	userID := cfg.(*config.Config).TachibanaUserID // 型アサーション
	password := cfg.(*config.Config).TachibanaPassword

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

	loggined, err := login(ctx, tc, userID, password)
	tc.loggined = loggined

	return err
}
