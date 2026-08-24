package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shipping"
)

// ShipmentRepository は shipping.Repository の GORM 実装である。
//
// アプリケーション層は shipping.Repository（インターフェース）にのみ
// 依存し、この構造体（具体的な GORM 実装）の存在を知らない
// （catalog.Repository / order.Repository と同じ依存性逆転の原則）。
type ShipmentRepository struct {
	db *gorm.DB
}

// NewShipmentRepository は ShipmentRepository を生成する。
func NewShipmentRepository(db *gorm.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

// コンパイル時に ShipmentRepository が shipping.Repository を満たすことを保証する。
var _ shipping.Repository = (*ShipmentRepository)(nil)

// FindByID は id に対応する配送を取得する。
func (r *ShipmentRepository) FindByID(ctx context.Context, id shipping.ShipmentID) (*shipping.Shipment, error) {
	db := DBFromContext(ctx, r.db)

	var model ShipmentModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("shipment %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	return model.toDomain()
}

// FindByOrderID は orderID に対応する配送を取得する。
func (r *ShipmentRepository) FindByOrderID(ctx context.Context, orderID shipping.OrderID) (*shipping.Shipment, error) {
	db := DBFromContext(ctx, r.db)

	var model ShipmentModel
	if err := db.WithContext(ctx).First(&model, "order_id = ?", orderID.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("shipment for order %s: %w", orderID, shared.ErrNotFound)
		}
		return nil, err
	}

	return model.toDomain()
}

// Save は配送を永続化する（新規作成・更新の両方を担う upsert）。
//
// GORM の Save は主キーが既に存在すれば UPDATE、存在しなければ INSERT を
// 行う。配送の ID はアプリケーション側（ドメイン層）で採番しているため、
// GORM の自動採番機能に頼らずこの挙動をそのまま利用できる
// （catalog_repository.go の ProductRepository.Save と同じ方針）。
func (r *ShipmentRepository) Save(ctx context.Context, s *shipping.Shipment) error {
	db := DBFromContext(ctx, r.db)
	model := shipmentModelFromDomain(s)
	return db.WithContext(ctx).Save(&model).Error
}
