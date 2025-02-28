// internal/infrastructure/persistence/tachibana/mthd_download_master_data.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"estore-trade/internal/domain" // 追加
	"fmt"
	"net/http"
	"time"
)

// DownloadMasterData はマスタデータをダウンロード
func (tc *TachibanaClientImple) DownloadMasterData(ctx context.Context) (*domain.MasterData, error) { // 戻り値の型を変更
	payload := map[string]string{
		"sCLMID":       clmidDownloadMasterData,
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

	req, cancel := withContextAndTimeout(req, 600*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3)
	if err != nil {
		return nil, fmt.Errorf("マスターデータのダウンロードに失敗: %w", err)
	}

	// MasterData インスタンスの作成
	m := &domain.MasterData{ //型を変更
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMap:                 make(map[string]domain.IssueMaster),
		IssueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		OperationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
	}

	if err := processResponse(response, m, tc); err != nil { //引数の型を変更
		return nil, err
	}

	tc.masterData = m // MasterData をセット
	return m, nil     // MasterDataを返す
}
