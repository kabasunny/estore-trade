// internal/infrastructure/persistence/tachibana/mthd_logout.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logout はAPIからログアウトします。
func (tc *TachibanaClientImple) Logout(ctx context.Context) error {
	tc.mu.Lock() // 他の操作と競合しないようにロック
	defer tc.mu.Unlock()

	// ログインしていない場合は何もしない (エラーにもしない)
	if !tc.loggined {
		return nil
	}

	payload := map[string]string{
		"sCLMID":    clmidLogoutRequest,
		"p_no":      tc.getPNo(),
		"p_sd_date": formatSDDate(time.Now()),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("ログアウトペイロードのJSONエンコードに失敗しました", zap.Error(err))
		return fmt.Errorf("ログアウトペイロードのJSONエンコードに失敗しました: %w", err)
	}

	// requestURL をそのまま使用.  認証のURL(/auth/)ではなく、通常の取引で使うURLを使う
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		tc.logger.Error("ログアウトリクエストの作成に失敗しました", zap.Error(err))
		return fmt.Errorf("ログアウトリクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3) // sendRequest 関数でリトライ処理を行う
	if err != nil {
		return fmt.Errorf("ログアウトに失敗しました: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		// 警告コードがある場合もログに出力 (PlaceOrder, GetOrderStatus に倣う)
		warnCode, _ := response["sWarningCode"].(string)
		warnText, _ := response["sWarningText"].(string)

		tc.logger.Error("ログアウトAPIがエラーを返しました",
			zap.String("result_code", resultCode),
			zap.String("result_text", response["sResultText"].(string)), // resultText を使用
			zap.String("warning_code", warnCode),
			zap.String("warning_text", warnText),
		)
		return fmt.Errorf("ログアウトAPIエラー: %s - %s", resultCode, response["sResultText"]) // 日本語でエラーを返す
	}
	tc.loggined = false     // ログアウト状態にする
	tc.requestURL = ""      // キャッシュをクリア
	tc.masterURL = ""       // キャッシュをクリア
	tc.priceURL = ""        // キャッシュをクリア
	tc.eventURL = ""        // キャッシュをクリア
	tc.expiry = time.Time{} // 有効期限を過去に設定

	return nil
}
