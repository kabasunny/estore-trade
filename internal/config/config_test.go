// internal/config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// 共通の設定 (全テストケースで共通)
	os.Setenv("TACHIBANA_BASE_URL", "https://example.com")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("LOG_LEVEL", "debug")

	// .env ファイルが存在しない場合
	t.Run("no .env file", func(t *testing.T) {
		// DB_PORT, HTTP_PORT は設定によって挙動が変わるので、一旦クリア
		os.Unsetenv("DB_PORT")
		os.Unsetenv("HTTP_PORT")

		cfg, err := LoadConfig("nonexistent.env") // 存在しないファイル
		assert.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "https://example.com", cfg.TachibanaBaseURL)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, 5432, cfg.DBPort) // デフォルト値
		assert.Equal(t, "debug", cfg.LogLevel)
		assert.Equal(t, 8080, cfg.HTTPPort) // デフォルト値
		// 他の必要な環境変数も同様に確認
	})

	// DB_PORT が設定されていない場合のテスト
	t.Run("DB_PORT not set", func(t *testing.T) {
		os.Unsetenv("DB_PORT") // DB_PORT を未設定にする
		cfg, err := LoadConfig("nonexistent.env")
		assert.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, 5432, cfg.DBPort) // デフォルト値が使用される
	})

	// DB_PORT が不正な値の場合のテスト
	t.Run("DB_PORT invalid", func(t *testing.T) {
		os.Setenv("DB_PORT", "invalid")            // 不正な値を設定
		defer os.Unsetenv("DB_PORT")               // テスト終了後に確実にクリア
		_, err := LoadConfig("nonexistent.env")    // 存在しないファイルを指定
		assert.Error(t, err)                       // エラーが発生するはず
		assert.Contains(t, err.Error(), "DB_PORT") // エラーメッセージに DB_PORT が含まれる
	})

	// HTTP_PORT が設定されていない場合のテスト (デフォルト値の確認)
	t.Run("HTTP_PORT not set", func(t *testing.T) {
		os.Unsetenv("HTTP_PORT") // HTTP_PORT を未設定にする
		cfg, err := LoadConfig("nonexistent.env")
		assert.NoError(t, err) // エラーは発生しない
		require.NotNil(t, cfg)
		assert.Equal(t, 8080, cfg.HTTPPort) // デフォルト値が使用されることを確認
	})

	// HTTP_PORT が不正な値の場合のテスト
	t.Run("HTTP_PORT invalid", func(t *testing.T) {
		os.Setenv("HTTP_PORT", "invalid") // 不正な値を設定
		defer os.Unsetenv("HTTP_PORT")
		_, err := LoadConfig("nonexistent.env")      // 存在しないファイルを指定
		assert.Error(t, err)                         // エラーが発生するはず
		assert.Contains(t, err.Error(), "HTTP_PORT") // エラーメッセージに HTTP_PORT が含まれる
	})
	// 正常な場合のテスト
	t.Run("valid config", func(t *testing.T) {
		// テスト用の環境変数を設定
		os.Setenv("TACHIBANA_BASE_URL", "https://example.com")
		os.Setenv("TACHIBANA_USER_ID", "testuser")
		os.Setenv("TACHIBANA_PASSWORD", "testpassword")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("LOG_LEVEL", "info")
		os.Setenv("EVENT_RID", "rid")
		os.Setenv("EVENT_BOARD_NO", "boardno")
		os.Setenv("EVENT_EVT_CMD", "cmd")
		os.Setenv("HTTP_PORT", "8080")

		// テスト終了後に環境変数を元に戻す
		defer func() {
			os.Unsetenv("TACHIBANA_BASE_URL")
			os.Unsetenv("TACHIBANA_USER_ID")
			os.Unsetenv("TACHIBANA_PASSWORD")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("EVENT_RID")
			os.Unsetenv("EVENT_BOARD_NO")
			os.Unsetenv("EVENT_EVT_CMD")
			os.Unsetenv("HTTP_PORT")
		}()

		cfg, err := LoadConfig("nonexistent.env") // 存在しないファイル
		assert.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "https://example.com", cfg.TachibanaBaseURL)
		assert.Equal(t, "testuser", cfg.TachibanaUserID)
		assert.Equal(t, "testpassword", cfg.TachibanaPassword)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, 5432, cfg.DBPort)
		assert.Equal(t, "testuser", cfg.DBUser)
		assert.Equal(t, "testpassword", cfg.DBPassword)
		assert.Equal(t, "testdb", cfg.DBName)
		assert.Equal(t, "info", cfg.LogLevel)
		assert.Equal(t, "rid", cfg.EventRid)
		assert.Equal(t, "boardno", cfg.EventBoardNo)
		assert.Equal(t, "cmd", cfg.EventEvtCmd)
		assert.Equal(t, 8080, cfg.HTTPPort)
	})
}
