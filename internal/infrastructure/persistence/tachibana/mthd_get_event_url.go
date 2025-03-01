package tachibana

import (
	"fmt"
	"time"
)

// GetEventURL はキャッシュされた仮想URL（EVENT）を返す
func (tc *TachibanaClientImple) GetEventURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.eventURL != "" {
		return tc.eventURL, nil
	}
	return "", fmt.Errorf("event URL not found, need to Login")
}
