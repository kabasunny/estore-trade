// internal/config/config.go
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TachibanaAPIKey    string
	TachibanaAPISecret string
	TachibanaBaseURL   string
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	LogLevel           string
	// --- ここから修正 ---
	// EVENT I/F 接続用パラメータ追加
	EventRid     string // p_rid
	EventBoardNo string // p_board_no
	EventEvtCmd  string // p_evt_cmd
	// --- ここまで修正 ---
}

func LoadConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		// .envがなくても続行。システム環境変数を使用するため
		// return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		TachibanaAPIKey:    os.Getenv("TACHIBANA_API_KEY"),
		TachibanaAPISecret: os.Getenv("TACHIBANA_API_SECRET"),
		TachibanaBaseURL:   os.Getenv("TACHIBANA_BASE_URL"), // "https://kabuka.e-shiten.jp/e_api_v4r5" など
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             port,
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		LogLevel:           os.Getenv("LOG_LEVEL"), // 環境変数からログレベルを取得
		// --- ここから修正 ---
		// EVENT I/F 接続用パラメータ
		EventRid:     os.Getenv("EVENT_RID"),      // 例: "0"
		EventBoardNo: os.Getenv("EVENT_BOARD_NO"), // 例: "1000"
		EventEvtCmd:  os.Getenv("EVENT_EVT_CMD"),  // 例: "ST,KP,EC,SS,US"
		// --- ここまで修正 ---
	}, nil
}
