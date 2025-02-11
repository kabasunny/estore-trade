package zap

import (
	"estore-trade/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger は設定に基づいて新しいzap.Loggerインスタンスを初期化
func NewZapLogger(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg zap.Config

	// ログレベルに基づいて設定を調整
	switch cfg.LogLevel {
	case "debug":
		zapCfg = zap.NewDevelopmentConfig()
	case "info":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "dpanic":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "panic":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		zapCfg = zap.NewProductionConfig() // デフォルトはproduction
	}

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601フォーマット
	zapCfg.OutputPaths = []string{"stdout"}                      // 標準出力
	zapCfg.ErrorOutputPaths = []string{"stderr"}                 // 標準エラー出力

	return zapCfg.Build()
}
