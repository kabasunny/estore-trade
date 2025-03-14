package tachibana

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// parseEvent は、受信したメッセージをパースして domain.OrderEvent に変換
func (es *EventStream) parseEvent(message []byte) (*domain.OrderEvent, error) {
	es.logger.Debug("parseEvent: input message", zap.ByteString("message", message))

	// メッセージを行ごとに分割 (\n で分割)
	lines := strings.Split(string(message), "\n")

	event := &domain.OrderEvent{}
	var eventType string // イベントタイプを一時保存

	for _, line := range lines {
		// 空行はスキップ
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 各行を ^A (0x01) で分割
		fields := strings.Split(line, "\x01")
		for _, field := range fields {
			// 各フィールドを ^B (0x02) で key と value に分割
			keyValue := strings.SplitN(field, "\x02", 2)
			if len(keyValue) != 2 {
				es.logger.Debug("Skipping invalid field", zap.String("field", field))
				continue // キーと値のペアでない場合はスキップ
			}
			key, value := keyValue[0], keyValue[1]
			es.logger.Debug("Processing field", zap.String("key", key), zap.String("value", value))

			switch key {
			case "p_cmd": // コマンド (イベントタイプ)
				eventType = value // EventType をここで取得
				event.EventType = value
				es.logger.Debug("p_cmd: eventType set to", zap.String("eventType", eventType))

			case "p_no": // イベント内の連番(イベント番号ではない)
				//parseEvent内では使わない

			case "p_date": // イベント発生時刻
				t, err := time.Parse("2006.01.02-15:04:05.000", value)
				if err != nil {
					es.logger.Warn("Failed to parse p_date", zap.String("value", value), zap.Error(err))
					continue
				}
				event.Timestamp = t

			case "p_ENO": // イベント番号
				event.EventNo = value
				es.logger.Debug("p_ENO: EventNo set to", zap.String("eventNo", event.EventNo))

			case "p_errno":
				if event.System == nil { //Systemを初期化
					event.System = &domain.SystemStatus{}
				}
				event.System.ErrNo = value
			case "p_err":
				if event.System == nil { //Systemを初期化
					event.System = &domain.SystemStatus{}
				}
				event.System.ErrMsg = value

			// 他の共通項目もここに追加 (必要に応じて)

			default:
				// eventType ごとの処理
				switch eventType {
				case "EC": // 注文約定通知
					if event.Order == nil {
						event.Order = &domain.Order{}
					}
					switch key {
					case "p_ON": // 注文番号
						event.Order.TachibanaOrderID = value
					case "p_IC": // 銘柄コード
						event.Order.Symbol = value
					case "p_BBKB": // 売買区分
						switch value {
						case "1":
							event.Order.Side = "short"
						case "3":
							event.Order.Side = "long"
						case "5": //現渡
							event.Order.Side = "short"
						case "7": //現引
							event.Order.Side = "long"
						default:
							es.logger.Warn("Invalid p_BBKB value", zap.String("value", value))
						}
					case "p_ODST": // 注文ステータス
						event.Order.Status = value
					case "p_CRPR": //注文価格
						price, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse p_CRPR to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Order.Price = price
					case "p_CRSR": //注文数量
						quantity, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse p_CRSR to int", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Order.Quantity = quantity
					case "p_NT": // 通知種別  ★ここを修正★
						event.Order.NotificationType = value // NotificationType に設定
					case "p_ED": // 営業日
						event.Order.CreatedAt, _ = time.Parse("20060102", value)
					case "p_OON": // 親注文番号
						// event.Order にフィールドがないため、処理しない

					case "p_OT": // 注文種別
						// 使用しない

					case "p_ST": // 商品種別
						// event.Order にフィールドがないため、処理しない

					case "p_THKB": // 取引区分
						switch value {
						case "0":
							event.Order.TradeType = "spot" // 現物
						case "2":
							event.Order.TradeType = "credit_open" // 信用新規 (制度)
						case "4":
							event.Order.TradeType = "credit_close" // 信用返済 (制度)
						case "6":
							event.Order.TradeType = "credit_open" // 信用新規 (一般)
						case "8":
							event.Order.TradeType = "credit_close" // 信用返済 (一般)
						case "5": //現渡
							event.Order.TradeType = "spot" //現渡
						case "7": //現引
							event.Order.TradeType = "spot" //現引
						default:
							es.logger.Warn("Invalid p_THKB value", zap.String("value", value))
						}

					case "p_CRSJ": // 執行条件
						switch value {
						case "0":
							event.Order.ExecutionType = "" // なし
						case "2":
							event.Order.ExecutionType = "opening"
						case "4":
							event.Order.ExecutionType = "closing"
						case "6":
							event.Order.ExecutionType = "market"
						}

					case "p_CRPRKB": // 注文値段区分

					case "p_CRTKSR": // 取消数量
						canceledQuantity, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse p_CRTKSR to int", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Order.FilledQuantity -= canceledQuantity

					case "p_CREPSR": // 失効数量
						expiredQuantity, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse p_CREPSR to int", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Order.FilledQuantity -= expiredQuantity

					case "p_CREXSR": // 約定済数量
						executedQuantity, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse p_CREXSR to int", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Order.FilledQuantity = executedQuantity

					case "p_KOFG": // 繰越フラグ
					case "p_TTST": // 訂正取消ステータス
					case "p_EXST": // 約定ステータス
					case "p_LMIT": // 有効期限
						event.Order.ExpireAt, _ = time.Parse("20060102", value)
					case "p_JKK": //譲渡益課税区分
					case "p_CHNL": //チャネル
					case "p_EPRC": //失効理由コード
					case "p_EXPR": //約定値段
						executionPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse p_EXPR to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Order.AveragePrice = executionPrice

					case "p_EXSR": //約定数量
						executionVolume, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse p_EXSR to int", zap.String("value", value), zap.Error(err))
							continue //数値に変換できない場合は処理しない
						}
						event.Order.FilledQuantity = executionVolume
					case "p_EXRC": //取引所エラーコード
					case "p_EXDT": //通知日時
						event.Order.UpdatedAt, _ = time.Parse("20060102150405", value)

					case "p_IN": //銘柄名称
						// event.Order にフィールドがないため、処理しない

						// 訂正関連のフィールドは省略
					}

				case "NS": // ニュース通知
					if event.News == nil {
						event.News = &domain.News{}
					}
					switch key {
					case "p_ID": // ニュースID
						event.News.NewsID = value
					case "p_DT": // ニュース日付
						event.News.NewsDate = value
					case "p_TM": // ニュース時刻
						event.News.NewsTime = value
					case "p_CGN": // ニュースカテゴリ数
						// 使用しない
					case "p_CGL": // ニュースカテゴリリスト
						event.News.CategoryList = strings.Split(value, "^C")
					case "p_GRN": // ニュースジャンル数
						// 使用しない
					case "p_GRL": // ニュースジャンルリスト
						event.News.GenreList = strings.Split(value, "^C")
					case "p_ISN": // 関連銘柄コードリスト数
						// 使用しない
					case "p_ISL": // 関連銘柄コードリスト
						event.News.Symbols = strings.Split(value, "^C")
					case "p_SKF": // 使用しない
					case "p_UPD": // 使用しない
					case "p_HDL": // ニュースタイトル
						event.News.Title = value
					case "p_TX": // ニュース本文
						event.News.Body = value
					}

				case "SS": // システムステータス
					fallthrough // SSとUSは共通の構造体を使用
				case "US": // 運用ステータス
					if event.System == nil {
						event.System = &domain.SystemStatus{}
					}
					switch key {
					case "p_CT": // 情報更新時間
						event.System.UpdateDate = value
					case "p_LK": // ログイン許可区分 (SS)
						event.System.LoginPermission = value //既存のフィールド
						event.System.LoginStatus = value     //念の為残す
					case "p_SS": // システムステータス (SS)
						event.System.SystemState = value //既存のフィールド
						//event.System.SystemStatus = value //念の為残す  <-- この行を削除
					case "p_MC":
						event.System.MarketCode = value
					case "p_UC": // 運用カテゴリー (US)
						event.System.UpdateCategory = value
					case "p_UU": // 運用ユニット (US)
						event.System.UpdateUnit = value
					case "p_US": // 運用ステータス (US)
						event.System.UpdateStatus = value
					case "p_GSCD": //原資産コード
					case "p_SHSB": //商品種別
					case "p_EDK": //営業日区分
					}
				case "FD": // 時価情報
					if event.Market == nil {
						event.Market = &domain.Market{}
					}

					// p_col の値を解析
					parts := strings.Split(key, "_")
					if len(parts) != 3 {
						es.logger.Warn("Invalid p_col format", zap.String("key", key))
						continue
					}

					//dataType := parts[0]  // "p", "t", "x" などの型情報 (今回は使用しない)
					rowNumber := parts[1] // 行番号
					infoCode := parts[2]  // 情報コード
					event.Market.RowNumber = rowNumber

					switch infoCode {
					case "DPP": // 現在値
						currentPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DPP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.CurrentPrice = value //string
						event.Market.CurrentPrice = currentPrice // 直接代入

					case "DPG": // 現在値前値比較
						event.Market.PriceStatus = value
					case "DV": // 出来高
						turnover, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse DV to int", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.Turnover = value //string
						event.Market.Turnover = turnover // 直接代入

					case "BV", "GBV1", "GBV2", "GBV3", "GBV4", "GBV5", "GBV6", "GBV7", "GBV8", "GBV9", "GBV10": // 買気配数量
						bidQuantity, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse GBV to int", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.BidQuantity = value //string
						event.Market.BidQuantity = bidQuantity // 複数回設定される可能性があるが、最後の値を保持
					case "QBP", "GBP1", "GBP2", "GBP3", "GBP4", "GBP5", "GBP6", "GBP7", "GBP8", "GBP9", "GBP10": // 買気配値
						bidPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse GBP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.BidPrice = value //string
						event.Market.BidPrice = bidPrice // 複数回設定される可能性があるが、最後の値を保持
					case "AV", "GAV1", "GAV2", "GAV3", "GAV4", "GAV5", "GAV6", "GAV7", "GAV8", "GAV9", "GAV10": // 売気配数量
						askQuantity, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse GAV to int", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.AskQuantity = value //string
						event.Market.AskQuantity = askQuantity // 複数回設定される可能性があるが、最後の値を保持
					case "QAP", "GAP1", "GAP2", "GAP3", "GAP4", "GAP5", "GAP6", "GAP7", "GAP8", "GAP9", "GAP10": // 売気配値
						askPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse GAP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.AskPrice = value
						event.Market.AskPrice = askPrice // 複数回設定される可能性があるが、最後の値を保持
					case "DOP": // 始値
						openingPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DOP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.OpeningPrice = value //string
						event.Market.OpeningPrice = openingPrice
					case "DHP": // 高値
						highPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DHP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.HighPrice = value //string
						event.Market.HighPrice = highPrice
					case "DLP": // 安値
						lowPrice, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DLP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.LowPrice = value //string
						event.Market.LowPrice = lowPrice
					case "DJ": // 売買代金
						tradingVolume, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DJ to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						//event.Market.TradingVolume = value //string
						event.Market.TradingVolume = tradingVolume
					case "DHF": // 日通し高値フラグ
						event.Market.DailyHighStatus = value
					case "DLF": // 日通し安値フラグ
						event.Market.DailyLowStatus = value
					case "QAS": //売気配値種類
						event.Market.AskQuoteType = value
					case "QBS": //買気配値種類
						event.Market.BidQuoteType = value
					case "VWAP": //VWAP
						vwap, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse VWAP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.VWAP = vwap
					case "DYRP": //騰落率
						ratio, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DYRP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.TurnoverRatio = ratio
					case "PRP": //前日終値
						close, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse PRP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.PreviousClose = close
					case "DYWP": //前日比
						change, err := strconv.ParseFloat(value, 64)
						if err != nil {
							es.logger.Warn("Failed to parse DYWP to float64", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.PreviousChange = change
					case "LISS": //所属
						event.Market.Listing = value
					case "DHP:T": //高値時刻
						highTime, err := time.Parse("15:04", value)
						if err != nil {
							es.logger.Warn("Failed to parse DHP:T", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.HighTime = highTime
					case "DLP:T": //安値時刻
						LowTime, err := time.Parse("15:04", value)
						if err != nil {
							es.logger.Warn("Failed to parse DLP:T", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.LowTime = LowTime
					case "DOP:T": //始値時刻
						openTime, err := time.Parse("15:04", value)
						if err != nil {
							es.logger.Warn("Failed to parse DOP:T", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.OpenTime = openTime
					case "DPP:T": //現在値時刻
						currentTime, err := time.Parse("15:04", value)
						if err != nil {
							es.logger.Warn("Failed to parse DPP:T", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.CurrentTime = currentTime
					case "QOV": //売-OVER
						overSell, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse QOV to int", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.OverSell = overSell
					case "QUV": //買-UNDER
						underBuy, err := strconv.Atoi(value)
						if err != nil {
							es.logger.Warn("Failed to parse QUV to int", zap.String("value", value), zap.Error(err))
							continue
						}
						event.Market.UnderBuy = underBuy

					default:
						es.logger.Debug("Unhandled FD infoCode", zap.String("infoCode", infoCode))
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
