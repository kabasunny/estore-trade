// utils.go
package tachibana

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// formatSDDate は time.Time を YYYY.MM.DD-HH:MM:SS.TTT 形式の文字列に変換
func formatSDDate(t time.Time) string {
	return t.Format("2006.01.02-15:04:05.000")
}

// withContextAndTimeout は、HTTP リクエストにコンテキストとタイムアウトを設定する
func withContextAndTimeout(req *http.Request, timeout time.Duration) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	return req.WithContext(ctx), cancel
}

// retryDo は、HTTP リクエストをリトライ付きで実行する
// retryFunc: 実際に行うHTTPリクエストを含む関数
// maxRetries: 最大リトライ回数
// initialBackoff: 最初のリトライまでの待ち時間
func retryDo(retryFunc func() (*http.Response, error), maxRetries int, initialBackoff time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		resp, err = retryFunc() // retryFunc() は、HTTPリクエストを送信し、(*http.Response, error)を返す関数

		if err == nil && resp.StatusCode == http.StatusOK {
			// 成功したらレスポンスを返す
			return resp, nil
		}

		// エラーの場合 (またはステータスコードが200でない場合)、リトライ
		if retries < maxRetries {
			// 指数バックオフを計算 (例: 2秒, 4秒, 8秒...)
			backoff := time.Duration(math.Pow(2, float64(retries))) * initialBackoff
			time.Sleep(backoff)

			// もしrespが存在し、Bodyがまだ閉じられていなければ閉じる。
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		} else {
			break // 最大リトライ回数に達したらループを抜ける
		}
	}

	// 最大リトライ回数を超えても成功しなかった場合
	if resp != nil {
		resp.Body.Close() // Bodyを閉じてリソースリークを防ぐ
		return nil, fmt.Errorf("HTTP request failed after %d retries: last error: %v, last status code: %d", maxRetries+1, err, resp.StatusCode)
	}
	return nil, fmt.Errorf("HTTP request failed after %d retries: last error: %w", maxRetries+1, err)
}

// isValidPrice は、注文価格が呼値の単位に従っているかをチェックする関数
func isValidPrice(price float64, callPrice CallPrice) bool {
	prices := [20]float64{
		callPrice.Price1, callPrice.Price2, callPrice.Price3, callPrice.Price4, callPrice.Price5,
		callPrice.Price6, callPrice.Price7, callPrice.Price8, callPrice.Price9, callPrice.Price10,
		callPrice.Price11, callPrice.Price12, callPrice.Price13, callPrice.Price14, callPrice.Price15,
		callPrice.Price16, callPrice.Price17, callPrice.Price18, callPrice.Price19, callPrice.Price20,
	}
	unitPrices := [20]float64{
		callPrice.UnitPrice1, callPrice.UnitPrice2, callPrice.UnitPrice3, callPrice.UnitPrice4, callPrice.UnitPrice5,
		callPrice.UnitPrice6, callPrice.UnitPrice7, callPrice.UnitPrice8, callPrice.UnitPrice9, callPrice.UnitPrice10,
		callPrice.UnitPrice11, callPrice.UnitPrice12, callPrice.UnitPrice13, callPrice.UnitPrice14, callPrice.UnitPrice15,
		callPrice.UnitPrice16, callPrice.UnitPrice17, callPrice.UnitPrice18, callPrice.UnitPrice19, callPrice.UnitPrice20,
	}

	for i := 0; i < len(prices); i++ {
		if price <= prices[i] {
			remainder := math.Mod(price, unitPrices[i])
			return remainder == 0
		}
	}
	return false // ここには到達しないはずだが、念のため
}

// contains は、スライスに特定の要素が含まれているかどうかをチェックするヘルパー関数
func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
