// internal/infrastructure/persistence/tachibana/test_helper.go
package tachibana

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/config"
	"estore-trade/internal/domain"

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
		//EventRid, EventBoardNo, EventEvtCmd を追加
		EventRid:     "testrid",
		EventBoardNo: "testboardno",
		EventEvtCmd:  "EC,TE", // カンマ区切りで複数指定可能
	}
	logger, _ := zap.NewDevelopment() // テスト用のロガー
	// モックの MasterData を作成
	mockMasterData := &domain.MasterData{}

	client := NewTachibanaClient(cfg, logger, mockMasterData).(*TachibanaClientImple) // masterData を渡す

	// DownloadMasterData メソッドのモック (より完全なデータ、API 仕様に準拠)
	//必要なデータのみ追加 + CallPrice関連の値追加
	mockMasterResponse := map[string]interface{}{
		"sResultCode": "0", // 0:正常
		"CLMIssueMstKabu": []interface{}{ // 配列に変更
			map[string]interface{}{"sCLMID": "CLMIssueMstKabu", "sIssueCode": "7974", "sIssueName": "任天堂"},
			map[string]interface{}{"sCLMID": "CLMIssueMstKabu", "sIssueCode": "9984", "sIssueName": "ソフトバンク"},
		},
		"CLMIssueSizyouMstKabu": []interface{}{
			// 7974と9984で異なる呼値単位番号が設定されるように修正
			map[string]interface{}{"sCLMID": "CLMIssueSizyouMstKabu", "sIssueCode": "7974", "sZyouzyouSizyou": "00", "sYobineTaniNumber": "101", "sYobineTaniNumberYoku": "101"}, //テストに必要な値を追加
			map[string]interface{}{"sCLMID": "CLMIssueSizyouMstKabu", "sIssueCode": "9984", "sYobineTaniNumber": "101", "sYobineTaniNumberYoku": "101"},                          //テストに必要な値を追加
		},
		//CLMIssueSizyouKiseiKabuを追加
		"CLMIssueSizyouKiseiKabu": []interface{}{
			map[string]interface{}{"sCLMID": "CLMIssueSizyouKiseiKabu", "sSystemKouzaKubun": "102", "sIssueCode": "7974", "sZyouzyouSizyou": "00", "sTeisiKubun": "1"},
			map[string]interface{}{"sCLMID": "CLMIssueSizyouKiseiKabu", "sSystemKouzaKubun": "102", "sIssueCode": "9984", "sZyouzyouSizyou": "00", "sTeisiKubun": "1"},
		},
		"CLMUnyouStatusKabu": []interface{}{
			map[string]interface{}{"sCLMID": "CLMUnyouStatusKabu", "sZyouzyouSizyou": "00", "sUnyouUnit": "0101", "sUnyouStatus": "001"}, //Status追加
		},
		// 他のマスタデータもモック (必要に応じて)
		"CLMSystemStatus": map[string]interface{}{
			"sCLMID":           "CLMSystemStatus",
			"sSystemStatusKey": "001",
			"sLoginKyokaKubun": "1",
			"sSystemStatus":    "1",
		},
		"CLMDateZyouhou": map[string]interface{}{
			"sCLMID":                "CLMDateZyouhou",
			"sDayKey":               "001",
			"sMaeEigyouDay_1":       "20231031",
			"sMaeEigyouDay_2":       "20231030",
			"sMaeEigyouDay_3":       "20231027",
			"sTheDay":               "20231101",
			"sYokuEigyouDay_1":      "20231102",
			"sYokuEigyouDay_2":      "20231106",
			"sYokuEigyouDay_3":      "20231107",
			"sYokuEigyouDay_4":      "20231108",
			"sYokuEigyouDay_5":      "20231109",
			"sYokuEigyouDay_6":      "20231110",
			"sYokuEigyouDay_7":      "20231113",
			"sYokuEigyouDay_8":      "20231114",
			"sYokuEigyouDay_9":      "20231115",
			"sYokuEigyouDay_10":     "20231116",
			"sKabuUkewatasiDay":     "20231106",
			"sKabuKariUkewatasiDay": "20231107", // 追加
			"sBondUkewatasiDay":     "20231106", // 追加
		},
		"CLMYobine": []interface{}{ // 配列に変更
			// sYobineTaniNumber=101 (当日) の呼値テーブル
			map[string]interface{}{
				"sCLMID":            "CLMYobine",
				"sYobineTaniNumber": "101",
				"sTekiyouDay":       "20140101", // ダミーの値
				"sKizunPrice_1":     "3000.0",
				"sKizunPrice_2":     "5000.0",
				"sKizunPrice_3":     "30000.0",
				"sKizunPrice_4":     "50000.0",
				"sKizunPrice_5":     "300000.0",
				"sKizunPrice_6":     "500000.0",
				"sKizunPrice_7":     "3000000.0",
				"sKizunPrice_8":     "5000000.0",
				"sKizunPrice_9":     "30000000.0",
				"sKizunPrice_10":    "50000000.0",
				"sKizunPrice_11":    "999999999.0",
				"sKizunPrice_12":    "0.0",
				"sKizunPrice_13":    "0.0",
				"sKizunPrice_14":    "0.0",
				"sKizunPrice_15":    "0.0",
				"sKizunPrice_16":    "0.0",
				"sKizunPrice_17":    "0.0",
				"sKizunPrice_18":    "0.0",
				"sKizunPrice_19":    "0.0",
				"sKizunPrice_20":    "0.0",
				"sYobineTanka_1":    "1.0",
				"sYobineTanka_2":    "5.0",
				"sYobineTanka_3":    "10.0",
				"sYobineTanka_4":    "50.0",
				"sYobineTanka_5":    "100.0",
				"sYobineTanka_6":    "500.0",
				"sYobineTanka_7":    "1000.0",
				"sYobineTanka_8":    "5000.0",
				"sYobineTanka_9":    "10000.0",
				"sYobineTanka_10":   "50000.0",
				"sYobineTanka_11":   "100000.0",
				"sYobineTanka_12":   "0.0",
				"sYobineTanka_13":   "0.0",
				"sYobineTanka_14":   "0.0",
				"sYobineTanka_15":   "0.0",
				"sYobineTanka_16":   "0.0",
				"sYobineTanka_17":   "0.0",
				"sYobineTanka_18":   "0.0",
				"sYobineTanka_19":   "0.0",
				"sYobineTanka_20":   "0.0",
				//sDecimal_1 など、,stringがついていないものは、あってもなくても良い
			},
		},
		"CLMEventDownloadComplete": map[string]interface{}{}, // ダウンロード完了通知
	}
	//モック追加
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

	// DownloadMasterData のモックの代わりに、モックデータを直接設定
	masterData := &domain.MasterData{
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMap:                 make(map[string]domain.IssueMaster),
		IssueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		OperationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
	}

	err := processResponse(mockMasterResponse, masterData, client)
	if err != nil {
		t.Fatalf("failed to process mock master response: %v", err)
	}

	// モックデータをクライアントに設定
	client.masterData = masterData
	client.masterData.SystemStatus = masterData.SystemStatus
	client.masterData.DateInfo = masterData.DateInfo
	client.masterData.CallPriceMap = masterData.CallPriceMap
	client.masterData.IssueMap = masterData.IssueMap
	client.masterData.IssueMarketMap = masterData.IssueMarketMap
	client.masterData.IssueMarketRegulationMap = masterData.IssueMarketRegulationMap
	client.masterData.OperationStatusKabuMap = masterData.OperationStatusKabuMap
	client.targetIssueCodes = []string{"7974", "9984"} //実装において、このフィールドを追加するロジックが必要

	return client, cfg
}

