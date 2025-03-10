// internal/infrastructure/persistence/tachibana/util_process_response.go
package tachibana

import (
	"encoding/json"
	"estore-trade/internal/domain"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func processResponse(body []byte, m *domain.MasterData, tc *TachibanaClientImple) error {
	// Shift-JIS から UTF-8 への変換
	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
	if err != nil {
		return fmt.Errorf("shift-jis decode error: %w", err)
	}

	// 連結されたJSONを分割
	jsonStr := string(bodyUTF8)
	// fmt.Println(jsonStr)
	jsonObjects := strings.Split(jsonStr, "}{") // "}{" で分割
	// 最初のオブジェクトに "{" を追加
	if len(jsonObjects) > 0 {
		jsonObjects[0] = jsonObjects[0] + "}" //最初
		if len(jsonObjects) > 1 {             //最後
			jsonObjects[len(jsonObjects)-1] = "{" + jsonObjects[len(jsonObjects)-1]
		}
	}
	if len(jsonObjects) > 2 {
		for i := 1; i < len(jsonObjects)-1; i++ { //途中
			jsonObjects[i] = "{" + jsonObjects[i] + "}"
		}
	}

	// 各JSONオブジェクトを処理
	for _, jsonObject := range jsonObjects {
		// fmt.Printf("Processing JSON object %d: %s\n", i+1, jsonObject) // デバッグ出力

		// JSON としてデコード
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(jsonObject), &item); err != nil {
			tc.logger.Warn("Failed to unmarshal line", zap.Error(err), zap.String("line", jsonObject))
			continue // デコードに失敗したらスキップ
		}

		// sCLMID キーの存在確認
		sCLMID, ok := item["sCLMID"].(string)
		if !ok {
			tc.logger.Warn("sCLMID not found in response item", zap.Any("item", item))
			continue
		}

		if sCLMID == "CLMEventDownloadComplete" {
			continue
		}
		// sCLMID の値に応じて処理
		switch sCLMID {
		case "CLMSystemStatus":
			var systemStatus domain.SystemStatus
			if err := mapToStruct(item, &systemStatus); err != nil {
				return fmt.Errorf("failed to map SystemStatus: %w", err)
			}
			m.SystemStatus = systemStatus
		case "CLMDateZyouhou":
			var dateInfo domain.DateInfo
			if err := mapToStruct(item, &dateInfo); err != nil {
				return fmt.Errorf("failed to map DateInfo: %w", err)
			}
			m.DateInfo = dateInfo
		case "CLMYobine":
			var callPrice domain.CallPrice
			if err := mapToStruct(item, &callPrice); err != nil {
				return fmt.Errorf("failed to map CallPrice: %w", err)
			}
			m.CallPriceMap[callPrice.UnitNumber] = callPrice
		case "CLMIssueMstKabu":
			var issueMaster domain.IssueMaster
			if err := mapToStruct(item, &issueMaster); err != nil {
				return fmt.Errorf("failed to map IssueMaster: %w", err)
			}
			m.IssueMap[issueMaster.IssueCode] = issueMaster
		case "CLMIssueSizyouMstKabu":
			var issueMarketMaster domain.IssueMarketMaster
			if err := mapToStruct(item, &issueMarketMaster); err != nil {
				return fmt.Errorf("failed to map IssueMarketMaster: %w", err)
			}
			if _, ok := m.IssueMarketMap[issueMarketMaster.IssueCode]; !ok {
				m.IssueMarketMap[issueMarketMaster.IssueCode] = make(map[string]domain.IssueMarketMaster)
			}
			m.IssueMarketMap[issueMarketMaster.IssueCode][issueMarketMaster.MarketCode] = issueMarketMaster

		case "CLMIssueSizyouKiseiKabu":
			var issueMarketRegulation domain.IssueMarketRegulation
			if err := mapToStruct(item, &issueMarketRegulation); err != nil {
				return fmt.Errorf("failed to map IssueMarketRegulation: %w", err)
			}
			if _, ok := m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode]; !ok {
				m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode] = make(map[string]domain.IssueMarketRegulation)
			}
			m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode][issueMarketRegulation.ListedMarket] = issueMarketRegulation
		case "CLMUnyouStatusKabu":
			var operationStatusKabu domain.OperationStatusKabu
			if err := mapToStruct(item, &operationStatusKabu); err != nil {
				return fmt.Errorf("failed to map OperationStatusKabu: %w", err)
			}
			if _, ok := m.OperationStatusKabuMap[operationStatusKabu.ListedMarket]; !ok {
				m.OperationStatusKabuMap[operationStatusKabu.ListedMarket] = make(map[string]domain.OperationStatusKabu)
			}
			m.OperationStatusKabuMap[operationStatusKabu.ListedMarket][operationStatusKabu.Unit] = operationStatusKabu
		default:
			// tc.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}
	return nil
}
