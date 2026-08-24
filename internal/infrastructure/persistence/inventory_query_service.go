package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	inventoryusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/inventory"
)

// StockQuery は inventoryusecase.StockQueryService の GORM 実装である。
//
// 型名を StockQueryService としなかったのは、inventoryusecase パッケージの
// インターフェース名と実装の型名が同名になり紛らわしくなるのを避けるため
// である（catalog.ProductQuery と同じ命名判断）。
type StockQuery struct {
	db *gorm.DB
}

// NewStockQuery は StockQuery を生成する。
func NewStockQuery(db *gorm.DB) *StockQuery {
	return &StockQuery{db: db}
}

// コンパイル時に StockQuery が inventoryusecase.StockQueryService を満たすことを保証する。
var _ inventoryusecase.StockQueryService = (*StockQuery)(nil)

// FindByProductID は指定商品の在庫を DTO として返す。
//
// 対象商品の在庫レコードがまだ存在しない場合でもエラーにはせず、
// 指定した productID を持つゼロ値の StockDTO を返す。これは
// inventoryusecase.StockQueryService のドキュメントコメントに記した通り、
// 未登録は「在庫0」と同じ意味に倒すという方針に基づく。
func (q *StockQuery) FindByProductID(ctx context.Context, productID string) (*inventoryusecase.StockDTO, error) {
	db := DBFromContext(ctx, q.db)

	var model StockModel
	err := db.WithContext(ctx).First(&model, "product_id = ?", productID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &inventoryusecase.StockDTO{
				ProductID: productID,
				Quantity:  0,
				Reserved:  0,
				Available: 0,
			}, nil
		}
		return nil, err
	}

	return &inventoryusecase.StockDTO{
		ProductID: model.ProductID,
		Quantity:  model.Quantity,
		Reserved:  model.Reserved,
		Available: model.Quantity - model.Reserved,
	}, nil
}
