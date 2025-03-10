// internal/infrastructure/persistence/tachibana/util_map_to_struct.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
	"strconv"
)

// mapToStruct は、map[string]interface{} を構造体にマッピングします。
// 今回は、IssueMaster, SystemStatus, DateInfo, OperationStatusKabu 構造体に特化した実装になっています。
func mapToStruct(data map[string]interface{}, result interface{}) error {
	switch v := result.(type) {
	case *domain.IssueMaster: // IssueMasterへのポインタの場合
		if value, ok := data["sIssueCode"].(string); ok {
			v.IssueCode = value
		}
		if value, ok := data["sIssueName"].(string); ok {
			v.IssueName = value
		}
		// 他の文字列フィールドも同様に...
		if value, ok := data["sTokuteiF"].(string); ok {
			v.TokuteiF = value
		}
		// 数値型フィールド (文字列として返ってくる場合)  ここは削除

	case *domain.SystemStatus:
		if value, ok := data["sSystemStatusKey"].(string); ok {
			v.SystemStatusKey = value
		}
		if value, ok := data["sLoginKyokaKubun"].(string); ok {
			v.LoginPermission = value
		}
		if value, ok := data["sSystemStatus"].(string); ok {
			v.SystemState = value
		}

	case *domain.DateInfo:
		if value, ok := data["sDayKey"].(string); ok {
			v.DateKey = value
		}
		// 他のフィールドも同様に...
		if value, ok := data["sMaeEigyouDay_1"].(string); ok {
			v.PrevBusinessDay1 = value
		}
		if value, ok := data["sTheDay"].(string); ok {
			v.TheDay = value
		}
		if value, ok := data["sYokuEigyouDay_1"].(string); ok {
			v.NextBusinessDay1 = value
		}
		if value, ok := data["sKabuUkewatasiDay"].(string); ok {
			v.StockDeliveryDate = value
		}

	case *domain.CallPrice: //callprice
		if value, ok := data["sTaniNo"].(string); ok {
			v.UnitNumber = value
		}
		if value, ok := data["sTekiyouDay"].(string); ok {
			v.ApplyDate = value
		}

		// 基準値段と呼値単価 (float64, string からの変換が必要)
		for i := 1; i <= 20; i++ {
			priceKey := fmt.Sprintf("sKizunPrice_%d", i)
			unitPriceKey := fmt.Sprintf("sYobineTanka_%d", i)
			decimalKey := fmt.Sprintf("sDecimal_%d", i)

			if priceStr, ok := data[priceKey].(string); ok {
				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %s: %w", priceKey, err)
				}
				switch i { // 構造体のフィールドに代入
				case 1:
					v.Price1 = price
				case 2:
					v.Price2 = price
				case 3:
					v.Price3 = price
				case 4:
					v.Price4 = price
				case 5:
					v.Price5 = price
				case 6:
					v.Price6 = price
				case 7:
					v.Price7 = price
				case 8:
					v.Price8 = price
				case 9:
					v.Price9 = price
				case 10:
					v.Price10 = price
				case 11:
					v.Price11 = price
				case 12:
					v.Price12 = price
				case 13:
					v.Price13 = price
				case 14:
					v.Price14 = price
				case 15:
					v.Price15 = price
				case 16:
					v.Price16 = price
				case 17:
					v.Price17 = price
				case 18:
					v.Price18 = price
				case 19:
					v.Price19 = price
				case 20:
					v.Price20 = price
				}
			}

			if unitPriceStr, ok := data[unitPriceKey].(string); ok {
				unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
				if err != nil {
					return fmt.Errorf("failed to parse %s: %w", unitPriceKey, err)
				}
				switch i { //構造体のフィールドに代入
				case 1:
					v.UnitPrice1 = unitPrice
				case 2:
					v.UnitPrice2 = unitPrice
				case 3:
					v.UnitPrice3 = unitPrice
				case 4:
					v.UnitPrice4 = unitPrice
				case 5:
					v.UnitPrice5 = unitPrice
				case 6:
					v.UnitPrice6 = unitPrice
				case 7:
					v.UnitPrice7 = unitPrice
				case 8:
					v.UnitPrice8 = unitPrice
				case 9:
					v.UnitPrice9 = unitPrice
				case 10:
					v.UnitPrice10 = unitPrice
				case 11:
					v.UnitPrice11 = unitPrice
				case 12:
					v.UnitPrice12 = unitPrice
				case 13:
					v.UnitPrice13 = unitPrice
				case 14:
					v.UnitPrice14 = unitPrice
				case 15:
					v.UnitPrice15 = unitPrice
				case 16:
					v.UnitPrice16 = unitPrice
				case 17:
					v.UnitPrice17 = unitPrice
				case 18:
					v.UnitPrice18 = unitPrice
				case 19:
					v.UnitPrice19 = unitPrice
				case 20:
					v.UnitPrice20 = unitPrice
				}
			}
			if decimalStr, ok := data[decimalKey].(string); ok {
				decimal, err := strconv.Atoi(decimalStr)
				if err != nil {
					return fmt.Errorf("failed to parse %s to int: %w", decimalKey, err)
				}
				switch i { //構造体のフィールドに代入
				case 1:
					v.Decimal1 = decimal
				case 2:
					v.Decimal2 = decimal
				case 3:
					v.Decimal3 = decimal
				case 4:
					v.Decimal4 = decimal
				case 5:
					v.Decimal5 = decimal
				case 6:
					v.Decimal6 = decimal
				case 7:
					v.Decimal7 = decimal
				case 8:
					v.Decimal8 = decimal
				case 9:
					v.Decimal9 = decimal
				case 10:
					v.Decimal10 = decimal
				case 11:
					v.Decimal11 = decimal
				case 12:
					v.Decimal12 = decimal
				case 13:
					v.Decimal13 = decimal
				case 14:
					v.Decimal14 = decimal
				case 15:
					v.Decimal15 = decimal
				case 16:
					v.Decimal16 = decimal
				case 17:
					v.Decimal17 = decimal
				case 18:
					v.Decimal18 = decimal
				case 19:
					v.Decimal19 = decimal
				case 20:
					v.Decimal20 = decimal
				}
			}
		}

	case *domain.IssueMarketMaster:
		if value, ok := data["sIssueCode"].(string); ok {
			v.IssueCode = value
		}
		if value, ok := data["sSizyouC"].(string); ok {
			v.MarketCode = value
		}

		// 他のフィールドも同様に...
		if value, ok := data["sNehabaMin"].(string); ok {
			if value == "" { //空文字チェック
				v.PriceRangeMin = 0.0
			} else {
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("failed to convert sNehabaMin to float64: %w", err)
				}
				v.PriceRangeMin = f
			}
		}
		if value, ok := data["sNehabaMax"].(string); ok {
			if value == "" { //空文字チェック
				v.PriceRangeMax = 0.0
			} else {
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("failed to convert sNehabaMax to float64: %w", err)
				}
				v.PriceRangeMax = f
			}
		}
		if value, ok := data["sSinyouC"].(string); ok {
			v.SinyouC = value
		}
		if value, ok := data["sZenzituOwarine"].(string); ok {
			if value == "" { //空文字チェック
				v.PreviousClose = 0.0
			} else {
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("failed to convert sZenzituOwarine to float64: %w", err)
				}
				v.PreviousClose = f
			}
		}
		if value, ok := data["sIssueKubunC"].(string); ok {
			v.IssueKubunC = value
		}
		if value, ok := data["sZyouzyouKubun"].(string); ok {
			v.ZyouzyouKubun = value
		}
		if value, ok := data["sYobineTaniNo"].(string); ok {
			v.CallPriceUnitNumber = value
		}
		if value, ok := data["sYobineTaniNoYoku"].(string); ok {
			v.CallPriceUnitNumberYoku = value
		}

	case *domain.IssueMarketRegulation:
		if value, ok := data["sIssueCode"].(string); ok {
			v.IssueCode = value
		}
		// 他のフィールドも同様に...

	case *domain.OperationStatusKabu:
		if value, ok := data["sZyouzyouSizyou"].(string); ok {
			v.ListedMarket = value
		}
		if value, ok := data["sUnyouUnit"].(string); ok {
			v.Unit = value
		}
		if value, ok := data["sUnyouStatus"].(string); ok {
			v.Status = value
		}

	default:
		return fmt.Errorf("unsupported type: %T", result)
	}

	return nil
}
