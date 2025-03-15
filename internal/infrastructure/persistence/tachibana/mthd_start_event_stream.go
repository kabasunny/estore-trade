package tachibana

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

// StartSample は EVENT I/F への接続を確立し、メッセージ受信ループを開始 (サンプル実装)
func (es *EventStream) Start(ctx context.Context) error {
	baseEventURL, err := es.tachibanaClient.GetEventURL()
	if err != nil {
		return fmt.Errorf("failed to get event URL: %w", err)
	}

	eventURL, err := url.Parse(baseEventURL)
	if err != nil {
		es.logger.Error("Failed to parse event base URL", zap.Error(err))
		return fmt.Errorf("failed to parse event base URL: %w", err)
	}

	values := url.Values{}
	values.Add("p_rid", es.config.EventRid)
	values.Add("p_board_no", es.config.EventBoardNo)
	values.Add("p_eno", "0")
	values.Add("p_evt_cmd", es.config.EventEvtCmd)
	eventURL.RawQuery = values.Encode()

	es.logger.Info("EventStream: eventURL", zap.String("url", eventURL.String()))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, eventURL.String(), nil)
	if err != nil {
		es.logger.Error("Failed to create event stream request", zap.Error(err))
		return fmt.Errorf("failed to create event stream request: %w", err)
	}

	// イベント受信ループを別のゴルーチンで開始
	go es.startEventLoop(ctx, req)

	return nil
}
