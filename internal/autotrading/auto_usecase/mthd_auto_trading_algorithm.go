package auto_usecase // 変更

import (
	// algorithmパッケージをインポート
	"estore-trade/internal/autotrading/auto_algorithm"
)

func (a *autoTradingUsecase) AutoTradingAlgorithm() *auto_algorithm.AutoTradingAlgorithm {
	return a.autoTradingAlgorithm
}
