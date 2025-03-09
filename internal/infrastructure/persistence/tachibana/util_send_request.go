package tachibana

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// sendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
func sendRequest(req *http.Request, maxRetries int) (map[string]interface{}, error) {
	var response map[string]interface{}

	// retryDoに渡す関数を定義。この関数がhttp clientの実行やデコード処理も行う
	retryFunc := func(client *http.Client, decodeFunc func([]byte, interface{}) error) (*http.Response, error) {

		resp, err := client.Do(req)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return resp, fmt.Errorf("API のステータスコードが200以外のためエラー: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // 読み込み終わったらすぐにクローズ
		if err != nil {
			return resp, fmt.Errorf("response body read error: %w", err)
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
