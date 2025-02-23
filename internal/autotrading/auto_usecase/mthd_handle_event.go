package auto_usecase // 変更

import (
	"context"
	"estore-trade/internal/domain"

	// algorithmパッケージをインポート
	"go.uber.org/zap"
)

func (a *autoTradingUsecase) HandleEvent(event domain.OrderEvent) {
	// 1. EventStreamからの約定通知などのイベントを受け取る

	// 2. イベントに基づいて、自動売買アルゴリズムを呼び出す
	signal, err := a.autoTradingAlgorithm.GenerateSignal(event) // シグナル生成
	if err != nil {
		a.logger.Error("Signal generate error", zap.Error(err))
	}

	if signal.ShouldTrade() { // シグナルに基づいて取引を行うか判断
		// 3. 資金リスク管理を行った上のポジションの決定
		position, err := a.autoTradingAlgorithm.CalculatePosition(signal)
		if err != nil {
			a.logger.Error("Position calculate error", zap.Error(err))
		}

		// 4. 既存の tradingUsecase を使って注文を送信
		order := domain.Order{ // domain.Order を作成
			// ... position の情報に基づいて必要な値を設定 ...
			// 例:
			Symbol:    position.Symbol,   // 銘柄コード (仮)
			Side:      position.Side,     // 売買区分 ("buy" or "sell") (仮)
			OrderType: "market",          // 指値・成行など (ここでは成行を仮定)
			Quantity:  position.Quantity, // 数量 (仮)
			// Price は成行注文の場合は設定しない (または 0 などの特別な値を設定)
		}
		// userID, password を削除
		if _, err := a.tradingUsecase.PlaceOrder(context.Background(), &order); err != nil {
			a.logger.Error("auto trading order error", zap.Error(err))
		}
	}
}
