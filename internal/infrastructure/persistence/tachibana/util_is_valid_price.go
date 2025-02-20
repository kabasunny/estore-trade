package tachibana

import (
	"math"
)

// isValidPrice は、注文価格が呼値の単位に従っているかをチェックする関数
func isValidPrice(price float64, callPrice CallPrice) bool {
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

	for i := 0; i < len(prices); i++ {
		if price <= prices[i] {
			remainder := math.Mod(price, unitPrices[i])
			return remainder == 0
		}
	}
	return false // ここには到達しないはずだが、念のため
}
