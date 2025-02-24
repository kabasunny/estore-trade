package auto_usecase

import "estore-trade/internal/domain"

// ShouldTrade はシグナルに基づいて取引を行うべきかどうかを判断する
func ShouldTrade(s *domain.Signal) bool {
	// シグナルに基づいて取引を行うか判断するロジック
	return true // 仮
}
