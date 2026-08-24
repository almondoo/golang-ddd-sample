package persistence

import (
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CouponModel は coupons テーブルに対応する GORM 用の構造体である。
//
// ドメインの Coupon 集約とは意図的に別の型として定義している。
// gorm タグはあくまで永続化の都合であり、ドメイン層の型にこれを持ち込むと
// 「ドメインが ORM の存在を知っている」ことになってしまうためである。
//
// Amount / RatePercent はどちらか一方だけが DiscountType に応じて意味を
// 持つ（ドメイン層の Coupon 構造体と同じ設計）。rate 型のクーポンでは
// Amount = 0, Currency = "" のまま永続化される。
type CouponModel struct {
	ID           string `gorm:"primaryKey"`
	Code         string `gorm:"uniqueIndex"`
	DiscountType string
	Amount       int64
	Currency     string
	RatePercent  int
	ExpiresAt    time.Time
	UsageLimit   int
	UsedCount    int
}

// TableName は GORM に対して物理テーブル名を明示する。
func (CouponModel) TableName() string {
	return "coupons"
}

// toDomain は永続化モデルからドメイン集約を復元する。
// DB から読み出した値は過去にドメイン層のバリデーションを通過済みという
// 前提のもと、ReconstructCoupon（検証を行わない再構築コンストラクタ）を使う。
func (m CouponModel) toDomain() (*coupon.Coupon, error) {
	id, err := coupon.NewCouponID(m.ID)
	if err != nil {
		return nil, err
	}
	code, err := coupon.NewCouponCode(m.Code)
	if err != nil {
		return nil, err
	}
	discountType, err := coupon.NewDiscountType(m.DiscountType)
	if err != nil {
		return nil, err
	}
	amount, err := shared.NewMoney(m.Amount, shared.Currency(m.Currency))
	if err != nil {
		return nil, err
	}
	return coupon.ReconstructCoupon(id, code, discountType, amount, m.RatePercent, m.ExpiresAt, m.UsageLimit, m.UsedCount), nil
}

// couponModelFromDomain はドメイン集約から永続化モデルを組み立てる。
func couponModelFromDomain(c *coupon.Coupon) CouponModel {
	return CouponModel{
		ID:           c.ID().String(),
		Code:         c.Code().String(),
		DiscountType: c.Type().String(),
		Amount:       c.Amount().Amount(),
		Currency:     string(c.Amount().Currency()),
		RatePercent:  c.RatePercent(),
		ExpiresAt:    c.ExpiresAt(),
		UsageLimit:   c.UsageLimit(),
		UsedCount:    c.UsedCount(),
	}
}
