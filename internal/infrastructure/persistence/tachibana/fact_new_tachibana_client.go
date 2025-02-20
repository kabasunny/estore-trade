package tachibana

import (
	"estore-trade/internal/config"
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
		baseURL:                  parsedURL,
		apiKey:                   cfg.TachibanaAPIKey,
		secret:                   cfg.TachibanaAPISecret,
		logger:                   logger,
		loggined:                 false,
		pNo:                      0,
		callPriceMap:             make(map[string]CallPrice),
		issueMap:                 make(map[string]IssueMaster),
		issueMarketMap:           make(map[string]map[string]IssueMarketMaster),
		issueMarketRegulationMap: make(map[string]map[string]IssueMarketRegulation),
		operationStatusKabuMap:   make(map[string]map[string]OperationStatusKabu),
		targetIssueCodes:         make([]string, 0),
	}
}
