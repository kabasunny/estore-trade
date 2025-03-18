// internal/infrastructure/persistence/tachibana/util_is_valid_price.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
	"math"
)

// isValidPrice は、注文価格が呼値の単位に従っているかをチェックする関数
func isValidPrice(price float64, callPrice domain.CallPrice) bool {
	// fmt.Printf("isValidPrice called with price: %f, callPrice: %+v\n", price, callPrice)

	prices := [20]float64{
		callPrice.Price1, callPrice.Price2, callPrice.Price3, callPrice.Price4, callPrice.Price5,
		callPrice.Price6, callPrice.Price7, callPrice.Price8, callPrice.Price9, callPrice.Price10,
		callPrice.Price11, callPrice.Price12, callPrice.Price13, callPrice.Price14, callPrice.Price15,
		callPrice.Price16, callPrice.Price17, callPrice.Price18, callPrice.Price19, callPrice.Price20,
	}
	unitPrices := [20]float64{
		callPrice.UnitPrice1, callPrice.UnitPrice2, callPrice.UnitPrice3, callPrice.UnitPrice4, callPrice.UnitPrice5,
		callPrice.UnitPrice6, callPrice.UnitPrice7, callPrice.UnitPrice8, callPrice.UnitPrice9, callPrice.UnitPrice10,
		callPrice.UnitPrice11, callPrice.UnitPrice12, callPrice.UnitPrice13, callPrice.UnitPrice14, callPrice.UnitPrice15,
		callPrice.UnitPrice16, callPrice.UnitPrice17, callPrice.UnitPrice18, callPrice.UnitPrice19, callPrice.UnitPrice20,
	}

	// priceがどの価格帯に属するか判定する
	for i := 0; i < len(prices); i++ {
		if i == 0 {
			// 最初の価格帯
			if price <= prices[i] {
				remainder := math.Mod(price, unitPrices[i])
				fmt.Printf("  First price range check: price <= prices[%d] (%f <= %f), unitPrices[%d]: %f, remainder: %f\n", i, price, prices[i], i, unitPrices[i], remainder) // ログ出力: i は %d
				return remainder == 0
			}
		} else {
			// それ以降の価格帯
			if prices[i-1] < price && price <= prices[i] {
				remainder := math.Mod(price-prices[i-1], unitPrices[i])
				fmt.Printf("  Price range check: prices[%d] < price <= prices[%d] (%f < %f <= %f), unitPrices[%d]: %f, remainder: %f\n", i-1, i, prices[i-1], price, prices[i], i, unitPrices[i], remainder) // ログ出力: i-1, i は %d
				return remainder == 0
			}
		}
	}

	// どの価格帯にも属さない場合は呼値テーブルの範囲外
	fmt.Printf("  Price out of range\n") // ログ出力
	return false
}
