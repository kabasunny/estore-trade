// internal/infrastructure/persistence/tachibana/mthd_get_call_price.go
package tachibana

import "estore-trade/internal/domain"

// GetCallPrice は指定された単位番号に対応する呼値情報を返します。
func (tc *TachibanaClientImple) GetCallPrice(unitNumber string) (domain.CallPrice, bool) {
	tc.mu.RLock() // Read Lock を取得
	defer tc.mu.RUnlock()

	if tc.masterData == nil { // masterData が nil の場合は、データがないことを示す
		return domain.CallPrice{}, false
	}
	callPrice, ok := tc.masterData.CallPriceMap[unitNumber] // masterData 経由でアクセス
	return callPrice, ok
}
