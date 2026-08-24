package controller

import (
	"net/http"

	"github.com/almondoo/golang-ddd-sample/internal/application/order/command"
	"github.com/almondoo/golang-ddd-sample/internal/application/order/query"
)

// OrderController は order コンテキストの HTTP エンドポイント群をまとめた
// プレゼンテーション層のコントローラである。
//
// CartController / CatalogController と同様、このコントローラはユースケース
// （アプリケーション層）だけに依存し、リポジトリや GORM の存在を
// 一切知らない。5 つのユースケース（登録・参照・支払い・発送・取消）を
// 束ねる薄い橋渡し役に徹する。
type OrderController struct {
	placeOrder  *command.PlaceOrderUseCase
	payOrder    *command.PayOrderUseCase
	shipOrder   *command.ShipOrderUseCase
	cancelOrder *command.CancelOrderUseCase
	getOrder    *query.GetOrderUseCase
}

// NewOrderController は OrderController を生成する。
func NewOrderController(
	placeOrder *command.PlaceOrderUseCase,
	payOrder *command.PayOrderUseCase,
	shipOrder *command.ShipOrderUseCase,
	cancelOrder *command.CancelOrderUseCase,
	getOrder *query.GetOrderUseCase,
) *OrderController {
	return &OrderController{
		placeOrder:  placeOrder,
		payOrder:    payOrder,
		shipOrder:   shipOrder,
		cancelOrder: cancelOrder,
		getOrder:    getOrder,
	}
}

// Register は order コンテキストが提供するルートを mux に登録する。
//
// 状態遷移（pay/ship/cancel）を専用のサブパスへの POST として表現して
// いるのは、これらが「注文というリソースの一部フィールドを PATCH する」
// のではなく「業務的に意味のある操作（コマンド）を実行する」ことを
// URL の形からも明確にするためである（いわゆる RPC 風のリソース設計）。
func (c *OrderController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", c.handlePlaceOrder)
	mux.HandleFunc("GET /orders/{id}", c.handleGetOrder)
	mux.HandleFunc("POST /orders/{id}/pay", c.handlePayOrder)
	mux.HandleFunc("POST /orders/{id}/ship", c.handleShipOrder)
	mux.HandleFunc("POST /orders/{id}/cancel", c.handleCancelOrder)
}

// placeOrderRequest は注文確定エンドポイントのリクエストボディである。
type placeOrderRequest struct {
	CustomerID string `json:"customerId"`
}

// placeOrderResponse は注文確定エンドポイントのレスポンスボディである。
type placeOrderResponse struct {
	OrderID string `json:"orderId"`
}

func (c *OrderController) handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req placeOrderRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	out, err := c.placeOrder.Execute(r.Context(), command.PlaceOrderInput{
		CustomerID: req.CustomerID,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, placeOrderResponse{OrderID: out.OrderID})
}

func (c *OrderController) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	dto, err := c.getOrder.Execute(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dto)
}

func (c *OrderController) handlePayOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.payOrder.Execute(r.Context(), command.PayOrderInput{OrderID: id}); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *OrderController) handleShipOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.shipOrder.Execute(r.Context(), command.ShipOrderInput{OrderID: id}); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *OrderController) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.cancelOrder.Execute(r.Context(), command.CancelOrderInput{OrderID: id}); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
