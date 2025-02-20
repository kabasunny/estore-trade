package tachibana

// GetDateInfo は日付情報を取得します。
func (tc *TachibanaClientImple) GetDateInfo() DateInfo {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.dateInfo
}
