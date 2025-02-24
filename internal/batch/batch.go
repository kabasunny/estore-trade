// internal/batch/batch.go
package batch

import (
	"context"
	"estore-trade/internal/app" // 追加
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/signal"
	"fmt"
	//"go.uber.org/zap" //mainでログ出力するので不要
)

func RunBatch(app *app.App) error {
	ctx := context.Background()

	// 1. 銘柄コードリスト取得 (ここでは省略 - 事前に取得済みとする)

	// 2. 株価・出来高取得
	allIssueCodes := []string{"7203", "8306", "9432"} // 仮
	marketData, err := ranking.GetMarketData(ctx, app.TachibanaClient, allIssueCodes)
	if err != nil {
		return fmt.Errorf("failed to get market data: %w", err)
	}

	// 3. 売買代金ランキング計算
	rank, err := ranking.CalculateRanking(ctx, marketData) // marketDataを渡す
	if err != nil {
		return fmt.Errorf("failed to calculate ranking: %w", err)
	}

	// 4. 銘柄リスト作成
	targetIssues := ranking.CreateTargetIssueList(rank, 10) // 例: 上位10銘柄

	// 5. シグナル生成 (外部 Python プロセス、または Go の関数を呼び出す)
	signals, err := Generate(ctx, app.Logger, targetIssues)
	if err != nil {
		return fmt.Errorf("failed to generate signals: %w", err)
	}

	// 6. シグナル保存 (signalRepository を使用)
	signalRepo := signal.NewSignalRepository(app.DB.DB())
	if err := signalRepo.SaveSignals(ctx, signals); err != nil {
		return fmt.Errorf("failed to save signals: %w", err)
	}

	// 7. ポジション計算 (auto_algorithm パッケージのメソッドを呼び出す)
	// autoTradingUsecase 経由で autoTradingAlgorithm にアクセス
	// autoTradingUsecase 経由で autoTradingAlgorithm にアクセス
	positions := make([]domain.Position, 0)
	for _, sig := range signals {
		position, err := app.AutoTradingUsecase.AutoTradingAlgorithm().CalculatePosition(&sig) //auto_model.Signalへのポインタを渡す
		if err != nil {
			// ポジション計算に失敗した場合でも、他の銘柄の処理は続行する
			// エラーはログに出力するなどの処理が必要
			fmt.Printf("failed to calculate position for signal %+v: %v\n", sig, err)
			continue //取引しない場合はcontinue
		}
		positions = append(positions, *position) //ここもポインタを外す
	}

	// (8. 注文実行は、取引開始時に別の処理で行う)

	return nil
}
