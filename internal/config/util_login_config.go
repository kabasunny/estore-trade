// internal/config/util_login_config.go
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf(".envファイルが見つかりませんでした: %v\n", err)
		//return nil, err // .env が見つからなくても、環境変数が設定されていれば続行
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		if os.Getenv("DB_PORT") == "" { // 環境変数が空の場合
			dbPort = 5432 // デフォルト値を設定
		} else { // 環境変数が空でないが数値に変換できない場合
			return nil, fmt.Errorf("DB_PORT の値が不正です: %w", err)
		}
	}

	httpPortStr := os.Getenv("HTTP_PORT")
	httpPort := 8080       // デフォルトポート
	if httpPortStr != "" { // 環境変数が設定されている場合のみ変換を試みる
		httpPort, err = strconv.Atoi(httpPortStr)
		if err != nil {
			// 数値に変換できない場合、エラーとする
			return nil, fmt.Errorf("HTTP_PORT の値が不正です: %w", err)
		}
	}
	return &Config{
		TachibanaAPIKey:    os.Getenv("TACHIBANA_API_KEY"),
		TachibanaAPISecret: os.Getenv("TACHIBANA_API_SECRET"),
		TachibanaBaseURL:   os.Getenv("TACHIBANA_BASE_URL"),
		TachibanaUserID:    os.Getenv("TACHIBANA_USER_ID"),
		TachibanaPassword:  os.Getenv("TACHIBANA_PASSWORD"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             dbPort,
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		LogLevel:           os.Getenv("LOG_LEVEL"),
		EventRid:           os.Getenv("EVENT_RID"),
		EventBoardNo:       os.Getenv("EVENT_BOARD_NO"),
		EventEvtCmd:        os.Getenv("EVENT_EVT_CMD"),
		HTTPPort:           httpPort, //HTTPPort
	}, nil
}
