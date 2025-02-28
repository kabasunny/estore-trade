// internal/infrastructure/persistence/tachibana/util_process_response.go
package tachibana

import (
	"estore-trade/internal/domain" //追加
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

func processResponse(response map[string]interface{}, m *domain.MasterData, tc *TachibanaClientImple) error { //m *masterDataManagerを*domain.MasterData
	// レスポンスからsCLMIDを取り出す
	for sCLMID, data := range response {
		// sCLMID でどのマスタデータか判別
		// dataをmap[string]interface{}に変換
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			tc.Logger.Error("Invalid data format in master data response", zap.String("sCLMID", sCLMID))
			continue // 処理をスキップして次のデータへ
		}
		switch sCLMID {
		case "CLMSystemStatus":
			var systemStatus domain.SystemStatus
			if err := mapToStruct(dataMap, &systemStatus); err != nil {
				return fmt.Errorf("failed to map SystemStatus: %w", err)
			}
			m.SystemStatus = systemStatus
		case "CLMDateZyouhou":
			var dateInfo domain.DateInfo
			if err := mapToStruct(dataMap, &dateInfo); err != nil {
				return fmt.Errorf("failed to map DateInfo: %w", err)
			}
			m.DateInfo = dateInfo
		case "CLMYobine":
			var callPrice domain.CallPrice
			if err := mapToStruct(dataMap, &callPrice); err != nil {
				return fmt.Errorf("failed to map CallPrice: %w", err)
			}
			m.CallPriceMap[strconv.Itoa(callPrice.UnitNumber)] = callPrice
		case "CLMIssueMstKabu":
			var issueMaster domain.IssueMaster
			if err := mapToStruct(dataMap, &issueMaster); err != nil {
				return fmt.Errorf("failed to map IssueMaster: %w", err)
			}
			m.IssueMap[issueMaster.IssueCode] = issueMaster
		case "CLMIssueSizyouMstKabu":
			var issueMarket domain.IssueMarketMaster
			if err := mapToStruct(dataMap, &issueMarket); err != nil {
				return fmt.Errorf("failed to map IssueMarketMaster: %w", err)
			}
			if _, ok := m.IssueMarketMap[issueMarket.IssueCode]; !ok {
				m.IssueMarketMap[issueMarket.IssueCode] = make(map[string]domain.IssueMarketMaster)
			}
			m.IssueMarketMap[issueMarket.IssueCode][issueMarket.MarketCode] = issueMarket
		case "CLMUnyouStatusKabu":
			var operationStatusKabu domain.OperationStatusKabu
			if err := mapToStruct(dataMap, &operationStatusKabu); err != nil {
				return fmt.Errorf("failed to map OperationStatusKabu: %w", err)
			}
			if _, ok := m.OperationStatusKabuMap[operationStatusKabu.ListedMarket]; !ok {
				m.OperationStatusKabuMap[operationStatusKabu.ListedMarket] = make(map[string]domain.OperationStatusKabu)
			}
			m.OperationStatusKabuMap[operationStatusKabu.ListedMarket][operationStatusKabu.Unit] = operationStatusKabu

		case "CLMIssueSizyouKiseiKabu":
			var issueMarketRegulation domain.IssueMarketRegulation
			if err := mapToStruct(dataMap, &issueMarketRegulation); err != nil {
				return fmt.Errorf("failed to map IssueMarketRegulation: %w", err)
			}
			if _, ok := m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode]; !ok {
				m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode] = make(map[string]domain.IssueMarketRegulation)
			}
			m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode][issueMarketRegulation.ListedMarket] = issueMarketRegulation

		case "CLMEventDownloadComplete":
			return nil

		default:
			tc.Logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}

	return nil
}
