package tachibana

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// retryDo は、HTTP リクエストをリトライ付きで実行する
func retryDo(
	retryFunc func(*http.Client, func([]byte, interface{}) error) (*http.Response, error), // decodeFuncの型修正
	maxRetries int,
	initialBackoff time.Duration,
	client *http.Client, // http.Client を引数で渡す
	decodeFunc func([]byte, interface{}) error, // デコード関数を引数で渡す []byteに変更
) (*http.Response, error) {
	var resp *http.Response
	var err error

	for retries := 1; retries <= maxRetries; retries++ {
		resp, err = retryFunc(client, decodeFunc)

		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil // 成功時: エラーがなく、ステータスコードが200の場合
		}

		if retries < maxRetries {
			// 指数バックオフを計算
			// 回数が増すごとに間隔が広くなる
			// 初期遅延時間に対して2の乗数でリトライ間隔を増加 (例: 2秒, 4秒, 8秒...)
			backoff := time.Duration(math.Pow(2, float64(retries))) * initialBackoff
			// 計算したリトライ間隔の時間だけ待機
			time.Sleep(backoff)

			// レスポンスが存在し、かつそのボディがまだ閉じられていない場合は閉じる
			// これはリソースリークを防ぐための重要なステップ
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}
	}

	if resp != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request failed after %d retries: last error: %v, last status code: %d", maxRetries+1, err, resp.StatusCode)
	}
	return nil, fmt.Errorf("HTTP request failed after %d retries: last error: %w", maxRetries+1, err)
}
