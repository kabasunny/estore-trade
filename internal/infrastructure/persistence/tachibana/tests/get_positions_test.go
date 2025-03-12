// internal/infrastructure/persistence/tachibana/tests/get_positions_test.go
package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestGetPositions(t *testing.T) {
	t.Run("正常系: ログイン後にGetPositionsを呼び出せること", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		positions, err := client.GetPositions(context.Background())
		assert.NoError(t, err)

		// 建玉の数をチェック (例: 0以上であることを確認)
		assert.GreaterOrEqual(t, len(positions), 0, "Expected at least 0 positions")

		// すべての建玉の内容をチェック
		for i, p := range positions {
			assert.NotEmpty(t, p.ID, "Position ID should not be empty for position %d", i)
			assert.NotEmpty(t, p.Symbol, "Symbol should not be empty for position %d", i)
			assert.NotEmpty(t, p.Side, "Side should not be empty for position %d", i)
			// 他のフィールドも必要に応じてチェック (例: p.Quantity, p.Price など)
		}
	})

	// t.Run("異常系: ログイン前にGetPositionsを呼び出すとエラー", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	// ログインしない

	// 	_, err := client.GetPositions(context.Background())
	// 	assert.Error(t, err)
	// 	assert.Equal(t, "not logged in", err.Error()) // ログインしていないときのエラーメッセージ
	// })
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/get_positions_test.go
