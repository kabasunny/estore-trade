package tachibana

// GetSystemStatus はシステムステータスを取得します。
func (tc *TachibanaClientImple) GetSystemStatus() SystemStatus {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.systemStatus
}
