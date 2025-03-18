// internal/infrastructure/persistence/tachibana/mthd_check_price_is_valid.go
package tachibana

import (
	"context"
	"fmt"
	"strconv"
)

// CheckPriceIsValid は指定された価格が正当であるかを確認
func (tc *TachibanaClientImple) CheckPriceIsValid(ctx context.Context, issueCode string, price float64, isNextDay bool) (bool, error) {
	// priceが0以下、かつ、指値、逆指値の場合のみエラーとする　　　　　　　　　　　　　　　　　　　　　　　　　←ここがポイント
	if price <= 0 {
		//return false, nil //これだとpriceが0以下の時に全てエラーになってしまう
		return false, fmt.Errorf("invalid price: %f", price) //priceの内容をエラーに含める
	}
	tc.mu.RLock()         // ミューテックスのロック
	defer tc.mu.RUnlock() // メソッド終了時にロック解除

	// 銘柄の呼値単位番号を取得 (翌営業日の場合は Yoku を使う)
	issueMarket, ok := tc.GetIssueMarketMaster(issueCode, "00") // Getter を使う
	if !ok {
		return false, fmt.Errorf("IssueMarketMaster not found for issueCode: %s", issueCode)
	}

	unitNumberStr := issueMarket.CallPriceUnitNumber
	if isNextDay {
		unitNumberStr = issueMarket.CallPriceUnitNumberYoku
	}

	// fmt.Printf("issueMarket.issueCode: %s \n", issueMarket.IssueCode)
	// fmt.Printf("issueMarket.CallPriceUnitNumber: %s \n", issueMarket.CallPriceUnitNumber)
	// fmt.Printf("issueMarket.CallPriceUnitNumberYoku: %s \n", issueMarket.CallPriceUnitNumberYoku)
	if unitNumberStr == "" {
		// 呼値情報がない場合はチェック不要 (またはエラーとする)
		return true, nil // または return false, fmt.Errorf(...)
	}

	unitNumber, err := strconv.Atoi(unitNumberStr)
	if err != nil {
		return false, fmt.Errorf("invalid CallPriceUnitNumber: %s", unitNumberStr)
	}

	callPrice, ok := tc.GetCallPrice(strconv.Itoa(unitNumber)) // Getter を使う
	if !ok {
		return false, fmt.Errorf("CallPrice not found for unitNumber: %d", unitNumber)
	}

	// isValidPrice 関数を使ってチェック
	return isValidPrice(price, callPrice), nil
}