// SetTargetIssueCodesForTest はテスト用に targetIssueCodes を設定する関数
func (tc *TachibanaClientImple) SetTargetIssueCodesForTest(issueCodes []string) {
	tc.targetIssueCodesMu.Lock() //ミューテックスでロック
	defer tc.targetIssueCodesMu.Unlock()
	tc.targetIssueCodes = issueCodes
}

// SetCallPriceMapForTest はテスト用に callPriceMap を設定する関数です。
func (tc *TachibanaClientImple) SetCallPriceMapForTest(m map[string]domain.CallPrice) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData.CallPriceMap = m
}

// SetIssueMarketMapForTest はテスト用に issueMarketMap を設定する関数
func (tc *TachibanaClientImple) SetIssueMarketMapForTest(m map[string]map[string]domain.IssueMarketMaster) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData.IssueMarketMap = m
}

// SetIssueMarketRegulationMapForTest はテスト用に issueMarketRegulationMap を設定する関数
func (tc *TachibanaClientImple) SetIssueMarketRegulationMapForTest(m map[string]map[string]domain.IssueMarketRegulation) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData.IssueMarketRegulationMap = m
}

// SetIssueMapForTestはテスト用に issueMap を設定する関数
func (tc *TachibanaClientImple) SetIssueMapForTest(m map[string]domain.IssueMaster) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData.IssueMap = m
}

