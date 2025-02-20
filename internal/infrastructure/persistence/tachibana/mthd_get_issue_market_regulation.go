package tachibana

// GetIssueMarketRegulation は銘柄コードと市場コードに対応する規制情報を返す
func (tc *TachibanaClientImple) GetIssueMarketRegulation(issueCode, marketCode string) (IssueMarketRegulation, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) {
			tc.targetIssueCodesMu.RUnlock()
			return IssueMarketRegulation{}, false // ターゲット銘柄でなければエラー
		}
	}
	tc.targetIssueCodesMu.RUnlock()

	marketMap, ok := tc.issueMarketRegulationMap[issueCode]
	if !ok {
		return IssueMarketRegulation{}, false
	}
	issueMarket, ok := marketMap[marketCode]
	return issueMarket, ok
}
