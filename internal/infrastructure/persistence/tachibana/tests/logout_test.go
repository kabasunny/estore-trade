// internal/infrastructure/persistence/tachibana/logout_test.go
package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)

	t.Run("正常系: ログイン後にログアウトできること", func(t *testing.T) {
		// ログイン
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, client.GetLogginedForTest())

		// ログアウト
		err = client.Logout(context.Background())
		assert.NoError(t, err)
		assert.False(t, client.GetLogginedForTest()) // ログアウト後なので false
	})

	t.Run("異常系: ログインしていない状態でログアウトを試みる", func(t *testing.T) {
		// ログインしていない状態でログアウト (エラーにならないことを確認)
		err := client.Logout(context.Background())
		assert.NoError(t, err) // ログアウト自体はエラーにならない
		assert.False(t, client.GetLogginedForTest())
	})

	t.Run("異常系: requestURLがない状態でログアウトを試みる", func(t *testing.T) {
		// ログイン
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, client.GetLogginedForTest())

		originalRequestURL := client.GetRequestURLForTest() // 正しいURLを保持

		// requestURLをクリア
		client.SetRequestURLForTest("")

		// ログアウト  エラーになる
		err = client.Logout(context.Background())
		assert.Error(t, err)                        // requestURLがないのでError
		assert.True(t, client.GetLogginedForTest()) // ログアウトに失敗するので、true

		// ログアウト (後処理)
		defer func() {
			client.SetRequestURLForTest(originalRequestURL) // URLを戻す
			client.Logout(context.Background())             // log out
		}()
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/logout_test.go
