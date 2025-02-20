package tachibana

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

func processResponse(response map[string]interface{}, m *masterDataManager, tc *TachibanaClientImple) error {
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
			if err := mapToStruct(dataMap, &systemStatus); err != nil {
				return fmt.Errorf("failed to map SystemStatus: %w", err)
			}
			m.systemStatus = systemStatus
		case "CLMDateZyouhou":
			var dateInfo DateInfo
			if err := mapToStruct(dataMap, &dateInfo); err != nil {
				return fmt.Errorf("failed to map DateInfo: %w", err)
			}
			m.dateInfo = dateInfo
		case "CLMYobine":
			var callPrice CallPrice
			if err := mapToStruct(dataMap, &callPrice); err != nil {
				return fmt.Errorf("failed to map CallPrice: %w", err)
			}
			m.callPriceMap[strconv.Itoa(callPrice.UnitNumber)] = callPrice
		case "CLMIssueMstKabu":
			var issueMaster IssueMaster
			if err := mapToStruct(dataMap, &issueMaster); err != nil {
				return fmt.Errorf("failed to map IssueMaster: %w", err)
			}
			m.issueMap[issueMaster.IssueCode] = issueMaster
		case "CLMIssueSizyouMstKabu":
			var issueMarket IssueMarketMaster
			if err := mapToStruct(dataMap, &issueMarket); err != nil {
				return fmt.Errorf("failed to map IssueMarketMaster: %w", err)
			}
			if _, ok := m.issueMarketMap[issueMarket.IssueCode]; !ok {
				m.issueMarketMap[issueMarket.IssueCode] = make(map[string]IssueMarketMaster)
			}
			m.issueMarketMap[issueMarket.IssueCode][issueMarket.MarketCode] = issueMarket
		case "CLMUnyouStatusKabu":
			var operationStatusKabu OperationStatusKabu
			if err := mapToStruct(dataMap, &operationStatusKabu); err != nil {
				return fmt.Errorf("failed to map OperationStatusKabu: %w", err)
			}
			if _, ok := m.operationStatusKabuMap[operationStatusKabu.ListedMarket]; !ok {
				m.operationStatusKabuMap[operationStatusKabu.ListedMarket] = make(map[string]OperationStatusKabu)
			}
			m.operationStatusKabuMap[operationStatusKabu.ListedMarket][operationStatusKabu.Unit] = operationStatusKabu

		case "CLMIssueSizyouKiseiKabu":
			var issueMarketRegulation IssueMarketRegulation
			if err := mapToStruct(dataMap, &issueMarketRegulation); err != nil {
				return fmt.Errorf("failed to map IssueMarketRegulation: %w", err)
			}
			if _, ok := m.issueMarketRegulationMap[issueMarketRegulation.IssueCode]; !ok {
				m.issueMarketRegulationMap[issueMarketRegulation.IssueCode] = make(map[string]IssueMarketRegulation)
			}
			m.issueMarketRegulationMap[issueMarketRegulation.IssueCode][issueMarketRegulation.ListedMarket] = issueMarketRegulation

		case "CLMEventDownloadComplete":
			return nil

		default:
			tc.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}

	return nil
}
