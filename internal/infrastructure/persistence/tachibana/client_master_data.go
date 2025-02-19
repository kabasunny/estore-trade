// internal/infrastructure/persistence/tachibana/client_master_data.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// マスタデータ保持用 (必要最低限に絞り込み)
type masterDataManager struct {
	systemStatus             SystemStatus
	dateInfo                 DateInfo
	callPriceMap             map[string]CallPrice                    // 呼値 (Key: sYobineTaniNumber)
	issueMap                 map[string]IssueMaster                  // 銘柄マスタ（株式）(Key: 銘柄コード)
	issueMarketMap           map[string]map[string]IssueMarketMaster // 株式銘柄市場マスタ (Key1: 銘柄コード, Key2: 上場市場)
	issueMarketRegulationMap map[string]map[string]IssueMarketRegulation
	operationStatusKabuMap   map[string]map[string]OperationStatusKabu // 運用ステータス（株）(Key1: 上場市場, Key2: 運用単位)
}

func (tc *TachibanaClientImple) DownloadMasterData(ctx context.Context) error {
	// sTargetCLMID を使用して、必要なマスタデータのみを要求
	payload := map[string]string{
		"sCLMID":       clmidDownloadMasterData,
		"sTargetCLMID": "CLMSystemStatus,CLMDateZyouhou,CLMYobine,CLMIssueMstKabu,CLMIssueSizyouMstKabu,CLMIssueSizyouKiseiKabu,CLMUnyouStatusKabu,CLMEventDownloadComplete", // 必要なマスタデータを指定
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal master data request payload: %w", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.masterURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create master data request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// コンテキストとタイムアウトの設定
	// 修正: withContextAndTimeoutの戻り値でreqを上書き
	req, cancel := withContextAndTimeout(req, 600*time.Second) // マスタデータダウンロードは時間がかかる可能性があるため、長めに設定
	defer cancel()

	// リトライ処理 sendRequest関数に統一  //client := &http.Client{}を削除
	response, err := sendRequest(ctx, tc, req) // reqを渡す
	if err != nil {
		return fmt.Errorf("failed to download master data: %w", err)
	}

	// マスタデータマネージャーの初期化
	tc.mu.Lock()
	m := &masterDataManager{
		callPriceMap:             make(map[string]CallPrice),
		issueMap:                 make(map[string]IssueMaster),
		issueMarketMap:           make(map[string]map[string]IssueMarketMaster),
		issueMarketRegulationMap: make(map[string]map[string]IssueMarketRegulation),
		operationStatusKabuMap:   make(map[string]map[string]OperationStatusKabu),
	}
	tc.mu.Unlock()

	// レスポンスからsCLMIDを取り出す
	for sCLMID, data := range response {
		// sCLMID でどのマスタデータか判別
		// dataをmap[string]interface{}に変換
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			tc.logger.Error("Invalid data format in master data response", zap.String("sCLMID", sCLMID))
			continue // 処理をスキップして次のデータへ
		}
		switch sCLMID {
		case "CLMSystemStatus":
			var systemStatus SystemStatus
			if err := mapToStruct(dataMap, &systemStatus); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map SystemStatus: %w", err)
			}
			tc.mu.Lock()
			m.systemStatus = systemStatus
			tc.mu.Unlock()
		case "CLMDateZyouhou":
			var dateInfo DateInfo
			if err := mapToStruct(dataMap, &dateInfo); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map DateInfo: %w", err)
			}
			tc.mu.Lock()
			m.dateInfo = dateInfo
			tc.mu.Unlock()
		case "CLMYobine":
			var callPrice CallPrice
			if err := mapToStruct(dataMap, &callPrice); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map CallPrice: %w", err)
			}
			tc.mu.Lock()
			m.callPriceMap[strconv.Itoa(callPrice.UnitNumber)] = callPrice
			tc.mu.Unlock()

		case "CLMIssueMstKabu":
			var issueMaster IssueMaster
			if err := mapToStruct(dataMap, &issueMaster); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map IssueMaster: %w", err)
			}
			tc.mu.Lock()
			m.issueMap[issueMaster.IssueCode] = issueMaster
			tc.mu.Unlock()

		case "CLMIssueSizyouMstKabu":
			var issueMarket IssueMarketMaster
			if err := mapToStruct(dataMap, &issueMarket); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map IssueMarketMaster: %w", err)
			}

			tc.mu.Lock()
			if _, ok := m.issueMarketMap[issueMarket.IssueCode]; !ok {
				m.issueMarketMap[issueMarket.IssueCode] = make(map[string]IssueMarketMaster)
			}
			m.issueMarketMap[issueMarket.IssueCode][issueMarket.MarketCode] = issueMarket
			tc.mu.Unlock()

		case "CLMUnyouStatusKabu":
			var operationStatusKabu OperationStatusKabu
			if err := mapToStruct(dataMap, &operationStatusKabu); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map OperationStatusKabu: %w", err)
			}
			tc.mu.Lock()
			if _, ok := m.operationStatusKabuMap[operationStatusKabu.ListedMarket]; !ok {
				m.operationStatusKabuMap[operationStatusKabu.ListedMarket] = make(map[string]OperationStatusKabu)
			}
			m.operationStatusKabuMap[operationStatusKabu.ListedMarket][operationStatusKabu.Unit] = operationStatusKabu
			tc.mu.Unlock()

		case "CLMIssueSizyouKiseiKabu":
			var issueMarketRegulation IssueMarketRegulation
			if err := mapToStruct(dataMap, &issueMarketRegulation); err != nil { // dataをdataMapに変更
				return fmt.Errorf("failed to map IssueMarketRegulation: %w", err)
			}
			tc.mu.Lock()
			if _, ok := m.issueMarketRegulationMap[issueMarketRegulation.IssueCode]; !ok {
				m.issueMarketRegulationMap[issueMarketRegulation.IssueCode] = make(map[string]IssueMarketRegulation)
			}
			m.issueMarketRegulationMap[issueMarketRegulation.IssueCode][issueMarketRegulation.ListedMarket] = issueMarketRegulation
			tc.mu.Unlock()

		case "CLMEventDownloadComplete":
			tc.mu.Lock()
			tc.systemStatus = m.systemStatus
			tc.dateInfo = m.dateInfo
			tc.issueMap = m.issueMap
			tc.callPriceMap = m.callPriceMap
			tc.issueMarketMap = m.issueMarketMap
			tc.issueMarketRegulationMap = m.issueMarketRegulationMap
			tc.operationStatusKabuMap = m.operationStatusKabuMap
			tc.mu.Unlock()
			return nil

		default:
			tc.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}

	// ダウンロード完了後、ターゲット銘柄が設定されていればフィルタリング
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// m から tc にデータをコピー (一旦、全データをコピー)
	tc.systemStatus = m.systemStatus
	tc.dateInfo = m.dateInfo
	tc.issueMap = m.issueMap
	tc.callPriceMap = m.callPriceMap
	tc.issueMarketMap = m.issueMarketMap
	tc.issueMarketRegulationMap = m.issueMarketRegulationMap
	tc.operationStatusKabuMap = m.operationStatusKabuMap

	if len(tc.targetIssueCodes) > 0 {
		if err := tc.SetTargetIssues(ctx, tc.targetIssueCodes); err != nil { // ターゲット銘柄でフィルタリング
			return err
		}
	}

	return nil
}

// mapToStruct は、map[string]interface{} を構造体にマッピングする汎用関数
func mapToStruct(data map[string]interface{}, result interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, result)
}
