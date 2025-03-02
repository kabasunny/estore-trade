// internal/infrastructure/persistence/tachibana/mthd_get_issue_market_regulation.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
)

// 銘柄コードと市場コードに対応する規制情報を返す
func (tc *TachibanaClientImple) GetIssueMarketRegulation(issueCode, marketCode string) (domain.IssueMarketRegulation, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	marketMap, ok := tc.masterData.IssueMarketRegulationMap[issueCode]
	if !ok {
		fmt.Printf("IssueCode %s not found in IssueMarketRegulationMap\n", issueCode)
		return domain.IssueMarketRegulation{}, false
	}

	issueMarketRegulation, ok := marketMap[marketCode]
	if !ok {
		fmt.Printf("MarketCode %s not found in marketMap for %s\n", marketCode, issueCode)
		return domain.IssueMarketRegulation{}, false
	}
	// ターゲット銘柄リストによるフィルタリング
	tc.targetIssueCodesMu.RLock()
	defer tc.targetIssueCodesMu.RUnlock()
	if len(tc.targetIssueCodes) > 0 && !contains(tc.targetIssueCodes, issueCode) {
		return domain.IssueMarketRegulation{}, false
	}

	return issueMarketRegulation, true
}
