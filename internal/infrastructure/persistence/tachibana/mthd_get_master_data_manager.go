package tachibana

import "estore-trade/internal/domain"

func (tc *TachibanaClientImple) GetMasterData() *domain.MasterData {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterData
}
