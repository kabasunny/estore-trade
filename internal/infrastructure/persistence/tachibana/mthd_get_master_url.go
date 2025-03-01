package tachibana

import (
	"fmt"
	"time"
)

// GetMasterURL はキャッシュされた仮想URL（Master）を返す
func (tc *TachibanaClientImple) GetMasterURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.masterURL != "" {
		return tc.masterURL, nil
	}
	return "", fmt.Errorf("master URL not found, need to Login")
}
