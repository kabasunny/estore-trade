package tachibana

// GetCallPrice は指定された単位番号に対応する呼値情報を返します。
func (tc *TachibanaClientImple) GetCallPrice(unitNumber string) (CallPrice, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	callPrice, ok := tc.callPriceMap[unitNumber]
	return callPrice, ok
}
