// internal/infrastructure/persistence/tachibana/login_test.go
package tachibana_test

import (
	"context"

	//"estore-trade/internal/domain" // 今回は不要
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// テスト用のMasterDataを作成  今回は不要
	//md := &domain.MasterData{} // 必要に応じてダミーデータを設定
	client := tachibana.CreateTestClient(t, nil)

	t.Run("正常系: 正しいIDとパスワードでログインできること", func(t *testing.T) {
		// client := tachibana.CreateTestClient(t, nil) //clientを使いまわすと、logginedがtrueになるので、分ける
		// Login
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)

		// ログイン状態を確認 (loggined フラグ、requestURL など)
		assert.True(t, client.GetLogginedForTest())
		requestURL := client.GetRequestURLForTest() //キャッシュされたrequestURLを取得
		assert.NotEmpty(t, requestURL)

		// ログアウト (後処理)
		err = client.Logout(context.Background())
		assert.NoError(t, err)
	})

	t.Run("異常系: 不正なIDとパスワードでログインできないこと", func(t *testing.T) {
		// 不正な認証情報でログインを試みる
		originalUserID := client.GetUserIDForTest()     //userID退避
		originalPassword := client.GetPasswordForTest() //password退避

		client.SetUserIDForTest("invalid_user")       // 存在しないユーザ
		client.SetPasswordForTest("invalid_password") // 間違ったパスワード
		err := client.Login(context.Background(), nil)
		assert.Error(t, err)

		// ログイン状態が false であることを確認
		assert.False(t, client.GetLogginedForTest())
		requestURL := client.GetRequestURLForTest() //キャッシュされたrequestURLを取得
		assert.Empty(t, requestURL)                 //requestURLがないはずなのでEmpty

		defer func() {
			client.SetUserIDForTest(originalUserID)     // UserIDを戻す
			client.SetPasswordForTest(originalPassword) // Passwordを戻す
		}()
	})

	// t.Run("異常系: APIがエラーレスポンスを返す場合", func(t *testing.T) {
	// 	// .env の TachibanaBaseURL を無効なURLに変更してテスト (意図的にAPIエラーを発生させる)
	// 	originalBaseURL := client.GetBaseURLForTest()            // オリジナルの baseURL を取得
	// 	client.SetBaseURLForTest("https://invalid.example.com/") // 無効な URL を設定
	// 	err := client.Login(context.Background(), nil)
	// 	assert.Error(t, err)
	// 	assert.False(t, client.GetLogginedForTest())
	// 	// エラーメッセージの検証 (任意)
	// 	// logout
	// 	defer func() {
	// 		client.SetBaseURLForTest(originalBaseURL) // URLを戻す
	// 		client.Logout(context.Background())       // log out
	// 	}()
	// })
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/login_test.go
