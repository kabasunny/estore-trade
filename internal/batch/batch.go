// internal/batch/batch.go

package batch

import (
	"context"
	"estore-trade/internal/app"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/batch/signal"
	"fmt"
)

func RunBatch(app *app.App) error {
	ctx := context.Background()

	// 1. ランキング計算
	rank, err := ranking.CalculateRanking(ctx, app.TachibanaClient)
	if err != nil {
		return fmt.Errorf("failed to calculate ranking: %w", err)
	}

	// 2. 銘柄リスト作成
	targetIssues := ranking.CreateTargetIssueList(rank, 10) // 例: 上位10銘柄

	// 3. シグナル生成 (仮: 外部 Python プロセスを呼び出す)
	signals, err := signal.Generate(ctx, app.Logger, targetIssues) //loggerを渡す
	if err != nil {
		return fmt.Errorf("failed to generate signals: %w", err)
	}

	// 4. シグナル保存
	if err := signal.NewSignalRepository(app.DB.DB()).SaveSignals(ctx, signals); err != nil { //SignalRepositoryを使う
		return fmt.Errorf("failed to save signals: %w", err)
	}

	// 5. ポジション計算 (auto_algorithm パッケージのメソッドを呼び出す)
	//    ポジション計算は、autoTradingUsecase が EventStream からのイベントを
	//    ハンドルする際に呼び出されるので、ここでは不要。
	//    ただし、バッチ処理で事前にポジションを計算しておきたい場合は、ここで呼び出す。

	// (オプション) 6.ランキングの履歴を保存する場合
	//if err := ranking.NewRepository(app.DB.DB()).SaveRanking(ctx, rank);err != nil{ //RankingRepositoryを使う
	//	return fmt.Errorf("failed to save ranking: %w", err)
	//}

	// (オプション) 7.取引開始前に注文を送信する場合
	// 翌日の寄り付きで注文を出す場合、ここで PlaceOrder を呼び出す必要はない
	// EventStream を利用している場合は、約定通知 (EC) を受信したら、autoTradingUsecase が
	// HandleEvent() -> autoTradingAlgorithm.CalculatePosition() -> tradingUsecase.PlaceOrder()
	// の順に呼び出して注文を出す。

	// ここでは、注文情報をDBに保存するだけ。
	// 実際の注文は、システムの起動時 (取引開始前) に、別の処理 (例: startTrading) で行う。

	return nil
}
