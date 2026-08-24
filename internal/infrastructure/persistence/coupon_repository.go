package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CouponRepository は coupon.Repository の GORM 実装である。
//
// アプリケーション層は coupon.Repository（インターフェース）にのみ依存し、
// この構造体（具体的な GORM 実装）の存在を知らない。
type CouponRepository struct {
	db *gorm.DB
}

// NewCouponRepository は CouponRepository を生成する。
func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

// コンパイル時に CouponRepository が coupon.Repository を満たすことを保証する。
var _ coupon.Repository = (*CouponRepository)(nil)

// FindByCode は指定コードのクーポンを取得する。
func (r *CouponRepository) FindByCode(ctx context.Context, code coupon.CouponCode) (*coupon.Coupon, error) {
	db := DBFromContext(ctx, r.db)

	var model CouponModel
	if err := db.WithContext(ctx).First(&model, "code = ?", code.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("coupon %s: %w", code.String(), shared.ErrNotFound)
		}
		return nil, err
	}

	return model.toDomain()
}

// Save はクーポンを永続化する（新規作成・更新の両方を担う upsert）。
//
// GORM の Save は主キーが既に存在すれば UPDATE、存在しなければ INSERT を
// 行う。クーポンの ID はアプリケーション側（ドメイン層）で採番しているため、
// GORM の自動採番機能に頼らずこの挙動をそのまま利用できる。
func (r *CouponRepository) Save(ctx context.Context, c *coupon.Coupon) error {
	db := DBFromContext(ctx, r.db)
	model := couponModelFromDomain(c)
	return db.WithContext(ctx).Save(&model).Error
}
