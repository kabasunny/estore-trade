// internal/infrastructure/persistence/tachibana/util_login.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func login(ctx context.Context, tc *TachibanaClientImple, userID, password string) (bool, error) {
	payload := map[string]string{
		"sCLMID":    clmidLogin,
		"sUserId":   userID,
		"sPassword": password,
		"p_no":      tc.getPNo(), // 呼び出し元でインクリメントするため、ここでは取得のみ
		"p_sd_date": formatSDDate(time.Now()),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("ログインペイロードのJSONエンコードに失敗しました", zap.Error(err))
		return false, fmt.Errorf("ログインペイロードのJSONエンコードに失敗しました: %w", err)
	}

	// URLを組み立て
	authURL, err := url.JoinPath(tc.baseURL.String(), "auth/") //baseURLにauth/を付け足す
	if err != nil {
		tc.logger.Error("認証URLの生成に失敗", zap.Error(err))
		return false, fmt.Errorf("認証URLの生成に失敗しました: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		tc.logger.Error("ログインリクエストの作成に失敗しました", zap.Error(err))
		return false, fmt.Errorf("ログインリクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3)
	if err != nil {
		return false, fmt.Errorf("ログインに失敗しました: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		tc.logger.Error("ログインAPIがエラーを返しました", zap.String("result_code", resultCode), zap.String("result_text", response["sResultText"].(string)))
		return false, fmt.Errorf("ログインAPIエラー: %s - %s", resultCode, response["sResultText"]) // 日本語でエラーを返す
	}

	requestURL, ok := response["sUrlRequest"]
	if !ok {
		tc.logger.Error("レスポンスにsUrlRequestが含まれていません")
		return false, fmt.Errorf("レスポンスにsUrlRequestが含まれていません")
	}
	masterURL, ok := response["sUrlMaster"]
	if !ok {
		tc.logger.Error("レスポンスにsUrlMasterが含まれていません")
		return false, fmt.Errorf("レスポンスにsUrlMasterが含まれていません")
	}
	priceURL, ok := response["sUrlPrice"]
	if !ok {
		tc.logger.Error("レスポンスにsUrlPriceが含まれていません")
		return false, fmt.Errorf("レスポンスにsUrlPriceが含まれていません")
	}
	eventURL, ok := response["sUrlEvent"]
	if !ok {
		tc.logger.Error("レスポンスにsUrlEventが含まれていません")
		return false, fmt.Errorf("レスポンスにsUrlEventが含まれていません")
	}

	// p_no はLogin APIのレスポンスで上書き
	if pNoStr, ok := response["p_no"].(string); ok {
		if pNo, err := strconv.ParseInt(pNoStr, 10, 64); err == nil {
			tc.pNo = pNo
		} else {
			tc.logger.Warn("p_noのパースに失敗しました", zap.String("p_no", pNoStr), zap.Error(err))
		}
	} else {
		tc.logger.Warn("レスポンスにp_noが含まれていないか、文字列ではありません", zap.Any("response", response))
	}

	tc.requestURL = requestURL.(string)
	tc.masterURL = masterURL.(string)
	tc.priceURL = priceURL.(string)
	tc.eventURL = eventURL.(string)
	tc.expiry = time.Now().Add(2 * time.Hour)
	tc.loggined = true

	return true, nil
}
