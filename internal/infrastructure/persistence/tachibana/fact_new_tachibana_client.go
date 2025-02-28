package tachibana

import (
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"net/url"

	"go.uber.org/zap"
)

// NewTachibanaClient 関数の定義
func NewTachibanaClient(cfg *config.Config, logger *zap.Logger) TachibanaClient {
	parsedURL, err := url.Parse(cfg.TachibanaBaseURL)
	if err != nil {
		logger.Fatal("Invalid Tachibana API base URL", zap.Error(err))
		return nil
	}
	return &TachibanaClientImple{
		BaseURL:                  parsedURL,
		ApiKey:                   cfg.TachibanaAPIKey,
		Secret:                   cfg.TachibanaAPISecret,
		Logger:                   logger,
		Loggined:                 false,
		PNo:                      0,
		callPriceMap:             make(map[string]domain.CallPrice),
		issueMap:                 make(map[string]domain.IssueMaster),
		issueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		issueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		operationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
		targetIssueCodes:         make([]string, 0),
	}
}
