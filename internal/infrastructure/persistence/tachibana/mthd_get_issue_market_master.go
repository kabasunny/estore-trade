package tachibana

// GetIssueMarketMaster は銘柄コードと市場コードに対応する市場情報を返します。
func (tc *TachibanaClientImple) GetIssueMarketMaster(issueCode, marketCode string) (IssueMarketMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) {
			tc.targetIssueCodesMu.RUnlock()
			return IssueMarketMaster{}, false // ターゲット銘柄でなければエラー
		}
	}
	tc.targetIssueCodesMu.RUnlock()

	marketMap, ok := tc.issueMarketMap[issueCode]
	if !ok {
		return IssueMarketMaster{}, false
	}
	issueMarket, ok := marketMap[marketCode]
	return issueMarket, ok
}
