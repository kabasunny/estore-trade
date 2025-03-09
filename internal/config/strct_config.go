// internal/config/strct_config.go
package config

type Config struct {
	TachibanaBaseURL  string `env:"TACHIBANA_BASE_URL"` // 立花証券APIベースURL
	TachibanaUserID   string `env:"TACHIBANA_USER_ID"`  // 立花証券ユーザーID
	TachibanaPassword string `env:"TACHIBANA_PASSWORD"` // 立花証券パスワード

	DBHost       string `env:"DB_HOST"`        // データベースホスト名
	DBPort       int    `env:"DB_PORT"`        // データベースポート番号
	DBUser       string `env:"DB_USER"`        // データベースユーザー名
	DBPassword   string `env:"DB_PASSWORD"`    // データベースパスワード
	DBName       string `env:"DB_NAME"`        // データベース名
	LogLevel     string `env:"LOG_LEVEL"`      // ログレベル (debug, info, warn, error など)
	EventRid     string `env:"EVENT_RID"`      // EVENT I/F p_rid
	EventBoardNo string `env:"EVENT_BOARD_NO"` // EVENT I/F p_board_no
	EventEvtCmd  string `env:"EVENT_EVT_CMD"`  // EVENT I/F p_evt_cmd
	HTTPPort     int    `env:"HTTP_PORT"`      // HTTPサーバーポート番号
}
