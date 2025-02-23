package config

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
