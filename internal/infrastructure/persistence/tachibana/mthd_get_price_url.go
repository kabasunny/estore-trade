package tachibana

import (
	"fmt"
	"time"
)

// GetPriceURL はキャッシュされた仮想URL（Price）を返す
func (tc *TachibanaClientImple) GetPriceURL() (string, error) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()
	if time.Now().Before(tc.Expiry) && tc.Loggined && tc.PriceURL != "" {
		return tc.PriceURL, nil
	}
	return "", fmt.Errorf("price URL not found, need to Login")
}
