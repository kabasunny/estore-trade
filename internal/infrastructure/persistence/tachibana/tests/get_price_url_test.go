// internal/infrastructure/persistence/tachibana/tests/get_price_url_test.go
package tachibana_test

import (
	"testing"
	"time"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetPriceURL(t *testing.T) {
	t.Run("URL is cached", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{}

		// テスト用の値を設定
		tachibana.SetLogginedForTest(client, true)
		tachibana.SetExpiryForTest(client, time.Now().Add(1*time.Hour))
		tachibana.SetPriceURLForTest(client, "https://example.com/price") // Setter を使う

		url, err := client.GetPriceURL()
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/price", url)
	})

	t.Run("URL is not cached - not logged in", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{}
		url, err := client.GetPriceURL()
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.EqualError(t, err, "price URL not found, need to Login")
	})

	t.Run("URL is not cached - expired", func(t *testing.T) {
		client := &tachibana.TachibanaClientImple{}
		tachibana.SetLogginedForTest(client, true)
		tachibana.SetPriceURLForTest(client, "https://example.com/price") // Setter を使う
		tachibana.SetExpiryForTest(client, time.Now().Add(-1*time.Hour))

		url, err := client.GetPriceURL()
		assert.Error(t, err)
		assert.Equal(t, "", url)
		assert.EqualError(t, err, "price URL not found, need to Login")
	})
}
