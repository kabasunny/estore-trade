// internal/infrastructure/logger/zapLogger/fact_new_zapLogger_test.go
package zapLogger

import (
	"testing"

	"estore-trade/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewZapLogger(t *testing.T) {
	// テスト用の設定 (ログレベル: info)
	cfg := &config.Config{LogLevel: "info"}

	// NewZapLogger を呼び出す
	logger, err := NewZapLogger(cfg)

	// エラーが発生していないことを確認
	assert.NoError(t, err)
	require.NotNil(t, logger)

	// Observer コアと ObservedLogs を作成
	observedCore, observedLogs := observer.New(zapcore.InfoLevel)

	// テスト用のロガーを作成 (Observer コアを使用)
	testLogger := zap.New(observedCore)

	// テスト用のログメッセージを出力
	testLogger.Debug("This should not be logged")  // Debug レベル (表示されないはず)
	testLogger.Info("This should be logged")       // Info レベル
	testLogger.Error("This should also be logged") // Error レベル

	// ログが期待通りに出力されたか確認
	assert.Equal(t, 2, observedLogs.Len(), "Expected 2 log entries")

	logs := observedLogs.All()
	assert.Equal(t, zapcore.InfoLevel, logs[0].Level, "First log entry should be INFO")
	assert.Equal(t, "This should be logged", logs[0].Message, "First log message mismatch")
	assert.Equal(t, zapcore.ErrorLevel, logs[1].Level, "Second log entry should be ERROR")
	assert.Equal(t, "This should also be logged", logs[1].Message, "Second log message mismatch")
}

func TestNewZapLogger_InvalidLevel(t *testing.T) {
	cfg := &config.Config{LogLevel: "invalid"} // 不正なログレベル
	logger, err := NewZapLogger(cfg)

	assert.Error(t, err)  // エラーが発生することを期待
	assert.Nil(t, logger) // ロガーは nil のはず
}
