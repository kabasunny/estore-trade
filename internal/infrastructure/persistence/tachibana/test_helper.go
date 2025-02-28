// internal/infrastructure/persistence/tachibana/test_helper.go
package tachibana

import (
	"context"
	"testing"

	"estore-trade/internal/config"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
)

// SetupTestClient はテスト用の TachibanaClientImple と Config を作成
func SetupTestClient(t *testing.T) (*TachibanaClientImple, *config.Config) {
	httpmock.Activate()

	cfg := &config.Config{
		TachibanaBaseURL:   "https://example.com/", // モックするのでダミーでOK
		TachibanaUserID:    "testuser",
		TachibanaPassword:  "testpassword",
		TachibanaAPIKey:    "testapikey", // モックなので使われない
		TachibanaAPISecret: "testapisecret",
	}
	logger, _ := zap.NewDevelopment() // テスト用のロガー

	client := NewTachibanaClient(cfg, logger).(*TachibanaClientImple)

	// DownloadMasterData メソッドのモック (より完全なデータ、API 仕様に準拠)
	//必要なデータのみ追加
	mockMasterResponse := map[string]interface{}{
		"sResultCode": "0",
		"CLMIssueMstKabu": []interface{}{ // 配列に変更
			map[string]interface{}{"sCLMID": "CLMIssueMstKabu", "sIssueCode": "7974", "sIssueName": "任天堂"},
			map[string]interface{}{"sCLMID": "CLMIssueMstKabu", "sIssueCode": "9984", "sIssueName": "ソフトバンク"},
		},
		"CLMIssueSizyouMstKabu": []interface{}{
			map[string]interface{}{"sCLMID": "CLMIssueSizyouMstKabu", "sIssueCode": "7974", "sZyouzyouSizyou": "00"},
			map[string]interface{}{"sCLMID": "CLMIssueSizyouMstKabu", "sIssueCode": "9984", "sZyouzyouSizyou": "00"},
		},
		//CLMIssueSizyouKiseiKabuを追加
		"CLMIssueSizyouKiseiKabu": []interface{}{
			map[string]interface{}{"sCLMID": "CLMIssueSizyouKiseiKabu", "sSystemKouzaKubun": "102", "sIssueCode": "7974", "sZyouzyouSizyou": "00", "sTeisiKubun": "1"},
			map[string]interface{}{"sCLMID": "CLMIssueSizyouKiseiKabu", "sSystemKouzaKubun": "102", "sIssueCode": "9984", "sZyouzyouSizyou": "00", "sTeisiKubun": "1"},
		},
		"CLMUnyouStatusKabu": []interface{}{
			map[string]interface{}{"sCLMID": "CLMUnyouStatusKabu", "sZyouzyouSizyou": "00", "sUnyouUnit": "0101"},
		},
		// 他のマスタデータもモック (必要に応じて)
		"CLMSystemStatus": map[string]interface{}{
			"sCLMID":           "CLMSystemStatus",
			"sSystemStatusKey": "001",
			"sLoginKyokaKubun": "1",
			"sSystemStatus":    "1",
		},
		"CLMDateZyouhou": map[string]interface{}{
			"sCLMID":  "CLMDateZyouhou",
			"sDayKey": "001",
			"sTheDay": "20231101",
		},
		"CLMYobine": []interface{}{ // 配列に変更
			map[string]interface{}{"sCLMID": "CLMYobine", "sYobineTaniNumber": "101"},
			map[string]interface{}{"sCLMID": "CLMYobine", "sYobineTaniNumber": "102"},
		},
		"CLMEventDownloadComplete": map[string]interface{}{}, // ダウンロード完了通知
	}

	httpmock.RegisterResponder("POST", "https://example.com/master", // MasterURL
		httpmock.NewJsonResponderOrPanic(200, mockMasterResponse),
	)
	// LoginをMockする（RequestURLなどを設定するため)
	loginMockResponse := `{
        "sResultCode": "0",
        "sUrlRequest": "https://example.com/request",
        "sUrlMaster": "https://example.com/master",
        "sUrlPrice": "https://example.com/price",
        "sUrlEvent": "https://example.com/event",
        "p_no": "12345"
        }`
	httpmock.RegisterResponder("POST", "https://example.com/login",
		httpmock.NewStringResponder(200, loginMockResponse))

	client.Login(context.Background(), cfg) // Login を実行して RequestURL などを設定

	// DownloadMasterData を呼び出して、モックデータがセットされることを確認
	_, err := client.DownloadMasterData(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return client, cfg
}
