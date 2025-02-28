package tachibana

import "estore-trade/internal/domain"

// GetCallPrice は指定された単位番号に対応する呼値情報を返します。
func (tc *TachibanaClientImple) GetCallPrice(unitNumber string) (domain.CallPrice, bool) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	callPrice, ok := tc.callPriceMap[unitNumber]
	return callPrice, ok
}
