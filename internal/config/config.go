package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// 環境変数から設定情報を読み込む機能

// 設定情報を保持
type Config struct {
	TachibanaAPIKey    string
	TachibanaAPISecret string
	TachibanaBaseURL   string // Base URL を追加
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	LogLevel           string // ログレベルを追加
}

// .env ファイルから環境変数を読み込み、Config 構造体に格納して返す
func LoadConfig(envPath string) (*Config, error) {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load(envPath)
	if err != nil {
		//.envがなくても続行。システム環境変数を使用するため
		//return nil, err // エラーを返さない
	}

	// 環境変数から設定値を読み込む
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		TachibanaAPIKey:    os.Getenv("TACHIBANA_API_KEY"),
		TachibanaAPISecret: os.Getenv("TACHIBANA_API_SECRET"),
		TachibanaBaseURL:   os.Getenv("TACHIBANA_BASE_URL"), // Base URL を追加
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             port,
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		LogLevel:           os.Getenv("LOG_LEVEL"), // 環境変数からログレベルを取得
	}, nil
}
