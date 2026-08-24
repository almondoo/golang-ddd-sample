package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ProductRepository は catalog.Repository の GORM 実装である。
//
// アプリケーション層は catalog.Repository（インターフェース）にのみ
// 依存しており、この構造体を直接参照しない。依存関係の向きは
// アプリケーション層 → catalog.Repository（インターフェース）← ProductRepository（実装）
// となっており、依存性逆転の原則（DIP）を体現している。
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository は ProductRepository を生成する。
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// コンパイル時に ProductRepository が catalog.Repository を満たすことを保証する。
var _ catalog.Repository = (*ProductRepository)(nil)

// FindByID は id に対応する商品を取得する。
func (r *ProductRepository) FindByID(ctx context.Context, id catalog.ProductID) (*catalog.Product, error) {
	db := DBFromContext(ctx, r.db)

	var model ProductModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	return model.toDomain()
}

// Save は商品を永続化する（新規作成・更新の両方を担う upsert）。
//
// GORM の Save は主キーが既に存在すれば UPDATE、存在しなければ INSERT を
// 行う。ここでは商品の ID をアプリケーション側（ドメイン層）で採番して
// いるため、GORM の自動採番機能に頼らずこの挙動をそのまま利用できる。
func (r *ProductRepository) Save(ctx context.Context, p *catalog.Product) error {
	db := DBFromContext(ctx, r.db)
	model := productModelFromDomain(p)
	return db.WithContext(ctx).Save(&model).Error
}
