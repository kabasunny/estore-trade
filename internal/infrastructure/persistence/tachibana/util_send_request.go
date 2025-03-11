// internal/infrastructure/persistence/tachibana/util_send_request.go
package tachibana

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// sendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
func sendRequest(req *http.Request, maxRetries int) (map[string]interface{}, error) {
	var response map[string]interface{}
	// 元のリクエストURLを保持
	originalURL := req.URL.String()
	// retryDoに渡す関数を定義。この関数がhttp clientの実行やデコード処理も行う
	retryFunc := func(client *http.Client, decodeFunc func([]byte, interface{}) error) (*http.Response, error) {
		// リトライごとにリクエストを再生成
		req, err := http.NewRequest(req.Method, originalURL, nil) //req.Clone(req.Context())
		if err != nil {
			return nil, fmt.Errorf("failed to recreate request for retry: %w", err)
		}
		// req.Header.Set("Content-Type", "application/json") // Content-Type を再設定

		// リクエスト内容をログに出力 (デバッグ用)
		fmt.Printf("Request URL: %s\n", req.URL.String())
		fmt.Printf("Request Method: %s\n", req.Method)

		// リクエストボディをログ出力 (Shift-JIS -> UTF-8 変換)
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 再設定
			bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), bodyBytes)
			if err != nil {
				fmt.Printf("Request Body Decode Error: %v\n", err)
			}
			fmt.Printf("Request Body (UTF-8): %s\n", string(bodyUTF8))
		}
		// fmt.Printf("Request Headers: %v\n", req.Header)

		resp, err := client.Do(req) //clientは、http.Client{}
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return resp, fmt.Errorf("API のステータスコードが200以外のためエラー: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		// fmt.Print(body)
		resp.Body.Close() // 読み込み終わったらすぐにクローズ
		if err != nil {
			return resp, fmt.Errorf("response body read error: %w", err)
		}

		// レスポンスボディをログに出力 (Shift-JIS -> UTF-8, 1行ずつ)
		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
		if err != nil {
			fmt.Printf("Raw Response Body Decode Error: %v\n", err)
			return resp, err //デコードに失敗したら、エラーを返す
		}
		fmt.Println("Raw Response Body (UTF-8, one line per JSON):")
		scanner := bufio.NewScanner(strings.NewReader(string(bodyUTF8)))
		for scanner.Scan() {
			line := scanner.Text()
			// 空行はスキップ
			if strings.TrimSpace(line) == "" {
				continue
			}
			// JSON として有効か確認 (簡易チェック)
			var js map[string]interface{}
			if json.Unmarshal([]byte(line), &js) == nil {
				fmt.Println(line)
			} else {
				fmt.Println("  Invalid JSON:", line)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning response body: %v\n", err)
		}

		if err := decodeFunc(body, &response); err != nil {
			return resp, fmt.Errorf("レスポンスのデコードに失敗: %w", err)
		}
		return resp, nil
	}

	// デコード関数を定義 (Shift-JIS から UTF-8 への変換)
	decodeFunc := func(body []byte, v interface{}) error { // 引数を io.Reader から []byte に変更
		// Shift-JISからUTF-8への変換
		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
		if err != nil {
			return fmt.Errorf("shift-jis decode error: %w", err)
		}
		return json.Unmarshal(bodyUTF8, v) // UTF-8 でデコード
	}

	// reqのTimeoutを使うので、ここではClientを生成しない
	resp, err := retryDo(retryFunc, maxRetries, 2*time.Second, &http.Client{}, decodeFunc) //空のClientを渡す
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return response, nil
}
