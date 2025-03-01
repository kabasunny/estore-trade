package tachibana

import "estore-trade/internal/domain"

// GetOperationStatusKabu は市場と単位に対応する運用ステータスを返す
func (tc *TachibanaClientImple) GetOperationStatusKabu(listedMarket string, unit string) (domain.OperationStatusKabu, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	marketMap, ok := tc.operationStatusKabuMap[listedMarket]
	if !ok {
		return domain.OperationStatusKabu{}, false
	}
	status, ok := marketMap[unit]
	return status, ok
}
