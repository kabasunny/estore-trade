// internal/infrastructure/logger/zapLogger/fact_new_zapLogger.go
package zapLogger

import (
	"go.uber.org/zap/zapcore"
)

// levelFromString は文字列から zapcore.Level を取得するヘルパー関数
func levelFromString(levelStr string) (zapcore.Level, error) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		return level, err
	}
	return level, nil
}
