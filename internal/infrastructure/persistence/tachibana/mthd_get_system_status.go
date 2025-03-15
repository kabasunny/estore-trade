// internal/infrastructure/persistence/tachibana/mthd_get_system_status.go
package tachibana

import (
	"context"
	"estore-trade/internal/domain"
)

// GetSystemStatus はシステムステータスを取得します。
func (tc *TachibanaClientImple) GetSystemStatus(ctx context.Context) domain.SystemStatus {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// masterData が nil の場合はデフォルト値を返す (追加)
	if tc.masterData == nil {
		return domain.SystemStatus{} // 空の SystemStatus
	}
	return tc.masterData.SystemStatus
}
