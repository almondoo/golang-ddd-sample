package controller

import (
	"net/http"

	shippingusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/shipping"
)

// ShippingController は shipping コンテキストの HTTP エンドポイント群を
// まとめたプレゼンテーション層のコントローラである。
//
// OrderController / CartController と同様、このコントローラはユースケース
// （アプリケーション層）だけに依存し、リポジトリや GORM の存在を
// 一切知らない。2 つのユースケース（配達完了・参照）を束ねる薄い
// 橋渡し役に徹する。
type ShippingController struct {
	deliverShipment *shippingusecase.DeliverShipmentUseCase
	getShipment     *shippingusecase.GetShipmentUseCase
}

// NewShippingController は ShippingController を生成する。
func NewShippingController(
	deliverShipment *shippingusecase.DeliverShipmentUseCase,
	getShipment *shippingusecase.GetShipmentUseCase,
) *ShippingController {
	return &ShippingController{
		deliverShipment: deliverShipment,
		getShipment:     getShipment,
	}
}

// Register は shipping コンテキストが提供するルートを mux に登録する。
//
// 状態遷移（deliver）を専用のサブパスへの POST として表現しているのは、
// order コンテキストの pay/ship/cancel と同じ理由（「配送というリソースの
// 一部フィールドを PATCH する」のではなく「業務的に意味のある操作
// （コマンド）を実行する」ことを URL の形からも明確にするため）である。
func (c *ShippingController) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /shipments/{id}", c.handleGetShipment)
	mux.HandleFunc("POST /shipments/{id}/deliver", c.handleDeliverShipment)
}

func (c *ShippingController) handleGetShipment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	dto, err := c.getShipment.Execute(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dto)
}

func (c *ShippingController) handleDeliverShipment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.deliverShipment.Execute(r.Context(), shippingusecase.DeliverShipmentInput{ShipmentID: id}); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
