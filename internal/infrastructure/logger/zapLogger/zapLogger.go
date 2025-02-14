package zapLogger

import (
	"estore-trade/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 各コンポーネント（PostgresDB, TachibanaClientIntImple, TradingHandler, tradingUsecase）が具体的なロガー実装（zap）に直接依存しないように
// ロギング機能を抽象化し、後で別のロガー（例えば logrus）に切り替えたり、テスト時にモックのロガーを注入したりすることを容易にする
// main.go はアプリケーションのエントリーポイントであり、具体的なロガー実装（zap）を直接使用しても、他のコンポーネントへの影響は少ないため、簡潔さを優先

// NewZapLogger は設定に基づいて新しいzap.Loggerインスタンスを初期化
func NewZapLogger(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg zap.Config

	// ログレベルに基づいて設定を調整
	switch cfg.LogLevel {
	case "debug":
		zapCfg = zap.NewDevelopmentConfig() // 開発環境用のデフォルト設定
	case "info":
		zapCfg = zap.NewProductionConfig()                     // 本番環境用のデフォルト設定
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // ログレベルを情報レベルに設定
	case "warn":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel) // ログレベルを警告レベルに設定
	case "error":
		zapCfg = zap.NewProductionConfig()

		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel) // ログレベルをエラーレベルに設定
	case "dpanic":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DPanicLevel) // ログレベルを致命的エラーレベル（DPanic）に設定
	case "panic":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel) // ログレベルをパニックレベルに設定
	case "fatal":
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel) // ログレベルを致命的エラーレベルに設定
	default:
		zapCfg = zap.NewProductionConfig() // デフォルトでは本番環境用のデフォルト設定
	}

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601フォーマット
	zapCfg.OutputPaths = []string{"stdout"}                      // 標準出力
	zapCfg.ErrorOutputPaths = []string{"stderr"}                 // 標準エラー出力

	return zapCfg.Build() //設定された全ての構成オプションを適用して、新しいロガーインスタンスを生成
}
