package tachibana

// マスタデータ保持用 (必要最低限に絞り込み)
type masterDataManager struct {
	systemStatus             SystemStatus
	dateInfo                 DateInfo
	callPriceMap             map[string]CallPrice
	issueMap                 map[string]IssueMaster
	issueMarketMap           map[string]map[string]IssueMarketMaster
	issueMarketRegulationMap map[string]map[string]IssueMarketRegulation
	operationStatusKabuMap   map[string]map[string]OperationStatusKabu
}
