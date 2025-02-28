package tachibana

import (
	"fmt"
	"time"
)

// GetMasterURL はキャッシュされた仮想URL（Master）を返す
func (tc *TachibanaClientImple) GetMasterURL() (string, error) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	if time.Now().Before(tc.Expiry) && tc.Loggined && tc.MasterURL != "" {
		return tc.MasterURL, nil
	}
	return "", fmt.Errorf("master URL not found, need to Login")
}
