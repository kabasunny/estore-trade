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
	tc.Mu.RLock()
	if time.Now().Before(tc.Expiry) && tc.Loggined && tc.RequestURL != "" && tc.MasterURL != "" && tc.PriceURL != "" && tc.EventURL != "" {
		tc.Mu.RUnlock()
		return nil
	}
	tc.Mu.RUnlock()

	// Write Lock: 新しいURLを取得
	tc.Mu.Lock()
	defer tc.Mu.Unlock()

	loggined, err := login(ctx, tc, userID, password)
	tc.Loggined = loggined

	return err
}
