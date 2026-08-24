package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// OrderRepository は order.Repository の GORM 実装である。
//
// アプリケーション層は order.Repository（インターフェース）にのみ依存し、
// この構造体（具体的な GORM 実装）の存在を知らない。
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository は OrderRepository を生成する。
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// コンパイル時に OrderRepository が order.Repository を満たすことを保証する。
var _ order.Repository = (*OrderRepository)(nil)

// FindByID は id に対応する注文を orders + order_items から復元する。
func (r *OrderRepository) FindByID(ctx context.Context, id order.OrderID) (*order.Order, error) {
	db := DBFromContext(ctx, r.db)

	var model OrderModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	var itemModels []OrderItemModel
	if err := db.WithContext(ctx).Where("order_id = ?", id.String()).Find(&itemModels).Error; err != nil {
		return nil, err
	}

	return orderFromModels(model, itemModels)
}

// Save は注文の現在の状態を永続化する（新規作成・更新の両方を担う upsert）。
//
// 【集約全体を単位とした永続化】
// orders 行は GORM の Save（主キーが存在すれば UPDATE、なければ INSERT）で
// upsert する。order_items は cart_items と同じ「対象注文の既存明細を
// 全削除してから、現在の明細を全件挿入し直す」という delete-all-then-insert
// 方式を採る。Order 集約は「注文確定後に明細を差し替える」という業務操作を
// 持たない（明細は注文確定時点で確定し、以後は不変）ため、実際には
// Pay/Ship/Cancel のような状態遷移操作では明細に変化がなく、この
// delete-then-insert は冪等な no-op に近い動きになる。それでも実装を
// 単純に保つため、Save は常に「集約全体を丸ごと置き換える」という
// 一貫した意味論を採用している。
//
// この関数は tx.Manager の Do の中から呼ばれることを想定しており、
// DBFromContext がトランザクション用の *gorm.DB を返すため、orders への
// upsert と order_items の削除・挿入は 1 つのトランザクションとして
// アトミックに行われる。
func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
	db := DBFromContext(ctx, r.db)

	model := orderModelFromDomain(o)
	if err := db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	if err := db.WithContext(ctx).
		Where("order_id = ?", o.ID().String()).
		Delete(&OrderItemModel{}).Error; err != nil {
		return err
	}

	itemModels := orderItemModelsFromDomain(o)
	if len(itemModels) == 0 {
		// NewOrder は明細 0 件の注文を作らせないが、万一に備えて
		// cart_repository.go と同様に空の場合は挿入をスキップする。
		return nil
	}

	return db.WithContext(ctx).Create(&itemModels).Error
}
