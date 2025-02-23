package auto_algorithm

import (
	"estore-trade/internal/autotrading/auto_model"
	"estore-trade/internal/domain"
)

func (a *AutoTradingAlgorithm) GenerateSignal(event domain.OrderEvent) (*auto_model.Signal, error) {
	// signalを生成
	return &auto_model.Signal{}, nil // 仮
}
