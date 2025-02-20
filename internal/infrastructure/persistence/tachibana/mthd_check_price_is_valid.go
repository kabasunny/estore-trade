package tachibana

import (
	"fmt"
	"strconv"
)

// CheckPriceIsValid は指定された価格が正当であるかを確認
func (tc *TachibanaClientImple) CheckPriceIsValid(issueCode string, price float64, isNextDay bool) (bool, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// 銘柄の呼値単位番号を取得 (翌営業日の場合は Yoku を使う)
	issueMarket, ok := tc.issueMarketMap[issueCode]["00"] // 例として市場コード"00" (東証) を使用
	if !ok {
		return false, fmt.Errorf("IssueMarketMaster not found for issueCode: %s", issueCode)
	}
	unitNumberStr := issueMarket.CallPriceUnitNumber
	if isNextDay {
		unitNumberStr = issueMarket.CallPriceUnitNumberYoku
	}
	if unitNumberStr == "" {
		// 呼値情報がない場合はチェック不要 (またはエラーとする)
		return true, nil // または return false, fmt.Errorf(...)
	}

	unitNumber, err := strconv.Atoi(unitNumberStr)
	if err != nil {
		return false, fmt.Errorf("invalid CallPriceUnitNumber: %s", unitNumberStr)
	}

	callPrice, ok := tc.callPriceMap[strconv.Itoa(unitNumber)]
	if !ok {
		return false, fmt.Errorf("CallPrice not found for unitNumber: %d", unitNumber)
	}

	// isValidPrice 関数を使ってチェック
	return isValidPrice(price, callPrice), nil
}
