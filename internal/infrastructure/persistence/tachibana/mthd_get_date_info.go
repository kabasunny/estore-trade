package tachibana

import "estore-trade/internal/domain"

// GetDateInfo は日付情報を取得します。
func (tc *TachibanaClientImple) GetDateInfo() domain.DateInfo {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterData.DateInfo
}
