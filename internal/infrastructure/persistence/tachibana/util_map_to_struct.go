package tachibana

import (
	"encoding/json"
	"fmt"
)

// mapToStruct は、map[string]interface{} を構造体にマッピングする汎用関数
func mapToStruct(data map[string]interface{}, result interface{}) error {
	// fmt.Printf("DEBUG: mapToStruct input: %+v\n", data)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal failed: %v\n", err)
		return err
	}
	// fmt.Printf("DEBUG: JSON output: %s\n", string(b)) // 追加

	err = json.Unmarshal(b, result)
	if err != nil {
		fmt.Printf("ERROR: json.Unmarshal failed: %v\n", err)
	}
	// fmt.Printf("DEBUG: mapToStruct output: %+v, error: %v\n", result, err)
	return err
}

// OperationStatusKabu 構造体と map[string]map[string]OperationStatusKabu] 型のデータを使った具体例で、mapToStruct の処理の流れ

// 1. 元のデータ (map[string]interface{}):
// processResponse 関数で response (map[string]interface{}) から取り出されたデータ (data) は、以下のような map[string]interface{} 型のデータになっている
// (立花証券API からのレスポンスの一部を想定）

// dataMap (map[string]interface{}) の例
// dataMap := map[string]interface{}{
// 	"sZyouzyouSizyou": "00", // 上場市場コード (例: 東証)
// 	"sUnyouUnit":      "1",  // 運用単位 (例: 単元)
// 	"sUnyouStatus":    "1",  // 運用ステータス (例: 売買可能)
// }

// 2. mapToStruct 関数による変換:
// mapToStruct(dataMap, &operationStatusKabu) が呼び出されると、以下の処理が行われる
// json.Marshal(data):
// dataMap (map[string]interface{}) が JSON 形式のバイト列 ([]byte) に変換される
// dataMap を Marshal した結果 (JSON)
// {
//     "sZyouzyouSizyou": "00",
//     "sUnyouUnit": "1",
//     "sUnyouStatus": "1"
// }
// JSON 形式のバイト列が、result (ここでは &operationStatusKabu、つまり OperationStatusKabu 構造体へのポインタ) が指す構造体にデコード
// json.Unmarshal は、JSON のキーと構造体のフィールドタグ (json:"...") を照合して、値を対応するフィールドに格納

// operationStatusKabu (OperationStatusKabu 構造体) の例
// type OperationStatusKabu struct {
//     ListedMarket string `json:"sZyouzyouSizyou"` // 上場市場
//     Unit         string `json:"sUnyouUnit"`      // 運用単位
//     Status       string `json:"sUnyouStatus"`    // 運用ステータス
// }

// Unmarshal 後の operationStatusKabu の値
// operationStatusKabu = OperationStatusKabu{
//     ListedMarket: "00",
//     Unit:         "1",
//     Status:       "1",
// }
