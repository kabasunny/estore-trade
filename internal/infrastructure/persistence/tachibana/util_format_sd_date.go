package tachibana

import (
	"time"
)

// formatSDDate は time.Time を YYYY.MM.DD-HH:MM:SS.TTT 形式の文字列に変換
func formatSDDate(t time.Time) string {
	return t.Format("2006.01.02-15:04:05.000")
}
