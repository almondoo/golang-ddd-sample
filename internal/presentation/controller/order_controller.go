package controller

import (
	"net/http"

	orderusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/order"
)

// OrderController は order コンテキストの HTTP エンドポイント群をまとめた
// プレゼンテーション層のコントローラである。
//
// CartController / CatalogController と同様、このコントローラはユースケース
// （アプリケーション層）だけに依存し、リポジトリや GORM の存在を
// 一切知らない。5 つのユースケース（登録・参照・支払い・発送・取消）を
// 束ねる薄い橋渡し役に徹する。
type OrderController struct {
	placeOrder  *orderusecase.PlaceOrderUseCase
	payOrder    *orderusecase.PayOrderUseCase
	shipOrder   *orderusecase.ShipOrderUseCase
	cancelOrder *orderusecase.CancelOrderUseCase
	getOrder    *orderusecase.GetOrderUseCase
}

// NewOrderController は OrderController を生成する。
func NewOrderController(
	placeOrder *orderusecase.PlaceOrderUseCase,
	payOrder *orderusecase.PayOrderUseCase,
	shipOrder *orderusecase.ShipOrderUseCase,
	cancelOrder *orderusecase.CancelOrderUseCase,
	getOrder *orderusecase.GetOrderUseCase,
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
	// CouponCode は適用したいクーポンのコードである（任意項目）。
	// 空文字列であればクーポンを適用しない。
	CouponCode string `json:"couponCode"`
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

	out, err := c.placeOrder.Execute(r.Context(), orderusecase.PlaceOrderInput{
		CustomerID: req.CustomerID,
		CouponCode: req.CouponCode,
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

	if err := c.payOrder.Execute(r.Context(), orderusecase.PayOrderInput{OrderID: id}); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// shipOrderResponse は発送エンドポイントのレスポンスボディである。
// 生成した Shipment の ID を返す。生成したリソースの ID を返さないと、
// クライアントは配達完了 API に到達できない（PlaceOrder が orderId を
// 返すのと同じ理由）。
type shipOrderResponse struct {
	ShipmentID string `json:"shipmentId"`
}

func (c *OrderController) handleShipOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	out, err := c.shipOrder.Execute(r.Context(), orderusecase.ShipOrderInput{OrderID: id})
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, shipOrderResponse{ShipmentID: out.ShipmentID})
}

func (c *OrderController) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.cancelOrder.Execute(r.Context(), orderusecase.CancelOrderInput{OrderID: id}); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
