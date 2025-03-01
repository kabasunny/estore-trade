// internal/infrastructure/persistence/tachibana/tests/get_system_status_test.go
package tachibana_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetSystemStatus(t *testing.T) {
	// SetupTestClient を使うと、DownloadMasterData のモックが実行され、
	// SystemStatus が設定される
	client, _ := tachibana.SetupTestClient(t)

	t.Run("SystemStatus is returned", func(t *testing.T) {
		expectedStatus := domain.SystemStatus{
			// SetupTestClient 内でモックデータが設定されているはず
			SystemState: "1", // 例
		}

		actualStatus := client.GetSystemStatus()

		// DeepEqual を使って構造体全体を比較
		assert.Equal(t, expectedStatus.SystemState, actualStatus.SystemState) //必要なものだけ比較
	})
}
