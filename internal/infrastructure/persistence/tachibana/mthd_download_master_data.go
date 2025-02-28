// internal/infrastructure/persistence/tachibana/mthd_download_master_data.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"estore-trade/internal/domain"
	"fmt"
	"net/http"
	"time"
)

// DownloadMasterData はマスタデータをダウンロードし、TachibanaClientImpleの
// 各mapに格納
func (tc *TachibanaClientImple) DownloadMasterData(ctx context.Context) (*domain.MasterData, error) {
	payload := map[string]string{
		"sCLMID": clmidDownloadMasterData,
		// 今回はすべて取得
		"sTargetCLMID": "CLMSystemStatus,CLMDateZyouhou,CLMYobine,CLMIssueMstKabu,CLMIssueSizyouMstKabu,CLMIssueSizyouKiseiKabu,CLMUnyouStatusKabu,CLMEventDownloadComplete",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ペイロードのJSONマーシャルに失敗: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.MasterURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, fmt.Errorf("マスターデータのリクエスト作成に失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req, cancel := withContextAndTimeout(req, 600*time.Second) // 10分
	defer cancel()

	response, err := sendRequest(req, 3) // リトライ
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
	if err := processResponse(response, m, tc); err != nil {
		return nil, err
	}

	// ダウンロードしたデータをTachibanaClientImpleにセット
	tc.MasterData = m

	return m, nil // 成功
}
