// internal/infrastructure/persistence/tachibana/mthd_get_call_price.go
package tachibana

import "estore-trade/internal/domain"

// GetCallPrice は指定された単位番号に対応する呼値情報を返します。
func (tc *TachibanaClientImple) GetCallPrice(unitNumber string) (domain.CallPrice, bool) {
	tc.Mu.RLock() // Read Lock を取得
	defer tc.Mu.RUnlock()
	callPrice, ok := tc.CallPriceMap[unitNumber]
	return callPrice, ok
}
