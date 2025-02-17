// internal/infrastructure/persistence/tachibana/master_data.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// マスタデータ保持用
type masterDataManager struct {
	systemStatus SystemStatus
	dateInfo     DateInfo
	callPriceMap map[string]CallPrice   // 呼値 (Key: sYobineTaniNumber)
	issueMap     map[string]IssueMaster // 銘柄マスタ（株式）(Key: 銘柄コード)
}

func (tc *TachibanaClientImple) DownloadMasterData(ctx context.Context) error {
	payload := map[string]string{
		"sCLMID":    clmidDownloadMasterData,  // マスタデータダウンロード用のsCLMID
		"p_no":      tc.getPNo(),              // v4.6以降不要か
		"p_sd_date": formatSDDate(time.Now()), // v4.6以降不要か
	}

	payloadJSON, err := json.Marshal(payload) // ここで payload を使用
	if err != nil {
		return fmt.Errorf("failed to marshal master data request payload: %w", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.masterURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create master data request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// コンテキストとタイムアウトの設定
	req, cancel := withContextAndTimeout(req, 600*time.Second) // マスタデータダウンロードは時間がかかる可能性があるため、長めに設定
	defer cancel()
	client := &http.Client{}    //  client := &http.Client{Timeout: 600 * time.Second}
	resp, err := client.Do(req) // リクエスト送信
	if err != nil {
		return fmt.Errorf("failed to send master data request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("master data API returned non-200 status code: %d", resp.StatusCode)
	}
	// Shift-JISからUTF-8に変換するリーダーを作成
	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())

	// マスタデータマネージャーの初期化
	tc.mu.Lock()
	m := &masterDataManager{ // masterDataManager のインスタンスを作成
		callPriceMap: make(map[string]CallPrice),
		issueMap:     make(map[string]IssueMaster),
	}
	tc.mu.Unlock()

	decoder := json.NewDecoder(reader)
	for { // マスタデータは複数回に分けて送られてくる可能性があるため
		var response map[string]interface{}
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break // データの終わりに達したらループを抜ける
			}
			return fmt.Errorf("failed to decode master data response: %w", err)
		}

		// sCLMID でどのマスタデータか判別
		sCLMID, ok := response["sCLMID"].(string)
		if !ok {
			tc.logger.Error("sCLMID not found in master data response") //追加
			return fmt.Errorf("sCLMID not found in master data response")
		}

		switch sCLMID {
		case "CLMSystemStatus":
			var systemStatus SystemStatus
			if err := mapToStruct(response, &systemStatus); err != nil {
				return fmt.Errorf("failed to map SystemStatus: %w", err)
			}
			tc.mu.Lock()
			m.systemStatus = systemStatus // m に格納
			tc.mu.Unlock()
		case "CLMDateZyouhou":
			var dateInfo DateInfo
			if err := mapToStruct(response, &dateInfo); err != nil {
				return fmt.Errorf("failed to map DateInfo: %w", err)
			}
			tc.mu.Lock()
			m.dateInfo = dateInfo // m に格納
			tc.mu.Unlock()
		case "CLMYobine":
			var callPrice CallPrice
			if err := mapToStruct(response, &callPrice); err != nil {
				return fmt.Errorf("failed to map CallPrice: %w", err)
			}
			// 呼値はキー(sYobineTaniNumber)でmapに格納
			tc.mu.Lock()
			m.callPriceMap[strconv.Itoa(callPrice.UnitNumber)] = callPrice // m に格納
			tc.mu.Unlock()

		case "CLMIssueMstKabu": // 銘柄マスタ(株式)
			var issueMaster IssueMaster
			if err := mapToStruct(response, &issueMaster); err != nil {
				return fmt.Errorf("failed to map IssueMaster: %w", err)
			}
			// 銘柄コードをキーにしてmapに格納
			tc.mu.Lock()
			m.issueMap[issueMaster.IssueCode] = issueMaster // m に格納
			tc.mu.Unlock()
		case "CLMEventDownloadComplete": // 初期ダウンロード完了通知
			// ダウンロード完了時に、m から tc にデータをコピー
			tc.mu.Lock()
			tc.systemStatus = m.systemStatus
			tc.dateInfo = m.dateInfo
			tc.issueMap = m.issueMap
			tc.callPriceMap = m.callPriceMap
			tc.mu.Unlock()
			return nil
			// 他のマスタデータも同様に処理
		default:
			tc.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID)) // ログは記録
		}
	}
	// ここで、ループを抜けた後にも tc にデータをコピーするように修正（データが途中で終わる場合に対応）
	tc.mu.Lock()
	tc.systemStatus = m.systemStatus
	tc.dateInfo = m.dateInfo
	tc.issueMap = m.issueMap
	tc.callPriceMap = m.callPriceMap
	tc.mu.Unlock()

	return nil
}

// mapToStruct は、map[string]interface{} を構造体にマッピングする汎用関数
func mapToStruct(data map[string]interface{}, result interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, result)
}

// (追加) マスタデータへのアクセス用メソッド (Getter) の追加:
func (tc *TachibanaClientImple) GetSystemStatus() SystemStatus {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.systemStatus
}

func (tc *TachibanaClientImple) GetDateInfo() DateInfo {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.dateInfo
}

// 呼値取得
func (tc *TachibanaClientImple) GetCallPrice(unitNumber string) (CallPrice, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	callPrice, ok := tc.callPriceMap[unitNumber]
	return callPrice, ok
}

// 銘柄情報取得
func (tc *TachibanaClientImple) GetIssueMaster(issueCode string) (IssueMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	issue, ok := tc.issueMap[issueCode]
	return issue, ok
}
