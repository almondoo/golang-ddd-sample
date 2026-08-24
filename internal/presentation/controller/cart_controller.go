package controller

import (
	"net/http"

	"github.com/almondoo/golang-ddd-sample/internal/application/cart/command"
	"github.com/almondoo/golang-ddd-sample/internal/application/cart/query"
)

// CartController は cart コンテキストの HTTP エンドポイント群をまとめた
// プレゼンテーション層のコントローラである。
//
// CatalogController と同様、このコントローラはユースケース（アプリケーション層）
// だけに依存し、リポジトリや GORM の存在を一切知らない。HTTP リクエストを
// ユースケースの入力形式に変換し、ユースケースの出力（または業務エラー）を
// HTTP レスポンスへ変換する橋渡しに徹する。
type CartController struct {
	addItem    *command.AddItemUseCase
	removeItem *command.RemoveItemUseCase
	getCart    *query.GetCartUseCase
}

// NewCartController は CartController を生成する。
func NewCartController(
	addItem *command.AddItemUseCase,
	removeItem *command.RemoveItemUseCase,
	getCart *query.GetCartUseCase,
) *CartController {
	return &CartController{
		addItem:    addItem,
		removeItem: removeItem,
		getCart:    getCart,
	}
}

// Register は cart コンテキストが提供するルートを mux に登録する。
// Go 1.22 以降の net/http のメソッド + パスパターンによるルーティングを使う。
func (c *CartController) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /carts/{customerID}", c.handleGetCart)
	mux.HandleFunc("POST /carts/{customerID}/items", c.handleAddItem)
	mux.HandleFunc("DELETE /carts/{customerID}/items/{productID}", c.handleRemoveItem)
}

func (c *CartController) handleGetCart(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	dto, err := c.getCart.Execute(r.Context(), customerID)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dto)
}

// addItemRequest はカートへの商品追加エンドポイントのリクエストボディである。
type addItemRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (c *CartController) handleAddItem(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var req addItemRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	err := c.addItem.Execute(r.Context(), command.AddItemInput{
		CustomerID: customerID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CartController) handleRemoveItem(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	productID := r.PathValue("productID")

	err := c.removeItem.Execute(r.Context(), command.RemoveItemInput{
		CustomerID: customerID,
		ProductID:  productID,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
