package tachibana

// IssueMaster 株式銘柄マスタ (必要最低限)
type IssueMaster struct {
	IssueCode   string `json:"sIssueCode"`         // 銘柄コード
	IssueName   string `json:"sIssueName"`         // 銘柄名称
	TradingUnit int    `json:"sBaibaiTani,string"` // 売買単位
	TokuteiF    string `json:"sTokuteiF"`          // 特定口座対象Ｃ

	// 他の情報は、必要になったら追加
}
