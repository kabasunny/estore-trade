// internal/infrastructure/persistence/tachibana/mthd_download_master_data.go
package tachibana

import (
	"context"
	"encoding/json"
	"estore-trade/internal/domain"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// DownloadMasterData はマスタデータをダウンロードし、TachibanaClientImpleの各mapに格納
func (tc *TachibanaClientImple) DownloadMasterData(ctx context.Context) (*domain.MasterData, error) {
	// リクエスト直前に p_sd_date を生成
	now := time.Now()
	payload := map[string]string{
		"sCLMID":    clmidDownloadMasterData,
		"p_no":      tc.getPNo(),
		"p_sd_date": formatSDDate(now), // ここで現在時刻を使用
		"sJsonOfmt": "4",
		// 今回はすべて取得
		"sTargetCLMID": "CLMSystemStatus,CLMDateZyouhou,CLMYobine,CLMIssueMstKabu,CLMIssueSizyouMstKabu,CLMIssueSizyouKiseiKabu,CLMUnyouStatusKabu,CLMEventDownloadComplete",
		// "sTargetCLMID": "CLMSystemStatus,CLMDateZyouhou", // 絞り込み
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ペイロードのJSONマーシャルに失敗: %w", err)
	}

	// URLエンコード (GETリクエスト)
	encodedPayload := url.QueryEscape(string(payloadJSON))
	requestURL := tc.masterURL + "?" + encodedPayload

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, fmt.Errorf("マスターデータのリクエスト作成に失敗: %w", err)
	}
	// req.Header.Set("Content-Type", "application/json")

	req, cancel := withContextAndTimeout(req, 600*time.Second) // 10分
	defer cancel()

	// sendRequest(req, 3) // リトライ  sendRequestの戻り値が変更
	body, err := sendMasterDataRequest(req) //専用の関数に変更
	if err != nil {
		return nil, fmt.Errorf("マスターデータのダウンロードに失敗: %w", err)
	}

	// domain.MasterData インスタンスの作成
	m := &domain.MasterData{
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMap:                 make(map[string]domain.IssueMaster),
		IssueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		OperationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
	}

	// レスポンスの処理
	// responseはmap[string]interface{}型
	// 各マスタデータは、sCLMIDをキーとして格納されている
	//if err := processResponse(response, m, tc); err != nil {　変更前
	if err := processResponse(body, m, tc); err != nil { // sendRequestの戻り値を引数に
		return nil, err
	}

	// ダウンロードしたデータをTachibanaClientImpleにセット
	tc.masterData = m

	return m, nil // 成功
}
