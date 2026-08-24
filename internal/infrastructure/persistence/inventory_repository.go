package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// StockRepository は inventory.Repository の GORM 実装である。
//
// アプリケーション層は inventory.Repository（インターフェース）にのみ
// 依存しており、この構造体（具体的な GORM 実装）の存在を知らない。
type StockRepository struct {
	db *gorm.DB
}

// NewStockRepository は StockRepository を生成する。
func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

// コンパイル時に StockRepository が inventory.Repository を満たすことを保証する。
var _ inventory.Repository = (*StockRepository)(nil)

// FindByProductID は指定商品の在庫を取得する。
func (r *StockRepository) FindByProductID(ctx context.Context, id inventory.ProductID) (*inventory.Stock, error) {
	db := DBFromContext(ctx, r.db)

	var model StockModel
	if err := db.WithContext(ctx).First(&model, "product_id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock for product %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	return model.toDomain()
}

// Save は在庫を永続化する（新規作成・更新の両方を担う upsert）。
//
// Stock 集約は「商品ごとに 1 行」というシンプルな単一行の集約であるため、
// cart.Cart のような delete-all-then-insert 方式は不要で、GORM の Save
// （主キーが存在すれば UPDATE、存在しなければ INSERT）をそのまま使う
// （catalog.ProductRepository.Save と同じ方式）。
func (r *StockRepository) Save(ctx context.Context, s *inventory.Stock) error {
	db := DBFromContext(ctx, r.db)
	model := stockModelFromDomain(s)
	return db.WithContext(ctx).Save(&model).Error
}
