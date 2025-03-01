// internal/infrastructure/persistence/tachibana/tests/get_request_url_test.go
package tachibana_test

import (
	"testing"
	"time"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetRequestURL(t *testing.T) {
	// test_helper.go の SetupTestClient は Login 処理と MasterData のモックを行うため、
	// 今回は直接 TachibanaClientImple のインスタンスを生成する

	t.Run("URL is cached", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{} // 直接インスタンスを生成

		// テスト用の値を設定 (モックの代わり)
		tachibana.SetLogginedForTest(client, true) //login状態
		tachibana.SetRequestURLForTest(client, "https://example.com/request")
		tachibana.SetExpiryForTest(client, time.Now().Add(1*time.Hour)) // 有効期限を未来に設定

		url, err := client.GetRequestURL()
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/request", url)
	})

	t.Run("URL is not cached - not logged in", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{} // 直接インスタンスを生成
		// Loggined = false (デフォルト) なので、ログインしていない状態
		url, err := client.GetRequestURL()
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.EqualError(t, err, "request URL not found, need to Login")
	})

	t.Run("URL is not cached - expired", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{}
		tachibana.SetLogginedForTest(client, true) //login状態
		tachibana.SetRequestURLForTest(client, "https://example.com/request")
		tachibana.SetExpiryForTest(client, time.Now().Add(-1*time.Hour)) // 有効期限切れ

		url, err := client.GetRequestURL()
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.EqualError(t, err, "request URL not found, need to Login")
	})
}
