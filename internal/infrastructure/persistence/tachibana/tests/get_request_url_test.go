// internal/infrastructure/persistence/tachibana/get_request_url_test.go
package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestURL(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)

	t.Run("正常系: ログイン後にRequestURLを取得できること", func(t *testing.T) {
		// ログイン
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)

		// GetRequestURL 呼び出し
		requestURL, err := client.GetRequestURL()
		assert.NoError(t, err)
		assert.NotEmpty(t, requestURL) // 空でない文字列が返される

		// ログアウト
		err = client.Logout(context.Background())
		assert.NoError(t, err)
	})

	t.Run("異常系: ログイン前にRequestURLを取得しようとするとエラーになること", func(t *testing.T) {
		//client := tachibana.CreateTestClient(t, nil) //ログインしていない状態

		// GetRequestURL 呼び出し
		_, err := client.GetRequestURL()
		assert.Error(t, err) // エラーが発生する

		// // ログアウト (ログインしていないので不要)
		// err = client.Logout(context.Background())
		// assert.NoError(t, err)
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/get_request_url_test.go
