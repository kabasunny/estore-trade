// internal/infrastructure/persistence/tachibana/tests/place_order_test.go
package tachibana_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTachibanaClientImple_PlaceOrder(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t) // tachibana パッケージの SetupTestClient を呼び出す
	defer httpmock.DeactivateAndReset()

	t.Run("Success - Buy Order Limit", func(t *testing.T) { //指値
		mockOrderID := "12345"

		httpmock.RegisterResponder("POST", "https://example.com/request", //tc.RequestURL
			func(req *http.Request) (*http.Response, error) {
				// リクエストボディの検証 (JSON デコード)
				var reqBody map[string]interface{}
				if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
					return httpmock.NewStringResponse(http.StatusBadRequest, "Invalid JSON"), nil
				}

				// 必須パラメータの存在チェック(簡易版)
				// ConvertOrderToPlaceOrderPayload で必須パラメータはセットされるはずなので、
				// ここでは、ConvertOrderToPlaceOrderPayload が正しく呼ばれているか確認できれば良い
				requiredParams := []string{"sCLMID", "sIssueCode", "sBaibaiKubun", "sOrderPrice", "sOrderSuryou", "sCondition"}
				for _, param := range requiredParams {
					if _, ok := reqBody[param]; !ok {
						return httpmock.NewStringResponse(http.StatusBadRequest, fmt.Sprintf("Missing parameter: %s", param)), nil
					}
				}

				//sOrderPriceが文字列になっているか確認
				if _, ok := reqBody["sOrderPrice"].(string); !ok {
					return httpmock.NewStringResponse(http.StatusBadRequest, "sOrderPrice should be string"), nil
				}

				// リクエスト内容に応じたモックレスポンスを作成
				respBody := map[string]interface{}{
					"sResultCode":  "0",
					"sOrderNumber": mockOrderID,
				}
				respJSON, _ := json.Marshal(respBody)
				return httpmock.NewBytesResponse(http.StatusOK, respJSON), nil
			},
		)

		order := &domain.Order{
			Symbol:    "7974",  // 任天堂
			Side:      "buy",   // 買い注文
			OrderType: "limit", // 指値
			Price:     7000,    // 指値
			Quantity:  1,       // 数量
			Status:    "pending",
		}

		returnedOrder, err := client.PlaceOrder(context.Background(), order)
		require.NoError(t, err)
		assert.Equal(t, mockOrderID, returnedOrder.ID)
		assert.Equal(t, "pending", returnedOrder.Status) //初期状態
	})

	t.Run("Success - Buy Order Market", func(t *testing.T) { //成行
		mockOrderID := "67890"

		httpmock.RegisterResponder("POST", "https://example.com/request",
			func(req *http.Request) (*http.Response, error) {
				var reqBody map[string]interface{}
				if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
					return httpmock.NewStringResponse(http.StatusBadRequest, "Invalid JSON"), nil
				}
				requiredParams := []string{"sCLMID", "sIssueCode", "sBaibaiKubun", "sOrderSuryou", "sCondition"}
				for _, param := range requiredParams {
					if _, ok := reqBody[param]; !ok {
						return httpmock.NewStringResponse(http.StatusBadRequest, fmt.Sprintf("Missing parameter: %s", param)), nil
					}
				}

				// sOrderPriceが"0"であることを確認
				if price, ok := reqBody["sOrderPrice"].(string); !ok || price != "0" {
					return httpmock.NewStringResponse(http.StatusBadRequest, "sOrderPrice should be 0 for market order"), nil
				}

				respBody := map[string]interface{}{
					"sResultCode":  "0",
					"sOrderNumber": mockOrderID,
				}
				respJSON, _ := json.Marshal(respBody)
				return httpmock.NewBytesResponse(http.StatusOK, respJSON), nil
			},
		)

		order := &domain.Order{
			Symbol:    "7974",
			Side:      "buy",
			OrderType: "market", // 成行
			Quantity:  1,
			Status:    "pending",
		}

		returnedOrder, err := client.PlaceOrder(context.Background(), order)
		require.NoError(t, err)
		assert.Equal(t, mockOrderID, returnedOrder.ID)
		assert.Equal(t, "pending", returnedOrder.Status)
	})
	t.Run("Success - Buy Order Stop", func(t *testing.T) {
		mockOrderID := "54321"

		httpmock.RegisterResponder("POST", "https://example.com/request",
			func(req *http.Request) (*http.Response, error) {
				var reqBody map[string]interface{}
				if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
					return httpmock.NewStringResponse(http.StatusBadRequest, "Invalid JSON"), nil
				}

				// 簡易チェック
				requiredParams := []string{"sCLMID", "sIssueCode", "sBaibaiKubun", "sOrderSuryou", "sCondition", "sGyakusasiOrderType", "sGyakusasiZyouken", "sGyakusasiPrice"}
				for _, param := range requiredParams {
					if _, ok := reqBody[param]; !ok {
						return httpmock.NewStringResponse(http.StatusBadRequest, fmt.Sprintf("Missing parameter: %s", param)), nil
					}
				}

				// sOrderPrice が "*" であることを確認
				if price, ok := reqBody["sOrderPrice"].(string); !ok || price != "*" {
					return httpmock.NewStringResponse(http.StatusBadRequest, "sOrderPrice should be * for stop order"), nil
				}

				if _, ok := reqBody["sGyakusasiZyouken"].(string); !ok {
					return httpmock.NewStringResponse(http.StatusBadRequest, "sGyakusasiZyouken should be string"), nil
				}
				if _, ok := reqBody["sGyakusasiPrice"].(string); !ok {
					return httpmock.NewStringResponse(http.StatusBadRequest, "sGyakusasiPrice should be string"), nil
				}

				respBody := map[string]interface{}{
					"sResultCode":  "0",
					"sOrderNumber": mockOrderID,
				}
				respJSON, _ := json.Marshal(respBody)
				return httpmock.NewBytesResponse(http.StatusOK, respJSON), nil
			},
		)

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "buy",
			OrderType:    "stop", //逆指値
			Quantity:     1,
			Status:       "pending",
			TriggerPrice: 7010, //トリガー価格
			Price:        7000,
		}

		returnedOrder, err := client.PlaceOrder(context.Background(), order)
		require.NoError(t, err)
		assert.Equal(t, mockOrderID, returnedOrder.ID)
		assert.Equal(t, "pending", returnedOrder.Status)
	})

	//OrderTypeはセットしなくて良い
	t.Run("API Error - Invalid Issue Code", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://example.com/request",
			httpmock.NewStringResponder(http.StatusOK, `{"sResultCode": "E999", "sResultText": "Invalid issue code", "sWarningCode": "W123", "sWarningText": "Check issue code"}`),
		)

		order := &domain.Order{
			Symbol:    "XXXX", // 無効な銘柄コード
			Side:      "buy",
			Price:     1000,
			OrderType: "limit",
		}

		_, err := client.PlaceOrder(context.Background(), order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "order API returned an error")
		assert.Contains(t, err.Error(), "E999") // エラーコードの確認
	})

	//OrderTypeはセットしなくて良い
	t.Run("API Error - Order Number Missing", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://example.com/request",
			httpmock.NewStringResponder(http.StatusOK, `{"sResultCode": "0"}`), // sOrderNumber がない
		)
		order := &domain.Order{
			Symbol:    "7974",
			Side:      "buy",
			Price:     1000,
			OrderType: "limit",
		}

		_, err := client.PlaceOrder(context.Background(), order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "order number not found")
	})

	//OrderTypeはセットしなくて良い
	t.Run("Network Error", func(t *testing.T) {
		httpmock.Reset() // Responder をクリアしてネットワークエラーを発生させる
		order := &domain.Order{
			Symbol:    "7974",
			Side:      "buy",
			Price:     1000,
			OrderType: "limit",
		}
		_, err := client.PlaceOrder(context.Background(), order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "place order failed")

	})
	t.Run("Convert Error", func(t *testing.T) {
		//OrderTypeはセットしなくて良い
		order := &domain.Order{
			Symbol:    "7974",
			Side:      "invalid", //無効なside
			Price:     1000,
			OrderType: "limit",
		}
		_, err := client.PlaceOrder(context.Background(), order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert order to payload: invalid order side")
	})
	// 他のテストケースも同様に追加...
}
