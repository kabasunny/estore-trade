// internal/infrastructure/persistence/tachibana/get_event_url_test.go
package tachibana_test

import (
	"context"
	"fmt"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestGetEventURL(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)

	t.Run("正常系: ログイン後にEventURLを取得できること", func(t *testing.T) {
		// ログイン
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)

		// GetEventURL 呼び出し
		eventURL, err := client.GetEventURL()
		fmt.Println(eventURL)
		assert.NoError(t, err)
		assert.NotEmpty(t, eventURL) // 空でない文字列が返される

		// ログアウト
		err = client.Logout(context.Background())
		assert.NoError(t, err)
	})

	t.Run("異常系: ログイン前にEventURLを取得しようとするとエラーになること", func(t *testing.T) {
		// 	client := tachibana.CreateTestClient(t, nil) //ログインしていない状態

		// GetEventURL 呼び出し
		_, err := client.GetEventURL()
		assert.Error(t, err) // エラーが発生する

		// // ログアウト (ログインしていないので不要)
		// err = client.Logout(context.Background())
		// assert.NoError(t, err)
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/get_event_url_test.go
