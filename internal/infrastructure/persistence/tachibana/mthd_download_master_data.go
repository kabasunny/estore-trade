package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DownloadMasterData はマスタデータをダウンロード
func (tc *TachibanaClientImple) DownloadMasterData(ctx context.Context) error {
	payload := map[string]string{ // 必要最低限のマスタデータを要求
		"sCLMID":       clmidDownloadMasterData,
		"sTargetCLMID": "CLMSystemStatus,CLMDateZyouhou,CLMYobine,CLMIssueMstKabu,CLMIssueSizyouMstKabu,CLMIssueSizyouKiseiKabu,CLMUnyouStatusKabu,CLMEventDownloadComplete",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ペイロードのJSONマーシャルに失敗: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.masterURL, bytes.NewBuffer(payloadJSON)) // マスターデータのリクエスト作成
	if err != nil {
		return fmt.Errorf("マスターデータのリクエスト作成に失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // HTTPリクエストヘッダーに "Content-Type" を設定し、ボディの内容が JSON フォーマットであることをサーバーに通知

	req, cancel := withContextAndTimeout(req, 600*time.Second) // 長めのタイムアウト付きでリクエストを送信
	defer cancel()

	response, err := sendRequest(req, 3) // 3回リトライ設定し送信
	if err != nil {
		return fmt.Errorf("マスターデータのダウンロードに失敗: %w", err)
	}

	// 不要なロックを削除
	m := &masterDataManager{ // マスタデータマネージャーの初期化 mapは初期化しておかないとパニックになる
		callPriceMap:             make(map[string]CallPrice),
		issueMap:                 make(map[string]IssueMaster),
		issueMarketMap:           make(map[string]map[string]IssueMarketMaster),
		issueMarketRegulationMap: make(map[string]map[string]IssueMarketRegulation),
		operationStatusKabuMap:   make(map[string]map[string]OperationStatusKabu),
	}

	if err := processResponse(response, m, tc); err != nil { // レスポンスデータの処理
		return err
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()

	// マスタデータの更新
	tc.systemStatus = m.systemStatus
	tc.dateInfo = m.dateInfo
	tc.issueMap = m.issueMap
	tc.callPriceMap = m.callPriceMap
	tc.issueMarketMap = m.issueMarketMap
	tc.issueMarketRegulationMap = m.issueMarketRegulationMap
	tc.operationStatusKabuMap = m.operationStatusKabuMap

	// ターゲット銘柄の設定　指定された銘柄コードのみを対象とするようにマスタデータをフィルタリングする
	if len(tc.targetIssueCodes) > 0 {
		if err := tc.SetTargetIssues(ctx, tc.targetIssueCodes); err != nil {
			return err
		}
	}

	return nil
}
