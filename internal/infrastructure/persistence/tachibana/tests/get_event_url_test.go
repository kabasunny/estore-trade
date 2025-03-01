// internal/infrastructure/persistence/tachibana/tests/get_event_url_test.go
package tachibana_test

import (
	"testing"
	"time"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetEventURL(t *testing.T) {
	t.Run("URL is cached", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{} // 直接インスタンスを生成

		// テスト用の値を設定
		tachibana.SetLogginedForTest(client, true)
		tachibana.SetExpiryForTest(client, time.Now().Add(1*time.Hour))
		tachibana.SetEventURLForTest(client, "https://example.com/event") // Setter を使う

		url, err := client.GetEventURL()
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/event", url)
	})

	t.Run("URL is not cached - not logged in", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{}
		url, err := client.GetEventURL()
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.EqualError(t, err, "event URL not found, need to Login")
	})

	t.Run("URL is not cached - expired", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{}
		tachibana.SetLogginedForTest(client, true)
		tachibana.SetEventURLForTest(client, "https://example.com/event") // Setter を使う
		tachibana.SetExpiryForTest(client, time.Now().Add(-1*time.Hour))

		url, err := client.GetEventURL()
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.EqualError(t, err, "event URL not found, need to Login")
	})
}
