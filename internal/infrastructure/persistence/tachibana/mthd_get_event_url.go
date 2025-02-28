package tachibana

import (
	"fmt"
	"time"
)

// GetEventURL はキャッシュされた仮想URL（EVENT）を返す
func (tc *TachibanaClientImple) GetEventURL() (string, error) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	if time.Now().Before(tc.Expiry) && tc.Loggined && tc.EventURL != "" {
		return tc.EventURL, nil
	}
	return "", fmt.Errorf("event URL not found, need to Login")
}
