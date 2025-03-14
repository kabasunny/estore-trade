// internal/infrastructure/persistence/tachibana/test_helpers.go
package tachibana

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"estore-trade/internal/config"
	"estore-trade/internal/domain"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// CreateTestClient はテスト用の TachibanaClientImple インスタンスを作成します。
func CreateTestClient(t *testing.T, md *domain.MasterData) *TachibanaClientImple {
	t.Helper()

	// .env ファイルのパスを修正
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get caller information")
	}
	// test_helpers.go から見た .env の相対パス (プロジェクトルート)
	envPath := filepath.Join(filepath.Dir(filename), "../../../../.env") // パスを修正

	// 設定ファイルの読み込み
	cfg, err := config.LoadConfig(envPath) // 絶対パスまたは相対パスを指定
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// ロガーの作成 (テスト用)
	logger := zaptest.NewLogger(t) // テストログを出力
	// logger := zap.NewNop()  // ログ出力を抑制する場合

	// デモ環境かどうかのチェックと表示
	if strings.Contains(cfg.TachibanaBaseURL, "demo") {
		fmt.Println("APIの　デモ　環境に接続")
	} else {
		fmt.Println("APIの　本番　環境に接続")
	}

	// TachibanaClientImple インスタンスの作成
	client := NewTachibanaClient(cfg, logger, md).(*TachibanaClientImple)

	return client
}

// GetClientFields は TachibanaClientImple のフィールドの値を map[string]string 形式で返します。
func GetClientFields(client *TachibanaClientImple) map[string]string {
	//muの追加
	client.mu.RLock()
	defer client.mu.RUnlock()
	return map[string]string{
		"baseURL":    client.baseURL.String(), // *url.URL は String() で文字列に
		"userID":     client.userID,           // 追加
		"password":   client.password,         // 追加
		"loggined":   fmt.Sprintf("%t", client.loggined),
		"requestURL": client.requestURL,
		"masterURL":  client.masterURL,
		"priceURL":   client.priceURL,
		"eventURL":   client.eventURL,
		"expiry":     client.expiry.Format(time.RFC3339Nano), // time.Time は Format() で文字列に
		"pNo":        fmt.Sprintf("%d", client.pNo),
		// 他のフィールドも必要に応じて追加
	}
}

// PrintClientFields は TachibanaClientImple のフィールドを整形して出力します。
func PrintClientFields(t *testing.T, client *TachibanaClientImple) {
	t.Helper()
	fields := GetClientFields(client) // フィールドの値を取得

	fmt.Println("TachibanaClientImple Fields:")
	for name, value := range fields {
		fmt.Printf("  %s: %s\n", name, value)
	}
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClientImple) GetBaseURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.baseURL.String() // 文字列で返す
}

// GetUserIDForTest はテスト用に userID を取得します。
func (tc *TachibanaClientImple) GetUserIDForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.userID
}

// GetPasswordForTest はテスト用に password を取得します。
func (tc *TachibanaClientImple) GetPasswordForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.password
}

// GetRequestURLForTest はテスト用に requestURL を取得 (テストヘルパー)
func (tc *TachibanaClientImple) GetRequestURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.requestURL
}

// GetLogginedForTest はテスト用にrequestURLを取得
func (tc *TachibanaClientImple) GetLogginedForTest() bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loggined
}

// SetBaseURLForTest はテスト用に baseURL を設定します。　削除
//func (tc *TachibanaClientImple) SetBaseURLForTest(baseURL string) {
//	tc.mu.Lock()
//	defer tc.mu.Unlock()
//	parsedURL, _ := url.Parse(baseURL)
//	tc.baseURL = parsedURL
//}

// SetUserIDForTest はテスト用に userID を設定します。
func (tc *TachibanaClientImple) SetUserIDForTest(userID string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.userID = userID
}

// SetPasswordForTest はテスト用に password を設定します。
func (tc *TachibanaClientImple) SetPasswordForTest(password string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.password = password
}

// SetRequestURLForTest はテスト用に requestURL を設定します。
func (tc *TachibanaClientImple) SetRequestURLForTest(requestURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.requestURL = requestURL
}

// GetMasterURLForTest はテスト用に masterURL を取得します。
func (tc *TachibanaClientImple) GetMasterURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.masterURL
}

// GetPositionSymbol はテスト用に Position の Symbol を取得します。
func GetPositionSymbol(p *domain.Position) string {
	return p.Symbol
}

// GetPositionQuantity はテスト用に Position の Quantity を取得します。
func GetPositionQuantity(p *domain.Position) int {
	return p.Quantity
}

// GetPositionID はテスト用に Position の ID (建玉番号) を取得します。
func GetPositionID(p *domain.Position) string {
	return p.ID
}

