package tachibana

import "estore-trade/internal/domain"

// GetSystemStatus はシステムステータスを取得します。
func (tc *TachibanaClientImple) GetSystemStatus() domain.SystemStatus {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterData.SystemStatus
}
