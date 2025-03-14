// internal/handler/trading.go
package handler

import (
	"encoding/json"
	"net/http"

	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// HandleTrade は、"/trade" エンドポイントへのPOSTリクエストを処理 (必要に応じて残す)
func (h *TradingHandler) HandleTrade(w http.ResponseWriter, r *http.Request) {
	// ... (既存のコード) ...
	// 1. リクエストボディのデコード (JSONをGoの構造体に変換)
	ctx := r.Context()
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
	// 認証情報はusecaseでconfigから取得するので、ここでは不要

	// 4. ユースケースの実行 (注文処理)
	// userID, password を削除
	placedOrder, err := h.tradingUsecase.PlaceOrder(ctx, &orderRequest)
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
	h.logger.Info("Order placed successfully", zap.String("order_id", placedOrder.UUID))
}