// GetOrderSymbol はテスト用に Order の Symbol を取得します。
func GetOrderSymbol(o *domain.Order) string {
	return o.Symbol
}

// GetOrderSide はテスト用に Order の Side を取得します。
func GetOrderSide(o *domain.Order) string {
	return o.Side
}

// GetOrderOrderType はテスト用に Order の OrderType を取得します。
func GetOrderOrderType(o *domain.Order) string {
	return o.OrderType
}

// GetOrderCondition はテスト用に Order の Condition を取得します。
func GetOrderCondition(o *domain.Order) string {
	return o.Condition
}

// GetOrderQuantity はテスト用に Order の Quantity を取得します。
func GetOrderQuantity(o *domain.Order) int {
	return o.Quantity
}

// GetOrderPrice はテスト用に Order の Price を取得します。
func GetOrderPrice(o *domain.Order) float64 {
	return o.Price
}

// GetOrderTriggerPrice はテスト用に Order の TriggerPriceを取得します
func GetOrderTriggerPrice(o *domain.Order) float64 {
	return o.TriggerPrice
}

// GetOrderMarketCode はテスト用に Order の MarketCode を取得します。
func GetOrderMarketCode(o *domain.Order) string {
	return o.MarketCode
}

// GetOrderTachibanaOrderID はテスト用に Order の TachibanaOrderID を取得します。
func GetOrderTachibanaOrderID(o *domain.Order) string {
	return o.TachibanaOrderID
}

// GetOrderStatus はテスト用に Order の Status を取得します。
func GetOrderStatus(o *domain.Order) string {
	return o.Status
}

// GetPNoForTest はテスト用に pNo を取得します。
func (tc *TachibanaClientImple) GetPNoForTest() string {
	tc.pNoMu.Lock()
	defer tc.pNoMu.Unlock()
	return strconv.FormatInt(tc.pNo, 10)
}

// GetConfig はテスト用に config.Config を取得します
func (tc *TachibanaClientImple) GetConfig() *config.Config {
	// .env ファイルのパスを修正
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get caller information") // テストヘルパー内なので panic で良い
	}
	// test_helpers.go から見た .env の相対パス (プロジェクトルート)
	envPath := filepath.Join(filepath.Dir(filename), "../../../../.env") // パスを修正

	// 設定ファイルの読み込み
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err)) // テストヘルパー内なので panic
	}
	return cfg
}

// GetLogger はLoggerを取得します
func (tc *TachibanaClientImple) GetLogger() *zap.Logger {
	return tc.logger
}

// CallParseEvent は、EventStream の parseEvent メソッドを呼び出すためのヘルパー関数です。
func CallParseEvent(es *EventStream, message []byte) (*domain.OrderEvent, error) {
	return es.parseEvent(message)
}

// NewTestEventStream は、テスト用の EventStream インスタンスを作成するヘルパー関数です。
func NewTestEventStream(logger *zap.Logger) *EventStream {
	return &EventStream{
		logger: logger, // 小文字の logger を使う
		// 他のフィールドは必要に応じて初期化 (今回は不要)
	}
}

// parseKPEvent は、KP メッセージをパースし、エラーがあればエラーを返す。
func parseKPEvent(message []byte) (*domain.OrderEvent, error) {
	fields := strings.Split(string(message), "\x01")

	event := &domain.OrderEvent{EventType: "KP"} // EventType を KP で初期化
	for _, field := range fields {
		keyValue := strings.SplitN(field, "\x02", 2)
		if len(keyValue) != 2 {
			continue // キーと値のペアでない場合はスキップ
		}
		key, value := keyValue[0], keyValue[1]

		switch key {
		case "p_no":
			event.EventNo = value
		case "p_date":
			t, err := time.Parse("2006.01.02-15:04:05.000", value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse p_date: %w", err)
			}
			event.Timestamp = t
		case "p_errno":
			//数値に変換して、0ならエラーなし
			code, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse p_errno to int: %w", err)
			}
			if code != 0 {
				if event.System == nil {
					event.System = &domain.SystemStatus{}
				}
				event.System.ErrNo = value
			}
		case "p_err":
			if event.System == nil {
				event.System = &domain.SystemStatus{}
			}
			event.System.ErrMsg = value
		case "p_cmd": //念の為
			if value != "KP" {
				return nil, fmt.Errorf("unexpected event type: %s", value)
			}
		}
	}
	return event, nil
}

// CallParseKPEvent は、テスト用に parseKPEvent メソッドを呼び出すためのヘルパー関数です。
func CallParseKPEvent(es *EventStream, message []byte) (*domain.OrderEvent, error) {
	return parseKPEvent(message)
}
