package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	shippingusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/shipping"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ShipmentQuery は shippingusecase.ShipmentQueryService の GORM 実装である。
//
// shipping.Repository を経由せず gorm.DB へ直接クエリを発行しているのは、
// order_query_service.go / catalog_query_service.go と同じ理由（読み取り
// 専用の問い合わせにドメイン集約の組み立てコストを払わせないための
// 軽量 CQRS の実践）による。
type ShipmentQuery struct {
	db *gorm.DB
}

// NewShipmentQuery は ShipmentQuery を生成する。
func NewShipmentQuery(db *gorm.DB) *ShipmentQuery {
	return &ShipmentQuery{db: db}
}

// コンパイル時に ShipmentQuery が shippingusecase.ShipmentQueryService を満たすことを保証する。
var _ shippingusecase.ShipmentQueryService = (*ShipmentQuery)(nil)

// FindByID は id に対応する配送を DTO として返す。
func (q *ShipmentQuery) FindByID(ctx context.Context, id string) (*shippingusecase.ShipmentDTO, error) {
	db := DBFromContext(ctx, q.db)

	var model ShipmentModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("shipment %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	return &shippingusecase.ShipmentDTO{
		ID:      model.ID,
		OrderID: model.OrderID,
		Address: model.Address,
		Status:  model.Status,
	}, nil
}
