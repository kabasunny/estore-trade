// internal/infrastructure/persistence/tachibana/tests/get_date_info_test.go
package tachibana_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetDateInfo(t *testing.T) {
	// テスト用の DateInfo データを作成
	expectedDateInfo := domain.DateInfo{
		DateKey:           "001",      // 当日基準
		PrevBusinessDay1:  "20231104", // 任意の値
		TheDay:            "20231105",
		NextBusinessDay1:  "20231106", // 任意の値
		StockDeliveryDate: "20231108", // 任意の値
	}

	// SetupTestClient を使うと、DownloadMasterData のモックが実行され、DateInfo が設定される
	client, _ := tachibana.SetupTestClient(t)

	// テスト用の DateInfo データを設定 (SetupTestClient で設定される値を上書き)
	tachibana.SetDateInfoForTest(client, expectedDateInfo)

	t.Run("DateInfo is returned", func(t *testing.T) {
		actualDateInfo := client.GetDateInfo()
		assert.Equal(t, expectedDateInfo, actualDateInfo) // 期待値と実際値が一致するか確認
	})
}
