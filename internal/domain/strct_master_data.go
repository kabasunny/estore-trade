// internal/domain/strct_master_data.go
package domain

// MasterData は、システムで扱うマスタデータを保持する構造体
type MasterData struct {
	SystemStatus             SystemStatus // システム状態
	DateInfo                 DateInfo     // 日付情報
	CallPriceMap             map[string]CallPrice
	IssueMap                 map[string]IssueMaster
	IssueMarketMap           map[string]map[string]IssueMarketMaster     // 銘柄コード -> 市場コード -> IssueMarketMaster
	IssueMarketRegulationMap map[string]map[string]IssueMarketRegulation // 銘柄コード -> 市場コード -> IssueMarketRegulation
	OperationStatusKabuMap   map[string]map[string]OperationStatusKabu   // 上場市場 -> 運用単位 -> OperationStatusKabu
}
