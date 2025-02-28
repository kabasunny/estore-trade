package tachibana

import (
	"fmt"
	"time"
)

// GetRequestURL はキャッシュされた仮想URL（REQUEST）を返す
func (tc *TachibanaClientImple) GetRequestURL() (string, error) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	if time.Now().Before(tc.Expiry) && tc.Loggined && tc.RequestURL != "" {
		return tc.RequestURL, nil
	}
	return "", fmt.Errorf("request URL not found, need to Login")
}
