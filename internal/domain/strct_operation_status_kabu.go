package domain

// OperationStatusKabu 運用ステータス（株）(必要最低限)
type OperationStatusKabu struct {
	ListedMarket string `json:"sZyouzyouSizyou"` // 上場市場
	Unit         string `json:"sUnyouUnit"`      // 運用単位
	Status       string `json:"sUnyouStatus"`    // 運用ステータス
}
