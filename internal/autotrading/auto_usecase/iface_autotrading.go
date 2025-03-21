package auto_usecase

import (
	"estore-trade/internal/autotrading/auto_algorithm"
	"estore-trade/internal/domain" // ドメインモデルをインポート
)

type AutoTradingUsecase interface {
	Start() error                        // 自動売買を開始 (EventStream からのデータ受信、シグナル生成、注文実行など)
	Stop() error                         // 自動売買を停止
	HandleEvent(event domain.OrderEvent) // EventStream からのイベントを処理
	AutoTradingAlgorithm() *auto_algorithm.AutoTradingAlgorithm
}
