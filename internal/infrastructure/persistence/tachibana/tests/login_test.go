// internal/infrastructure/persistence/tachibana/tests/login_test.go
package tachibana_test

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/config"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTachibanaClientImple_Login_Success(t *testing.T) {
	// httpmock をアクティブ化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// モックのレスポンスを設定 (正常系: ログイン成功)
	mockResponse := `{
        "sResultCode": "0",
        "sUrlRequest": "https://example.com/request",
        "sUrlMaster": "https://example.com/master",
        "sUrlPrice": "https://example.com/price",
        "sUrlEvent": "https://example.com/event",
        "p_no": "12345"
        }`
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(200, mockResponse))

	// テスト用の設定
	cfg := &config.Config{
		TachibanaBaseURL:  "https://example.com/",
		TachibanaUserID:   "testuser",
		TachibanaPassword: "testpassword",
	}
	logger, _ := zap.NewDevelopment() // テスト用のロガー

	// TachibanaClientImple のインスタンスを作成
	client := tachibana.NewTachibanaClient(cfg, logger)
	//p_noの初期値を確認
	tc := client.(*tachibana.TachibanaClientImple)
	assert.Equal(t, int64(0), tc.PNo)

	// Login メソッドを呼び出す
	err := client.Login(context.Background(), cfg)

	// 結果を検証
	assert.NoError(t, err)

	// キャッシュされたURLが正しく設定されているか確認
	assert.True(t, tc.Loggined)
	assert.Equal(t, "https://example.com/request", tc.RequestURL)
	assert.Equal(t, "https://example.com/master", tc.MasterURL)
	assert.Equal(t, "https://example.com/price", tc.PriceURL)
	assert.Equal(t, "https://example.com/event", tc.EventURL)
	assert.Equal(t, int64(12345), tc.PNo)
	assert.True(t, time.Now().Before(tc.Expiry))
}

// login_test.go
func TestTachibanaClientImple_Login_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// モックのレスポンスを設定 (異常系: ログイン失敗)
	mockResponse := `{"sResultCode": "E001", "sResultText": "Login failed"}`
	// ステータスコードを 401 (Unauthorized) に変更
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(401, mockResponse))

	cfg := &config.Config{
		TachibanaBaseURL:  "https://example.com/",
		TachibanaUserID:   "testuser",
		TachibanaPassword: "wrongpassword",
	}
	logger, _ := zap.NewDevelopment()

	client := tachibana.NewTachibanaClient(cfg, logger)
	//p_noの初期値は0
	tc := client.(*tachibana.TachibanaClientImple)
	assert.Equal(t, int64(0), tc.PNo)

	err := client.Login(context.Background(), cfg)

	// ログイン失敗時は、pNoが1になっているはず (TachibanaClientImple 側でインクリメント)
	assert.Equal(t, int64(1), tc.PNo)
	assert.ErrorContains(t, err, "login failed") // エラーメッセージの確認。ここを修正
	assert.False(t, tc.Loggined)                 //ログイン状態

	// URLは空文字列
	assert.Equal(t, "", tc.RequestURL)
	assert.Equal(t, "", tc.MasterURL)
	assert.Equal(t, "", tc.PriceURL)
	assert.Equal(t, "", tc.EventURL)

}

func TestTachibanaClientImple_Login_NetworkError(t *testing.T) {
	// httpmock.Activate() // Activate しないことで、ネットワークエラーを発生させる
	// (httpmock.DeactivateAndReset()は不要)

	cfg := &config.Config{
		TachibanaBaseURL:  "https://invalid.example.com/", // 存在しないURL
		TachibanaUserID:   "testuser",
		TachibanaPassword: "testpassword",
	}
	logger, _ := zap.NewDevelopment()
	client := tachibana.NewTachibanaClient(cfg, logger)
	//p_noの初期値は0
	tc := client.(*tachibana.TachibanaClientImple)
	assert.Equal(t, int64(0), tc.PNo)
	// ネットワークエラーが発生することを期待
	err := client.Login(context.Background(), cfg)
	// ログイン失敗時は、pNoが1になっているはず (TachibanaClientImple 側でインクリメント)
	assert.Equal(t, int64(1), tc.PNo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "login failed") // ネットワークエラー
	assert.False(t, tc.Loggined)                    //ログイン状態

	// URLは空文字列
	assert.Equal(t, "", tc.RequestURL)
	assert.Equal(t, "", tc.MasterURL)
	assert.Equal(t, "", tc.PriceURL)
	assert.Equal(t, "", tc.EventURL)
}

func TestTachibanaClientImple_Login_ExpiredURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 正常なレスポンスを返すモックを設定 (1回目)
	mockResponse1 := `{
        "sResultCode": "0",
        "sUrlRequest": "https://example.com/request1",
        "sUrlMaster": "https://example.com/master1",
        "sUrlPrice": "https://example.com/price1",
        "sUrlEvent": "https://example.com/event1",
        "p_no": "1"
        }`
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(200, mockResponse1))

	cfg := &config.Config{
		TachibanaBaseURL:  "https://example.com/",
		TachibanaUserID:   "testuser",
		TachibanaPassword: "testpassword",
	}
	logger, _ := zap.NewDevelopment()

	client := tachibana.NewTachibanaClient(cfg, logger)

	//p_noの初期値は0
	tc := client.(*tachibana.TachibanaClientImple)
	assert.Equal(t, int64(0), tc.PNo)
	// 1回目のログイン (成功)
	err := client.Login(context.Background(), cfg)
	assert.NoError(t, err)

	// URL がキャッシュされていることを確認
	assert.Equal(t, "https://example.com/request1", tc.RequestURL)
	//p_noが1になっていることを確認
	assert.Equal(t, int64(1), tc.PNo)

	// 有効期限を過去に設定
	tc.Expiry = time.Now().Add(-1 * time.Hour)

	// 正常なレスポンスを返すモックを設定 (2回目)
	mockResponse2 := `{
        "sResultCode": "0",
        "sUrlRequest": "https://example.com/request2",
        "sUrlMaster": "https://example.com/master2",
        "sUrlPrice": "https://example.com/price2",
        "sUrlEvent": "https://example.com/event2",
        "p_no": "2"
        }`
	// URL は同じだが、2回目の呼び出しとして認識させる
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(200, mockResponse2))

	// 2回目のログイン (キャッシュが無効なので、再度ログイン処理が行われる)
	err = client.Login(context.Background(), cfg)
	assert.NoError(t, err)

	// URL が更新されていることを確認
	assert.Equal(t, "https://example.com/request2", tc.RequestURL)
	//p_noが2になっていることを確認
	assert.Equal(t, int64(2), tc.PNo)
}
