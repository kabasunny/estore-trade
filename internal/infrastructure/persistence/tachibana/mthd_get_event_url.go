package tachibana

import (
	"context"
	"fmt"
	"time"
)

// getEventURL はキャッシュされた仮想URL（EVENT）を返す
func (tc *TachibanaClientImple) getEventURL(ctx context.Context) (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.eventURL != "" {
		return tc.eventURL, nil
	}
	return "", fmt.Errorf("event URL not found, need to Login")
}
