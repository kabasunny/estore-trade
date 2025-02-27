// internal/infrastructure/logger/zapLogger/fact_new_zapLogger.go
package zapLogger

import (
	"estore-trade/internal/config"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger は設定に基づいて新しいzap.Loggerインスタンスを初期化
func NewZapLogger(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg zap.Config

	// levelFromString を switch の外で呼び出す
	level, err := levelFromString(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// ログレベルに基づいて設定を調整
	switch cfg.LogLevel {
	case "debug":
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
	case "info", "warn", "error", "dpanic", "panic", "fatal": //本番環境で適切なもの
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level) // ログレベルを設定
	default: // もしLogLevelが空か、上記以外だったら
		// default case は不要になる
		// 有効なログレベルが設定されている場合は、ここには到達しない
		// 無効なログレベルの場合は、levelFromString でエラーが返される
	}

	// ... (その他の設定は変更なし) ...
	//エンコーディングをJSONに
	zapCfg.Encoding = "json"
	// タイムスタンプのフォーマットを ISO8601 に設定
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//標準出力と標準エラー出力
	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}
	return logger, nil
}
