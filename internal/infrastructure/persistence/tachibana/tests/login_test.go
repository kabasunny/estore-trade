// internal/infrastructure/persistence/tachibana/tests/login_test.go
package tachibana_test

import (
	"context"
	"estore-trade/internal/domain" // 追加
	"testing"
	"time"

	"estore-trade/internal/config"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTachibanaClientImple_Login_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	cfg := &config.Config{
		TachibanaBaseURL:  "https://example.com/",
		TachibanaUserID:   "testuser",
		TachibanaPassword: "testpassword",
	}
	logger, _ := zap.NewDevelopment()

	// モックの MasterData を作成
	mockMasterData := &domain.MasterData{}

	client := tachibana.NewTachibanaClient(cfg, logger, mockMasterData).(*tachibana.TachibanaClientImple) // masterData を渡す

	assert.Equal(t, int64(0), tachibana.GetPNoForTest(client)) // 関数を使用

	err := client.Login(context.Background(), cfg)
	assert.NoError(t, err)

	assert.True(t, tachibana.IsLogginedForTest(client))                                    // 関数を使用
	assert.Equal(t, "https://example.com/request", tachibana.GetRequestURLForTest(client)) // 関数
	assert.Equal(t, "https://example.com/master", tachibana.GetMasterURLForTest(client))   // 関数
	assert.Equal(t, "https://example.com/price", tachibana.GetPriceURLForTest(client))     // 関数
	assert.Equal(t, "https://example.com/event", tachibana.GetEventURLForTest(client))     // 関数
	assert.Equal(t, int64(12345), tachibana.GetPNoForTest(client))                         // 関数
	assert.True(t, time.Now().Before(tachibana.GetExpiryForTest(client)))                  // 関数
}

func TestTachibanaClientImple_Login_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResponse := `{"sResultCode": "E001", "sResultText": "Login failed"}`
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(401, mockResponse))

	cfg := &config.Config{
		TachibanaBaseURL:  "https://example.com/",
		TachibanaUserID:   "testuser",
		TachibanaPassword: "wrongpassword",
	}
	logger, _ := zap.NewDevelopment()

	// モックの MasterData を作成
	mockMasterData := &domain.MasterData{}

	client := tachibana.NewTachibanaClient(cfg, logger, mockMasterData).(*tachibana.TachibanaClientImple) // masterDataを渡す
	assert.Equal(t, int64(0), tachibana.GetPNoForTest(client))                                            // 関数を使用

	err := client.Login(context.Background(), cfg)

	assert.Equal(t, int64(1), tachibana.GetPNoForTest(client))  // 関数を使用
	assert.ErrorContains(t, err, "login failed")                // エラーメッセージ (APIエラー)
	assert.False(t, tachibana.IsLogginedForTest(client))        // 関数を使用
	assert.Equal(t, "", tachibana.GetRequestURLForTest(client)) // 関数を使用
	assert.Equal(t, "", tachibana.GetMasterURLForTest(client))  // 関数を使用
	assert.Equal(t, "", tachibana.GetPriceURLForTest(client))   // 関数を使用
	assert.Equal(t, "", tachibana.GetEventURLForTest(client))   // 関数を使用
}

func TestTachibanaClientImple_Login_NetworkError(t *testing.T) {

	cfg := &config.Config{
		TachibanaBaseURL:  "https://invalid.example.com/", // 存在しないURL
		TachibanaUserID:   "testuser",
		TachibanaPassword: "testpassword",
	}
	logger, _ := zap.NewDevelopment()
	// モックの MasterData を作成
	mockMasterData := &domain.MasterData{}

	client := tachibana.NewTachibanaClient(cfg, logger, mockMasterData).(*tachibana.TachibanaClientImple) //masterDataを渡す
	assert.Equal(t, int64(0), tachibana.GetPNoForTest(client))                                            // 関数を使用
	err := client.Login(context.Background(), cfg)
	assert.Equal(t, int64(1), tachibana.GetPNoForTest(client))  // 関数を使用
	assert.Error(t, err)                                        // エラー
	assert.Contains(t, err.Error(), "login failed")             // エラーメッセージ (ネットワークエラー)
	assert.False(t, tachibana.IsLogginedForTest(client))        // 関数を使用
	assert.Equal(t, "", tachibana.GetRequestURLForTest(client)) // 関数を使用
	assert.Equal(t, "", tachibana.GetMasterURLForTest(client))  // 関数を使用
	assert.Equal(t, "", tachibana.GetPriceURLForTest(client))   // 関数を使用
	assert.Equal(t, "", tachibana.GetEventURLForTest(client))   // 関数を使用
}

func TestTachibanaClientImple_Login_ExpiredURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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
	// モックの MasterData を作成
	mockMasterData := &domain.MasterData{}

	client := tachibana.NewTachibanaClient(cfg, logger, mockMasterData).(*tachibana.TachibanaClientImple) //masterDataを渡す

	assert.Equal(t, int64(0), tachibana.GetPNoForTest(client)) // 関数を使用
	err := client.Login(context.Background(), cfg)
	assert.NoError(t, err)

	assert.Equal(t, "https://example.com/request1", tachibana.GetRequestURLForTest(client)) // 関数を使用
	assert.Equal(t, int64(1), tachibana.GetPNoForTest(client))                              // 関数を使用

	// 有効期限を過去に設定 (テストヘルパーの関数を使用)
	tachibana.SetExpiryForTest(client, time.Now().Add(-1*time.Hour))

	mockResponse2 := `{
        "sResultCode": "0",
        "sUrlRequest": "https://example.com/request2",
        "sUrlMaster": "https://example.com/master2",
        "sUrlPrice": "https://example.com/price2",
        "sUrlEvent": "https://example.com/event2",
        "p_no": "2"
        }`
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(200, mockResponse2))

	err = client.Login(context.Background(), cfg)
	assert.NoError(t, err)

	assert.Equal(t, "https://example.com/request2", tachibana.GetRequestURLForTest(client)) // 関数を使用
	assert.Equal(t, int64(2), tachibana.GetPNoForTest(client))                              // 関数を使用
}
