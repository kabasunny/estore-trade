package tachibana

import (
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"net/url"

	"go.uber.org/zap"
)

// NewTachibanaClient 関数の定義
func NewTachibanaClient(cfg *config.Config, logger *zap.Logger, masterData *domain.MasterData) TachibanaClient {
	parsedURL, err := url.Parse(cfg.TachibanaBaseURL)
	if err != nil {
		logger.Fatal("Invalid Tachibana API base URL", zap.Error(err))
		return nil
	}
	return &TachibanaClientImple{
		baseURL:          parsedURL,
		apiKey:           cfg.TachibanaAPIKey,
		secret:           cfg.TachibanaAPISecret,
		logger:           logger,
		loggined:         false,
		pNo:              0,
		targetIssueCodes: make([]string, 0),
		masterData:       masterData, // 引数で受け取った masterData を設定
	}

}
