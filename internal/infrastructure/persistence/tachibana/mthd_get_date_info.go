package tachibana

import "estore-trade/internal/domain"

// GetDateInfo は日付情報を取得します。
func (tc *TachibanaClientImple) GetDateInfo() domain.DateInfo {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	return tc.dateInfo
}
