package tachibana

import "estore-trade/internal/domain"

// GetSystemStatus はシステムステータスを取得します。
func (tc *TachibanaClientImple) GetSystemStatus() domain.SystemStatus {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	return tc.systemStatus
}
