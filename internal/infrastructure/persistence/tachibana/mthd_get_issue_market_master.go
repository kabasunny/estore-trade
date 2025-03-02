// internal/infrastructure/persistence/tachibana/mthd_get_issue_market_master.go
package tachibana

import "estore-trade/internal/domain"

// GetIssueMarketMaster は銘柄コードと市場コードに対応する市場情報を返します。
func (tc *TachibanaClientImple) GetIssueMarketMaster(issueCode, marketCode string) (domain.IssueMarketMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	// ターゲット銘柄リストが空でない、かつ、指定された銘柄コードが含まれていない場合のみ false を返す
	if len(tc.targetIssueCodes) > 0 && !contains(tc.targetIssueCodes, issueCode) {
		tc.targetIssueCodesMu.RUnlock()
		return domain.IssueMarketMaster{}, false // ターゲット銘柄でなければエラー
	}
	tc.targetIssueCodesMu.RUnlock()

	marketMap, ok := tc.masterData.IssueMarketMap[issueCode]
	if !ok {
		return domain.IssueMarketMaster{}, false
	}
	issueMarket, ok := marketMap[marketCode]
	return issueMarket, ok
}
