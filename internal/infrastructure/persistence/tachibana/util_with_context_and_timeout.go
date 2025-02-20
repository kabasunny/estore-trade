package tachibana

import (
	"context"
	"net/http"
	"time"
)

// withContextAndTimeout は、HTTP リクエストにコンテキストとタイムアウトを設定する
func withContextAndTimeout(req *http.Request, timeout time.Duration) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	return req.WithContext(ctx), cancel
}
