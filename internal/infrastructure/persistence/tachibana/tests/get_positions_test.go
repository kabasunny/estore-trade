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
		client := tachibana.CreateTestClient(t, nil) // MasterDataは不要
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		positions, err := client.GetPositions(context.Background())
		assert.NoError(t, err)

		// 返ってくる建玉情報の数や内容をチェック (APIの応答に依存)
		// デモ環境では建玉がない可能性があるので、空の場合も許容する
		if len(positions) > 0 {
			// 建玉がある場合の追加のチェック (例)
			assert.NotEmpty(t, positions[0].ID)
			assert.NotEmpty(t, positions[0].Symbol)
			assert.NotEmpty(t, positions[0].Side)
		}
	})

	t.Run("異常系: ログイン前にGetPositionsを呼び出すとエラー", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		// ログインしない

		_, err := client.GetPositions(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "not logged in", err.Error()) // ログインしていないときのエラーメッセージ
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/get_positions_test.go
