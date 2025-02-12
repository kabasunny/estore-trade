package zapLogger

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
		// 開発環境用のデフォルト設定
		zapCfg = zap.NewDevelopmentConfig()
	case "info":
		// 本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
		// ログレベルを情報レベルに設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		// 本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
		// ログレベルを警告レベルに設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		// 本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
		// ログレベルをエラーレベルに設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "dpanic":
		// 本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
		// ログレベルを致命的エラーレベル（DPanic）に設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "panic":
		// 本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
		// ログレベルをパニックレベルに設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal":
		// 本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
		// ログレベルを致命的エラーレベルに設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		// デフォルトでは本番環境用のデフォルト設定
		zapCfg = zap.NewProductionConfig()
	}

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601フォーマット
	zapCfg.OutputPaths = []string{"stdout"}                      // 標準出力
	zapCfg.ErrorOutputPaths = []string{"stderr"}                 // 標準エラー出力

	return zapCfg.Build() //設定された全ての構成オプションを適用して、新しいロガーインスタンスを生成
}
