package controller

import (
	"net/http"
	"time"

	couponusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/coupon"
)

// CouponController は coupon コンテキストの HTTP エンドポイント群をまとめた
// プレゼンテーション層のコントローラである。
//
// CartController / CatalogController と同様、このコントローラはユースケース
// （アプリケーション層）だけに依存し、リポジトリや GORM の存在を一切知らない。
// HTTP リクエストをユースケースの入力形式に変換し、ユースケースの出力
// （または業務エラー）を HTTP レスポンスへ変換する橋渡しに徹する。
type CouponController struct {
	issueCoupon *couponusecase.IssueCouponUseCase
	getCoupon   *couponusecase.GetCouponUseCase
}

// NewCouponController は CouponController を生成する。
func NewCouponController(
	issueCoupon *couponusecase.IssueCouponUseCase,
	getCoupon *couponusecase.GetCouponUseCase,
) *CouponController {
	return &CouponController{
		issueCoupon: issueCoupon,
		getCoupon:   getCoupon,
	}
}

// Register は coupon コンテキストが提供するルートを mux に登録する。
// Go 1.22 以降の net/http のメソッド + パスパターンによるルーティングを使う。
func (c *CouponController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /coupons", c.handleIssueCoupon)
	mux.HandleFunc("GET /coupons/{code}", c.handleGetCoupon)
}

// issueCouponRequest はクーポン発行エンドポイントのリクエストボディである。
type issueCouponRequest struct {
	Code         string    `json:"code"`
	DiscountType string    `json:"discountType"`
	Amount       int64     `json:"amount"`
	RatePercent  int       `json:"ratePercent"`
	ExpiresAt    time.Time `json:"expiresAt"`
	UsageLimit   int       `json:"usageLimit"`
}

// issueCouponResponse はクーポン発行エンドポイントのレスポンスボディである。
type issueCouponResponse struct {
	CouponID string `json:"couponId"`
}

func (c *CouponController) handleIssueCoupon(w http.ResponseWriter, r *http.Request) {
	var req issueCouponRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteBadRequest(w, "invalid request body")
		return
	}

	out, err := c.issueCoupon.Execute(r.Context(), couponusecase.IssueCouponInput{
		Code:         req.Code,
		DiscountType: req.DiscountType,
		Amount:       req.Amount,
		RatePercent:  req.RatePercent,
		ExpiresAt:    req.ExpiresAt,
		UsageLimit:   req.UsageLimit,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, issueCouponResponse{CouponID: out.CouponID})
}

func (c *CouponController) handleGetCoupon(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	dto, err := c.getCoupon.Execute(r.Context(), code)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dto)
}
