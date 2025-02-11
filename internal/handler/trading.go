package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"estore-trade/internal/domain"  // ドメインモデルをインポート
	"estore-trade/internal/usecase" // ユースケースインターフェースをインポート

	"go.uber.org/zap" //zapロガー
)

// TradingHandler は、取引関連のHTTPリクエストを処理するハンドラです。
type TradingHandler struct {
	tradingUsecase usecase.TradingUsecase
	logger         *zap.Logger // ロガーへのポインタ
}

// NewTradingHandler は、TradingHandlerの新しいインスタンスを作成します。
func NewTradingHandler(tradingUsecase usecase.TradingUsecase, logger *zap.Logger) *TradingHandler {
	return &TradingHandler{
		tradingUsecase: tradingUsecase,
		logger:         logger,
	}
}

// HandleTrade は、"/trade" エンドポイントへのPOSTリクエストを処理します (例)。
func (h *TradingHandler) HandleTrade(w http.ResponseWriter, r *http.Request) {
	// リクエストのコンテキストを取得 (キャンセルやタイムアウトの処理に使う)
	ctx := r.Context()

	// 1. リクエストボディのデコード (JSONをGoの構造体に変換)
	var orderRequest domain.Order
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 2. リクエストのバリデーション (必要に応じて)
	if err := validateOrderRequest(&orderRequest); err != nil {
		h.logger.Error("Invalid order request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. 認証情報などの取得 (通常はリクエストヘッダや認証トークンから)
	//    ここでは簡略化のため、固定値を使用。実際には、認証ミドルウェアなどを使う。
	userID := "your_user_id"    // 実際にはリクエストから取得
	password := "your_password" // 実際にはセキュアな方法で管理

	// 4. ユースケースの実行 (注文処理)
	placedOrder, err := h.tradingUsecase.PlaceOrder(ctx, userID, password, &orderRequest)
	if err != nil {
		h.logger.Error("Failed to place order", zap.Error(err))
		// エラーの種類に応じて適切なHTTPステータスコードを返す
		http.Error(w, "Failed to place order", http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	// 5. レスポンスの作成 (Goの構造体をJSONに変換)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created (注文が正常に作成された場合)
	if err := json.NewEncoder(w).Encode(placedOrder); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Order placed successfully", zap.String("order_id", placedOrder.ID))
}

// 他のハンドラ関数 (例: GetOrderStatusHandler, CancelOrderHandler) をここに追加

// validateOrderRequest は、注文リクエストのバリデーションを行います (例)。
func validateOrderRequest(order *domain.Order) error {
	// ここで、注文リクエストの内容をチェックします (例: 数量が正であるか、銘柄コードが有効かなど)。
	if order.Quantity <= 0 {
		return fmt.Errorf("invalid order quantity: %d", order.Quantity)
	}
	// 他のバリデーションルール...
	return nil
}
