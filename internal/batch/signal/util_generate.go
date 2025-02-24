// internal/batch/signal/generate.go
package signal

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// 銘柄リストに基づいて売買シグナルを生成する
func Generate(ctx context.Context, logger *zap.Logger, targetIssues []domain.TargetIssue) ([]domain.Signal, error) {
	// TODO: 実際には、外部の Python プロセスを呼び出すか、
	//       auto_algorithm.AutoTradingAlgorithm.GenerateSignal() を呼び出す

	// ここでは、仮のシグナルデータを生成
	logger.Info("Generate Signal!")
	var signals []domain.Signal
	for i, issue := range targetIssues {
		signals = append(signals, domain.Signal{
			ID:        i + 1,
			IssueCode: issue.IssueCode,
			Side:      "buy", // 仮
			Priority:  1,     // 仮
			CreatedAt: time.Now(),
		})
		fmt.Println(signals)
	}

	return signals, nil
}
