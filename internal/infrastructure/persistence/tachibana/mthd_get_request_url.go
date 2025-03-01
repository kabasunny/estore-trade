package tachibana

import (
	"fmt"
	"time"
)

// GetRequestURL はキャッシュされた仮想URL（REQUEST）を返す
func (tc *TachibanaClientImple) GetRequestURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.requestURL != "" {
		return tc.requestURL, nil
	}
	return "", fmt.Errorf("request URL not found, need to Login")
}
