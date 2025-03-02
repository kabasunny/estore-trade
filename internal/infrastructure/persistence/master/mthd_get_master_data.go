// internal/infrastructure/persistence/master/mthd_get_master_data.go
package master

import (
	"context"
	"estore-trade/internal/domain"
)

func (r *masterDataRepository) GetMasterData(ctx context.Context) (*domain.MasterData, error) {
	masterData := &domain.MasterData{
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMap:                 make(map[string]domain.IssueMaster),
		IssueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		OperationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
	}

	// SystemStatus を取得
	systemStatus, err := r.getSystemStatus(ctx)
	if err != nil {
		return nil, err
	}
	masterData.SystemStatus = *systemStatus

	// DateInfoを取得
	dateInfo, err := r.getDateInfo(ctx)
	if err != nil {
		return nil, err
	}
	masterData.DateInfo = *dateInfo

	// CallPriceMapを取得
	callPrices, err := r.getAllCallPrices(ctx)
	if err != nil {
		return nil, err
	}
	for _, cp := range callPrices {
		masterData.CallPriceMap[cp.UnitNumber] = cp //修正
	}

	// IssueMapを取得
	issues, err := r.getAllIssueMasters(ctx)
	if err != nil {
		return nil, err
	}
	for _, issue := range issues {
		masterData.IssueMap[issue.IssueCode] = issue
	}

	// IssueMarketMapを取得
	issueMarkets, err := r.getAllIssueMarketMasters(ctx)
	if err != nil {
		return nil, err
	}
	for _, im := range issueMarkets {
		if _, ok := masterData.IssueMarketMap[im.IssueCode]; !ok {
			masterData.IssueMarketMap[im.IssueCode] = make(map[string]domain.IssueMarketMaster)
		}
		masterData.IssueMarketMap[im.IssueCode][im.MarketCode] = im
	}

	// IssueMarketRegulationMapを取得
	issueRegulations, err := r.getAllIssueMarketRegulations(ctx)
	if err != nil {
		return nil, err
	}
	for _, ir := range issueRegulations {
		if _, ok := masterData.IssueMarketRegulationMap[ir.IssueCode]; !ok {
			masterData.IssueMarketRegulationMap[ir.IssueCode] = make(map[string]domain.IssueMarketRegulation)
		}
		masterData.IssueMarketRegulationMap[ir.IssueCode][ir.ListedMarket] = ir
	}

	// OperationStatusKabuMapを取得
	operationStatuses, err := r.getAllOperationStatusKabu(ctx)
	if err != nil {
		return nil, err
	}
	for _, os := range operationStatuses {
		if _, ok := masterData.OperationStatusKabuMap[os.ListedMarket]; !ok {
			masterData.OperationStatusKabuMap[os.ListedMarket] = make(map[string]domain.OperationStatusKabu)
		}
		masterData.OperationStatusKabuMap[os.ListedMarket][os.Unit] = os
	}

	return masterData, nil
}
