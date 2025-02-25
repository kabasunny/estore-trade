// internal/domain/strct_issue_master.go
package domain

// IssueMaster 株式銘柄マスタ (必要最低限)
type IssueMaster struct {
	IssueCode   string `json:"issue_code"`   // 銘柄コード
	IssueName   string `json:"issue_name"`   // 銘柄名称
	TradingUnit int    `json:"trading_unit"` // 売買単位
	TokuteiF    string `json:"tokutei_f"`    // 特定口座対象Ｃ

	// 他の情報は、必要になったら追加
}
