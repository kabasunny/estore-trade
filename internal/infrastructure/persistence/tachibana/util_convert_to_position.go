// internal/infrastructure/persistence/tachibana/util_convert_to_position.go
package tachibana

import (
	"fmt"
	"strconv"
	"time"

	"estore-trade/internal/domain"
)

// convertToPosition は map[string]interface{} を domain.Position に変換するヘルパー関数
func convertToPosition(data map[string]interface{}) (*domain.Position, error) {
	position := &domain.Position{}

	if id, ok := data["sOrderTategyokuNumber"].(string); ok {
		position.ID = id // 建玉番号
	}
	if symbol, ok := data["sOrderIssueCode"].(string); ok {
		position.Symbol = symbol
	}

	if side, ok := data["sOrderBaibaiKubun"].(string); ok {
		switch side {
		case "1": // 売り
			position.Side = "short" // "sell" から "short" へ変更
		case "3": // 買い
			position.Side = "long" // "buy" から "long" へ変更
		default:
			return nil, fmt.Errorf("invalid sOrderBaibaiKubun: %s", side)
		}
	}

	if quantityStr, ok := data["sOrderTategyokuSuryou"].(string); ok {
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sOrderTategyokuSuryou: %w", err)
		}
		position.Quantity = quantity
	}

	if priceStr, ok := data["sOrderTategyokuTanka"].(string); ok {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sOrderTategyokuTanka: %w", err)
		}
		position.Price = price
	}

	if dateStr, ok := data["sOrderTategyokuDay"].(string); ok {
		// "YYYYMMDD" 形式の日付を time.Time に変換
		t, err := time.Parse("20060102", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sOrderTategyokuDay: %w", err)
		}
		position.OpenDate = t // OpenDate に設定
	}

	if dueDateStr, ok := data["sOrderTategyokuKizituDay"].(string); ok {
		position.DueDate = dueDateStr // sOrderTategyokuKizituDay の値をそのまま DueDate に設定
	}

	if tradeTypeStr, ok := data["sOrderBensaiKubun"].(string); ok {
		position.MarginTradeType = tradeTypeStr // sOrderBensaiKubun の値をそのまま MarginTradeType に設定
	}

	return position, nil
}
