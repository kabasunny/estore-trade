// internal/infrastructure/persistence/tachibana/util_process_response.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

func processResponse(response map[string]interface{}, m *domain.MasterData, tc *TachibanaClientImple) error {
	for sCLMID, data := range response {
		// "sResultCode" は処理しない
		if sCLMID == "sResultCode" {
			continue
		}

		switch sCLMID {
		case "CLMSystemStatus", "CLMDateZyouhou":
			// これらは単一のオブジェクトなので、今まで通り処理
			dataMap, ok := data.(map[string]interface{})
			if !ok {
				tc.logger.Error("Invalid data format in master data response", zap.String("sCLMID", sCLMID))
				continue
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
			}

			//配列として定義
		case "CLMYobine", "CLMIssueMstKabu", "CLMIssueSizyouMstKabu", "CLMIssueSizyouKiseiKabu", "CLMUnyouStatusKabu":
			dataArray, ok := data.([]interface{}) // 配列として扱う
			if !ok {
				tc.logger.Error("Invalid data format in master data response", zap.String("sCLMID", sCLMID))
				continue
			}
			for _, item := range dataArray { // 配列の各要素を処理
				itemMap, ok := item.(map[string]interface{}) // map[string]interface{} に変換
				if !ok {
					tc.logger.Error("Invalid data format in master data response", zap.String("sCLMID", sCLMID))
					continue
				}
				// ここで、sCLMID に応じて適切な構造体にマッピング
				switch sCLMID {
				case "CLMYobine":
					var callPrice domain.CallPrice
					// fmt.Printf("itemMap: %v\n", itemMap)

					if err := mapToStruct(itemMap, &callPrice); err != nil {
						return fmt.Errorf("failed to map CallPrice: %w", err)
					}
					m.CallPriceMap[callPrice.UnitNumber] = callPrice
					// fmt.Printf("callPrice: %v\n", callPrice)
				case "CLMIssueMstKabu":
					var issueMaster domain.IssueMaster
					if err := mapToStruct(itemMap, &issueMaster); err != nil {
						return fmt.Errorf("failed to map IssueMaster: %w", err)
					}
					// fmt.Printf("DEBUG: Mapping issueCode: %s, issueName: %s\n", issueMaster.IssueCode, issueMaster.IssueName)
					m.IssueMap[issueMaster.IssueCode] = issueMaster
					// fmt.Printf("DEBUG: Current IssueMap: %+v\n", m.IssueMap)

				case "CLMIssueSizyouMstKabu":
					var issueMarket domain.IssueMarketMaster
					if err := mapToStruct(itemMap, &issueMarket); err != nil {
						return fmt.Errorf("failed to map IssueMarketMaster: %w", err)
					}
					if _, ok := m.IssueMarketMap[issueMarket.IssueCode]; !ok {
						m.IssueMarketMap[issueMarket.IssueCode] = make(map[string]domain.IssueMarketMaster)
					}
					m.IssueMarketMap[issueMarket.IssueCode][issueMarket.MarketCode] = issueMarket

				case "CLMUnyouStatusKabu":
					var operationStatusKabu domain.OperationStatusKabu
					if err := mapToStruct(itemMap, &operationStatusKabu); err != nil {
						return fmt.Errorf("failed to map OperationStatusKabu: %w", err)
					}
					if _, ok := m.OperationStatusKabuMap[operationStatusKabu.ListedMarket]; !ok {
						m.OperationStatusKabuMap[operationStatusKabu.ListedMarket] = make(map[string]domain.OperationStatusKabu)
					}
					m.OperationStatusKabuMap[operationStatusKabu.ListedMarket][operationStatusKabu.Unit] = operationStatusKabu

				case "CLMIssueSizyouKiseiKabu":
					var issueMarketRegulation domain.IssueMarketRegulation
					if err := mapToStruct(itemMap, &issueMarketRegulation); err != nil {
						return fmt.Errorf("failed to map IssueMarketRegulation: %w", err)
					}
					if _, ok := m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode]; !ok {
						m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode] = make(map[string]domain.IssueMarketRegulation)
					}
					m.IssueMarketRegulationMap[issueMarketRegulation.IssueCode][issueMarketRegulation.ListedMarket] = issueMarketRegulation
				}
			}
		case "CLMEventDownloadComplete": //CLMEventDownloadComplete
		//continue
		// 最後の sCLMID だった場合は、ここで return nil
		//return nil
		default: //sResultCodeをはじく
			tc.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}

	return nil
}
