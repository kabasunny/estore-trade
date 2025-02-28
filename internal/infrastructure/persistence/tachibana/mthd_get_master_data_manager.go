package tachibana

import "estore-trade/internal/domain"

func (tc *TachibanaClientImple) GetMasterData() *domain.MasterData {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	return tc.masterData
}
