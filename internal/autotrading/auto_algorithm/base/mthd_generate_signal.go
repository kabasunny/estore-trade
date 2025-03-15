package auto_algorithm

import (
	"estore-trade/internal/domain"
)

func (a *AutoTradingAlgorithm) GenerateSignal(event domain.OrderEvent) (*domain.Signal, error) {
	// signalを生成
	return &domain.Signal{}, nil // 仮
}
