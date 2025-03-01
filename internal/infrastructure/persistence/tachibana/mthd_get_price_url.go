package tachibana

import (
	"fmt"
	"time"
)

// GetPriceURL はキャッシュされた仮想URL（Price）を返す
func (tc *TachibanaClientImple) GetPriceURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.priceURL != "" {
		return tc.priceURL, nil
	}
	return "", fmt.Errorf("price URL not found, need to Login")
}