// SetOperationStatusKabuMapForTest はテスト用に operationStatusKabuMap を設定する関数
func (tc *TachibanaClientImple) SetOperationStatusKabuMapForTest(m map[string]map[string]domain.OperationStatusKabu) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData.OperationStatusKabuMap = m
}

// 以下、テスト用の Getter 関数

// GetPNoForTest はテスト用に pNo を返す
func GetPNoForTest(tc *TachibanaClientImple) int64 {
	tc.pNoMu.Lock()
	defer tc.pNoMu.Unlock()
	return tc.pNo
}

// IsLogginedForTest はテスト用に loggined の状態を返す
func IsLogginedForTest(tc *TachibanaClientImple) bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loggined
}

// SetLogginedForTest はテスト用にlogginedを設定する
func SetLogginedForTest(tc *TachibanaClientImple, loggined bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.loggined = loggined
}

// GetRequestURLForTest はテスト用に requestURL を返す
func GetRequestURLForTest(tc *TachibanaClientImple) string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.requestURL
}

// SetRequestURLForTest はテスト用にrequestURLを設定する
func SetRequestURLForTest(tc *TachibanaClientImple, requestURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.requestURL = requestURL
}

// GetMasterURLForTest はテスト用に masterURL を返す
func GetMasterURLForTest(tc *TachibanaClientImple) string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterURL
}

// GetPriceURLForTest はテスト用に priceURL を返す
func GetPriceURLForTest(tc *TachibanaClientImple) string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.priceURL
}

// GetEventURLForTest はテスト用に eventURL を返す
func GetEventURLForTest(tc *TachibanaClientImple) string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.eventURL
}

// SetEventURLForTest はテスト用に eventURL を設定する関数
func SetEventURLForTest(tc *TachibanaClientImple, eventURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.eventURL = eventURL
}

// GetExpiryForTest はテスト用に expiryを返す
func GetExpiryForTest(tc *TachibanaClientImple) time.Time {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.expiry
}

// SetExpiryForTest はテスト用に expiry を設定する関数
func SetExpiryForTest(tc *TachibanaClientImple, t time.Time) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.expiry = t
}

// GetSystemStatusForTest はテスト用に systemStatus を返す関数
func GetSystemStatusForTest(tc *TachibanaClientImple) domain.SystemStatus {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterData.SystemStatus
}

// GetDateInfoForTest はテスト用に dateInfo を返す関数
func GetDateInfoForTest(tc *TachibanaClientImple) domain.DateInfo {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterData.DateInfo
}

// SetDateInfoForTest はテスト用に dateInfo を設定する関数
func SetDateInfoForTest(tc *TachibanaClientImple, dateInfo domain.DateInfo) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData.DateInfo = dateInfo
}

// GetMasterDataForTest はテスト用に masterData を返す関数
func GetMasterDataForTest(tc *TachibanaClientImple) *domain.MasterData {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterData
}

// SetMasterDataForTest はテスト用に masterData を設定する関数
func SetMasterDataForTest(tc *TachibanaClientImple, masterData *domain.MasterData) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterData = masterData
}

// SetMasterURLForTest はテスト用に masterURL を設定する関数
func SetMasterURLForTest(tc *TachibanaClientImple, masterURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.masterURL = masterURL
}

// SetPriceURLForTest はテスト用に priceURL を設定する関数
func SetPriceURLForTest(tc *TachibanaClientImple, priceURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.priceURL = priceURL
}

// ParseEventForTest はテスト用に EventStream.parseEvent を呼び出すヘルパー関数
func ParseEventForTest(es *EventStream, message []byte) (*domain.OrderEvent, error) {
	return es.parseEvent(message)
}

// SendEventForTest はテスト用に EventStream.sendEvent を呼び出すヘルパー関数
func SendEventForTest(es *EventStream, event *domain.OrderEvent) {
	es.sendEvent(event)
}
