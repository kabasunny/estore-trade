package tachibana

import (
	"fmt"
	"strings"
	"time"

	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// parseEvent は、受信したメッセージをパースして domain.OrderEvent に変換 (メインロジック)
func (es *EventStream) parseEvent(message []byte) (*domain.OrderEvent, error) {
	es.logger.Debug("parseEvent: input message", zap.ByteString("message", message))

	lines := strings.Split(string(message), "\n")
	event := &domain.OrderEvent{}
	var eventType string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\x01")
		for _, field := range fields {
			keyValue := strings.SplitN(field, "\x02", 2)
			if len(keyValue) != 2 {
				es.logger.Debug("Skipping invalid field", zap.String("field", field))
				continue
			}
			key, value := keyValue[0], keyValue[1]
			es.logger.Debug("Processing field", zap.String("key", key), zap.String("value", value))

			switch key {
			case "p_cmd":
				eventType = value
				event.EventType = value
				es.logger.Debug("p_cmd: eventType set to", zap.String("eventType", eventType))

			case "p_PV":
				event.Provider = value

			case "p_no":
				//parseEvent内では使わない
			case "p_date":
				t, err := time.Parse("2006.01.02-15:04:05.000", value)
				if err != nil {
					es.logger.Warn("Failed to parse p_date", zap.String("value", value), zap.Error(err))
					continue
				}
				event.Timestamp = t

			case "p_ENO":
				event.EventNo = value
				es.logger.Debug("p_ENO: EventNo set to", zap.String("eventNo", event.EventNo))

			case "p_errno":
				if event.System == nil {
					event.System = &domain.SystemStatus{}
				}
				event.System.ErrNo = value
			case "p_err":
				if event.System == nil {
					event.System = &domain.SystemStatus{}
				}
				event.System.ErrMsg = value
			default:
				// eventType ごとの処理を呼び出す
				switch eventType {
				case "EC":
					if err := es.parseEC(event, key, value); err != nil {
						return nil, err
					}
				case "NS":
					if err := es.parseNS(event, key, value); err != nil {
						return nil, err
					}
				case "SS", "US":
					if err := es.parseSSUS(event, key, value); err != nil {
						return nil, err
					}
				case "FD":
					if err := es.parseFD(event, key, value); err != nil {
						return nil, err
					}
				case "ST": //STはp_errnoとp_errの処理するため、ここでは処理しない。
				case "KP": //KPはp_errnoとp_errの処理するため、ここでは処理しない。
				case "RR": // アプリケーション専用のため、処理しない。
				case "FC": // アプリケーション専用のため、処理しない。
				default:
					es.logger.Warn("Unknown event type", zap.String("eventType", eventType))
				}
			}
		}
	}

	if eventType == "" {
		return nil, fmt.Errorf("event type is empty: %s", string(message))
	}
	es.logger.Debug("parseEvent: returning event", zap.Any("event", event))
	return event, nil
}
