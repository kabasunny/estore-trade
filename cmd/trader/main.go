// cmd/trader/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"estore-trade/internal/app" // internal/app パッケージをインポート
	"estore-trade/internal/batch"

	"go.uber.org/zap"
)

func main() {
	// app パッケージの初期化関数を呼び出す
	app, err := app.InitApp() // 各インスタンスや設定の初期化
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close() // クリーンアップ処理

	errCh := make(chan error, 1) // エラーチャネル

	ctx := context.Background()

	// EventStream の開始 (ゴルーチンで実行)
	go func() { // EventStreamの開始
		if err := app.EventStream.Start(ctx); err != nil { //app経由で呼び出し
			app.Logger.Error("EventStream error", zap.Error(err))
			errCh <- err
		}
	}()

	// AutoTradingUsecase を開始 (ゴルーチンで実行)
	go func() { // AutoTradingUsecaseを開始
		if err := app.AutoTradingUsecase.Start(); err != nil { //app経由で呼び出し
			app.Logger.Error("AutoTradingUsecase error", zap.Error(err))
			errCh <- err
		}
	}()

	// HTTP サーバーの起動
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.Logger.Info("Starting server on port", zap.Int("port", app.Config.HTTPPort)) //app経由で呼び出し
	go func() {
		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed { //app経由で呼び出し
			app.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// バッチ処理の実行 (スケジューラで実行)
	go func() {
		if err := batch.RunBatch(app); err != nil { //バッチ処理を呼び出し
			app.Logger.Error("Batch processing error", zap.Error(err))
			errCh <- err
		}
	}()

	// シグナル処理とサーバーのシャットダウン
	select {
	case <-ctx.Done():
	case err := <-errCh:
		app.Logger.Error("Received error from goroutine", zap.Error(err))
	}

	app.Logger.Info("Shutting down server...")

	if err := app.EventStream.Stop(); err != nil { //app経由で呼び出し
		app.Logger.Error("Failed to stop EventStream", zap.Error(err))
	}
	// AutoTradingUsecase停止
	if err := app.AutoTradingUsecase.Stop(); err != nil {
		app.Logger.Error("Failed to stop AutoTradingUsecase", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.HTTPServer.Shutdown(ctx); err != nil { //app経由で呼び出し
		app.Logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	app.Logger.Info("Server exiting")
}
