package controller

import (
	"net/http"

	inventoryusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/inventory"
)

// InventoryController は inventory コンテキストの HTTP エンドポイント群をまとめた
// プレゼンテーション層のコントローラである。
//
// CartController / CatalogController と同様、このコントローラはユースケース
// （アプリケーション層）だけに依存し、リポジトリや GORM の存在を一切知らない。
type InventoryController struct {
	setStock *inventoryusecase.SetStockUseCase
	getStock *inventoryusecase.GetStockUseCase
}

// NewInventoryController は InventoryController を生成する。
func NewInventoryController(
	setStock *inventoryusecase.SetStockUseCase,
	getStock *inventoryusecase.GetStockUseCase,
) *InventoryController {
	return &InventoryController{
		setStock: setStock,
		getStock: getStock,
	}
}

// Register は inventory コンテキストが提供するルートを mux に登録する。
//
// ここで登録するパス（/products/{id}/stock）は catalog コンテキストが
// 所有する /products/{id} 空間を拡張するものだが、catalog が登録する
// パターン（"GET /products/{id}" 等）とは末尾の "/stock" セグメントで
// 区別されるため、Go 1.22 以降の net/http のパターンマッチング上は
// 衝突しない（より詳細なパターンが優先されるルーティング規則により、
// 両コントローラが同じ ServeMux に安全に共存できる）。
func (c *InventoryController) Register(mux *http.ServeMux) {
	mux.HandleFunc("PUT /products/{id}/stock", c.handleSetStock)
	mux.HandleFunc("GET /products/{id}/stock", c.handleGetStock)
}

// setStockRequest は在庫数設定エンドポイントのリクエストボディである。
type setStockRequest struct {
	Quantity int `json:"quantity"`
}

func (c *InventoryController) handleSetStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req setStockRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	err := c.setStock.Execute(r.Context(), inventoryusecase.SetStockInput{
		ProductID: id,
		Quantity:  req.Quantity,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *InventoryController) handleGetStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	dto, err := c.getStock.Execute(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dto)
}
