package controller

import (
	"net/http"

	"github.com/almondoo/golang-ddd-sample/internal/application/catalog/command"
	"github.com/almondoo/golang-ddd-sample/internal/application/catalog/query"
)

// CatalogController は catalog コンテキストの HTTP エンドポイント群をまとめた
// プレゼンテーション層のコントローラである。
//
// このコントローラはユースケース（アプリケーション層）だけに依存しており、
// リポジトリや GORM の存在を一切知らない。プレゼンテーション層の責務は
// あくまで「HTTP リクエストをユースケースの入力に変換し、ユースケースの
// 出力（または業務エラー）を HTTP レスポンスに変換する」という橋渡しに
// 限定している。controller は HTTP を usecase の入出力へ変換して呼び出す
// だけで、業務ルールを持たない。
type CatalogController struct {
	registerProduct *command.RegisterProductUseCase
	changePrice     *command.ChangePriceUseCase
	listProducts    *query.ListProductsUseCase
	getProduct      *query.GetProductUseCase
}

// NewCatalogController は CatalogController を生成する。
func NewCatalogController(
	registerProduct *command.RegisterProductUseCase,
	changePrice *command.ChangePriceUseCase,
	listProducts *query.ListProductsUseCase,
	getProduct *query.GetProductUseCase,
) *CatalogController {
	return &CatalogController{
		registerProduct: registerProduct,
		changePrice:     changePrice,
		listProducts:    listProducts,
		getProduct:      getProduct,
	}
}

// Register は catalog コンテキストが提供するルートを mux に登録する。
//
// Go 1.22 以降の net/http はメソッド + パスパターンによるルーティング
// （"POST /products" のような書式）と、{id} のようなパスパラメータを
// 標準ライブラリだけでサポートするようになった。本サンプルでは
// 学習目的でこれを採用し、外部ルーティングライブラリへの依存を避けている。
func (c *CatalogController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", c.handleRegisterProduct)
	mux.HandleFunc("GET /products", c.handleListProducts)
	mux.HandleFunc("GET /products/{id}", c.handleGetProduct)
	mux.HandleFunc("PUT /products/{id}/price", c.handleChangePrice)
}

// registerProductRequest は商品登録エンドポイントのリクエストボディである。
type registerProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceAmount int64  `json:"priceAmount"`
}

// registerProductResponse は商品登録エンドポイントのレスポンスボディである。
type registerProductResponse struct {
	ProductID string `json:"productId"`
}

func (c *CatalogController) handleRegisterProduct(w http.ResponseWriter, r *http.Request) {
	var req registerProductRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	out, err := c.registerProduct.Execute(r.Context(), command.RegisterProductInput{
		Name:        req.Name,
		Description: req.Description,
		PriceAmount: req.PriceAmount,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, registerProductResponse{ProductID: out.ProductID})
}

// listProductsResponse は商品一覧エンドポイントのレスポンスボディである。
type listProductsResponse struct {
	Products []query.ProductDTO `json:"products"`
}

func (c *CatalogController) handleListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.listProducts.Execute(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, listProductsResponse{Products: products})
}

func (c *CatalogController) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	product, err := c.getProduct.Execute(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, product)
}

// changePriceRequest は価格変更エンドポイントのリクエストボディである。
type changePriceRequest struct {
	NewPriceAmount int64 `json:"newPriceAmount"`
}

func (c *CatalogController) handleChangePrice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req changePriceRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	err := c.changePrice.Execute(r.Context(), command.ChangePriceInput{
		ProductID:      id,
		NewPriceAmount: req.NewPriceAmount,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
