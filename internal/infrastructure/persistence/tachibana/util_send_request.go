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
	retryFunc := func(client *http.Client, decodeFunc func(io.Reader, interface{}) error) (*http.Response, error) {

		resp, err := client.Do(req)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return resp, fmt.Errorf("API のステータスコードが200以外のためエラー: %d", resp.StatusCode)
		}

		if err := decodeFunc(resp.Body, &response); err != nil {
			resp.Body.Close()
			return resp, fmt.Errorf("レスポンスのデコードに失敗: %w", err)
		}
		return resp, nil
	}

	// デコード関数を定義 (Shift-JIS 固定)
	decodeFunc := func(r io.Reader, v interface{}) error {
		reader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
		return json.NewDecoder(reader).Decode(v)
	}

	// reqのTimeoutを使うので、ここではClientを生成しない
	resp, err := retryDo(retryFunc, maxRetries, 2*time.Second, &http.Client{}, decodeFunc) //空のClientを渡す
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return response, nil
}
