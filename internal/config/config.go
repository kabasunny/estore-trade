package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TachibanaAPIKey    string
	TachibanaAPISecret string
	TachibanaBaseURL   string
	TachibanaUserID    string
	TachibanaPassword  string
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	LogLevel           string
	EventRid           string // p_rid
	EventBoardNo       string // p_board_no
	EventEvtCmd        string // p_evt_cmd
	HTTPPort           int
}

func LoadConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		// .envファイルが見つからない場合にエラーメッセージを表示して終了
		fmt.Printf(".envファイルが見つかりませんでした: %v\n", err)
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	// HTTP_PORT を読み込む (エラー処理も追加)
	httpPortStr := os.Getenv("HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		// 数値に変換できない、または環境変数が設定されていない場合、デフォルト値を使用
		httpPort = 8080 // デフォルトポート
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
		HTTPPort:           httpPort, // 追加
	}, nil
}
