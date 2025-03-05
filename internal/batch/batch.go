// internal/batch/batch.go
package batch

import (
	"context"
	"estore-trade/internal/app"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/signal"
	"fmt"
)

func RunBatch(app *app.App) error {
	ctx := context.Background()

	// 1. マスターデータのダウンロードと保存 (最初に実行)
	md, err := app.TachibanaClient.DownloadMasterData(ctx) //変更
	if err != nil {
		return fmt.Errorf("failed to download master data: %w", err)
	}
	if err := app.MasterDataRepo.SaveMasterData(ctx, md); err != nil { //変更
		return fmt.Errorf("failed to save master data: %w", err)
	}

	// 2. 銘柄コードリスト取得 (DB から)
	allIssueCodes, err := app.MasterDataRepo.GetAllIssueCodes(ctx) //変更
	if err != nil {
		return fmt.Errorf("failed to get all issue codes: %w", err)
	}

	// GetMarketData を呼び出して、必要なデータを取得
	marketData, err := ranking.GetMarketData(ctx, app.TachibanaClient, allIssueCodes)
	if err != nil {
		return fmt.Errorf("failed to get market data: %w", err)
	}

	// 3. 売買代金ランキング計算
	rank, err := ranking.CalculateRanking(ctx, marketData) //tachibanaClientを削除
	if err != nil {
		return fmt.Errorf("failed to calculate ranking: %w", err)
	}

	// (以下、既存のコード)
	// 4. 銘柄リスト作成
	targetIssues := ranking.CreateTargetIssueList(rank, 10)

	// 5. シグナル生成
	signals, err := Generate(ctx, app.Logger, targetIssues)
	if err != nil {
		return fmt.Errorf("failed to generate signals: %w", err)
	}

	// 6. シグナル保存
	signalRepo := signal.NewSignalRepository(app.DB.DB())
	if err := signalRepo.SaveSignals(ctx, signals); err != nil {
		return fmt.Errorf("failed to save signals: %w", err)
	}
	// 7. ポジション計算
	positions := make([]domain.Position, 0)
	for _, sig := range signals {
		position, err := app.AutoTradingUsecase.AutoTradingAlgorithm().CalculatePosition(&sig)
		if err != nil {
			fmt.Printf("failed to calculate position for signal %+v: %v\n", sig, err)
			continue
		}
		positions = append(positions, *position)
	}
	return nil
}
