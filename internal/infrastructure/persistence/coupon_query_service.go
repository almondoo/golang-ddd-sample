package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	couponusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/coupon"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CouponQuery は couponusecase.CouponQueryService の GORM 実装である。
//
// 型名を CouponQueryService としなかったのは、couponusecase パッケージの
// インターフェース名（couponusecase.CouponQueryService）とこの実装の型名が
// 同名になり紛らわしくなるのを避けるためである。command 側の
// CouponRepository が「coupon.Repository の実装」であるのと対比して、
// こちらは「couponusecase.CouponQueryService の実装」であることを var _ の
// アサーションで明示している。
type CouponQuery struct {
	db *gorm.DB
}

// NewCouponQuery は CouponQuery を生成する。
func NewCouponQuery(db *gorm.DB) *CouponQuery {
	return &CouponQuery{db: db}
}

// コンパイル時に CouponQuery が couponusecase.CouponQueryService を満たすことを保証する。
var _ couponusecase.CouponQueryService = (*CouponQuery)(nil)

// FindByCode は指定コードのクーポンを DTO として返す。
func (q *CouponQuery) FindByCode(ctx context.Context, code string) (*couponusecase.CouponDTO, error) {
	db := DBFromContext(ctx, q.db)

	var model CouponModel
	if err := db.WithContext(ctx).First(&model, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("coupon %s: %w", code, shared.ErrNotFound)
		}
		return nil, err
	}

	dto := toCouponDTO(model)
	return &dto, nil
}

// toCouponDTO は永続化モデルを問い合わせ用 DTO に変換する。
// ドメイン集約を経由しないため、ここでの変換はバリデーションを伴わない
// 単純なフィールドのコピーである。
func toCouponDTO(m CouponModel) couponusecase.CouponDTO {
	return couponusecase.CouponDTO{
		ID:           m.ID,
		Code:         m.Code,
		DiscountType: m.DiscountType,
		Amount:       m.Amount,
		RatePercent:  m.RatePercent,
		ExpiresAt:    m.ExpiresAt,
		UsageLimit:   m.UsageLimit,
		UsedCount:    m.UsedCount,
	}
}
