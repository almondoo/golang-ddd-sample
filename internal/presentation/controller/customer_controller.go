package controller

import (
	"net/http"

	customerusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/customer"
)

// CustomerController は customer コンテキストの HTTP エンドポイント群を
// まとめたプレゼンテーション層のコントローラである。
//
// CartController / OrderController と同様、このコントローラはユースケース
// （アプリケーション層）だけに依存し、リポジトリや GORM の存在を
// 一切知らない。5 つのユースケース（登録・参照・住所追加・デフォルト
// 住所変更・住所削除）を束ねる薄い橋渡し役に徹する。
type CustomerController struct {
	registerCustomer     *customerusecase.RegisterCustomerUseCase
	addAddress           *customerusecase.AddAddressUseCase
	changeDefaultAddress *customerusecase.ChangeDefaultAddressUseCase
	removeAddress        *customerusecase.RemoveAddressUseCase
	getCustomer          *customerusecase.GetCustomerUseCase
}

// NewCustomerController は CustomerController を生成する。
func NewCustomerController(
	registerCustomer *customerusecase.RegisterCustomerUseCase,
	addAddress *customerusecase.AddAddressUseCase,
	changeDefaultAddress *customerusecase.ChangeDefaultAddressUseCase,
	removeAddress *customerusecase.RemoveAddressUseCase,
	getCustomer *customerusecase.GetCustomerUseCase,
) *CustomerController {
	return &CustomerController{
		registerCustomer:     registerCustomer,
		addAddress:           addAddress,
		changeDefaultAddress: changeDefaultAddress,
		removeAddress:        removeAddress,
		getCustomer:          getCustomer,
	}
}

// Register は customer コンテキストが提供するルートを mux に登録する。
// Go 1.22 以降の net/http のメソッド + パスパターンによるルーティングを使う。
//
// デフォルト住所変更を「住所リソースの部分更新」ではなく専用サブパスへの
// PUT として表現しているのは、これが「isDefault フィールドを PATCH する」
// のではなく「デフォルトを切り替える」という業務的に意味のある操作（コマンド）
// であることを URL の形からも明確にするためである
// （order コンテキストの pay/ship/cancel と同じ設計判断）。
func (c *CustomerController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /customers", c.handleRegisterCustomer)
	mux.HandleFunc("GET /customers/{id}", c.handleGetCustomer)
	mux.HandleFunc("POST /customers/{id}/addresses", c.handleAddAddress)
	mux.HandleFunc("PUT /customers/{id}/addresses/{addressID}/default", c.handleChangeDefaultAddress)
	mux.HandleFunc("DELETE /customers/{id}/addresses/{addressID}", c.handleRemoveAddress)
}

// registerCustomerRequest は顧客登録エンドポイントのリクエストボディである。
type registerCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// registerCustomerResponse は顧客登録エンドポイントのレスポンスボディである。
type registerCustomerResponse struct {
	CustomerID string `json:"customerId"`
}

func (c *CustomerController) handleRegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req registerCustomerRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	out, err := c.registerCustomer.Execute(r.Context(), customerusecase.RegisterCustomerInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, registerCustomerResponse{CustomerID: out.CustomerID})
}

func (c *CustomerController) handleGetCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	dto, err := c.getCustomer.Execute(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dto)
}

// addAddressRequest は住所追加エンドポイントのリクエストボディである。
type addAddressRequest struct {
	PostalCode string `json:"postalCode"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Line       string `json:"line"`
}

// addAddressResponse は住所追加エンドポイントのレスポンスボディである。
type addAddressResponse struct {
	AddressID string `json:"addressId"`
}

func (c *CustomerController) handleAddAddress(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req addAddressRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	out, err := c.addAddress.Execute(r.Context(), customerusecase.AddAddressInput{
		CustomerID: id,
		PostalCode: req.PostalCode,
		Prefecture: req.Prefecture,
		City:       req.City,
		Line:       req.Line,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, addAddressResponse{AddressID: out.AddressID})
}

func (c *CustomerController) handleChangeDefaultAddress(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	addressID := r.PathValue("addressID")

	err := c.changeDefaultAddress.Execute(r.Context(), customerusecase.ChangeDefaultAddressInput{
		CustomerID: id,
		AddressID:  addressID,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CustomerController) handleRemoveAddress(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	addressID := r.PathValue("addressID")

	err := c.removeAddress.Execute(r.Context(), customerusecase.RemoveAddressInput{
		CustomerID: id,
		AddressID:  addressID,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
