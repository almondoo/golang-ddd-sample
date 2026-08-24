package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	orderusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// OrderQuery は orderusecase.OrderQueryService の GORM 実装である。
//
// order.Repository を経由せず gorm.DB へ直接クエリを発行しているのは、
// catalog_query_service.go / cart_query_service.go に書いた通り、
// 読み取り専用の問い合わせにドメイン集約の組み立てコスト（Reconstruct*
// を通した値オブジェクトの再生成やバリデーション）を払わせないための
// 軽量 CQRS の実践である。orders・order_items は本来カートと違い
// products テーブルとの JOIN を必要としない点に注意する。OrderItem は
// 注文確定時点の商品名・単価をスナップショットとして自身の行に
// 持っているため、表示のために catalog を参照する必要が一切ない。
// これは order/README.md に記した「スナップショット設計」がクエリ側の
// 実装をも単純化するという副次的な利点を示す例でもある。
type OrderQuery struct {
	db *gorm.DB
}

// NewOrderQuery は OrderQuery を生成する。
func NewOrderQuery(db *gorm.DB) *OrderQuery {
	return &OrderQuery{db: db}
}

// コンパイル時に OrderQuery が orderusecase.OrderQueryService を満たすことを保証する。
var _ orderusecase.OrderQueryService = (*OrderQuery)(nil)

// FindByID は id に対応する注文を、明細を含む DTO として返す。
func (q *OrderQuery) FindByID(ctx context.Context, id string) (*orderusecase.OrderDTO, error) {
	db := DBFromContext(ctx, q.db)

	var model OrderModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	var itemModels []OrderItemModel
	if err := db.WithContext(ctx).Where("order_id = ?", id).Find(&itemModels).Error; err != nil {
		return nil, err
	}

	items := make([]orderusecase.OrderItemDTO, 0, len(itemModels))
	var total int64
	for _, im := range itemModels {
		subtotal := im.UnitPriceAmount * int64(im.Quantity)
		items = append(items, orderusecase.OrderItemDTO{
			ProductID:       im.ProductID,
			ProductName:     im.ProductName,
			UnitPriceAmount: im.UnitPriceAmount,
			Quantity:        im.Quantity,
			Subtotal:        subtotal,
		})
		total += subtotal
	}

	// 支払うべき金額 = 合計金額 - 割引額。クーポン未適用（CouponCode が
	// 空文字列）の場合は DiscountAmount が常に 0 なので、payable は
	// total とそのまま一致する。
	payable := total - model.DiscountAmount

	return &orderusecase.OrderDTO{
		ID:             model.ID,
		CustomerID:     model.CustomerID,
		Status:         model.Status,
		TotalAmount:    total,
		PlacedAt:       model.PlacedAt,
		Items:          items,
		CouponCode:     model.CouponCode,
		DiscountAmount: model.DiscountAmount,
		PayableAmount:  payable,
	}, nil
}
