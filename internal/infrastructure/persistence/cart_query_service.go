package persistence

import (
	"context"

	"gorm.io/gorm"

	cartusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/cart"
)

// CartQuery は cartusecase.CartQueryService の GORM 実装である。
//
// cart.Repository を経由せず gorm.DB へ直接クエリを発行しているのは、
// カート集約が持たない「商品名」「価格」を画面表示のために取得したい
// ためである。cart_items（cart コンテキスト所有）と products（catalog
// コンテキスト所有）を SQL レベルで直接 JOIN する。これはドメイン層を
// 経由しない読み取り専用の参照であり、catalog のドメインロジックを
// 呼び出しているわけではないため、コンテキスト間の疎結合を破らない
// （あくまで表示専用のクエリ側の最適化である）。
type CartQuery struct {
	db *gorm.DB
}

// NewCartQuery は CartQuery を生成する。
func NewCartQuery(db *gorm.DB) *CartQuery {
	return &CartQuery{db: db}
}

// コンパイル時に CartQuery が cartusecase.CartQueryService を満たすことを保証する。
var _ cartusecase.CartQueryService = (*CartQuery)(nil)

// cartItemJoinRow は cart_items × products の JOIN 結果を受け取るための
// 内部専用の一時構造体である。DTO と似た形をしているが、
// SQL の SELECT 句と 1 対 1 対応させるためにあえて別の型として定義している。
type cartItemJoinRow struct {
	ProductID   string
	ProductName string
	PriceAmount int64
	Quantity    int
}

// FindByCustomerID は指定顧客のカートを、商品名・価格を結合した DTO として返す。
//
// カートが存在しない場合（cart_items 行が 0 件）や、明細はあるが金額計算の
// 結果 0 円になる場合であっても、エラーにはせず空の CartDTO を返す。
// これは cartusecase.CartQueryService のドキュメントコメントに記した通り、
// 参照側にとって「カートがない」と「カートは空」は区別する意味がないためである。
func (q *CartQuery) FindByCustomerID(ctx context.Context, customerID string) (*cartusecase.CartDTO, error) {
	db := DBFromContext(ctx, q.db)

	var rows []cartItemJoinRow
	err := db.WithContext(ctx).
		Table("cart_items").
		Select("cart_items.product_id AS product_id, products.name AS product_name, products.price_amount AS price_amount, cart_items.quantity AS quantity").
		Joins("JOIN products ON products.id = cart_items.product_id").
		Where("cart_items.customer_id = ?", customerID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]cartusecase.CartItemDTO, 0, len(rows))
	var total int64
	for _, r := range rows {
		subtotal := r.PriceAmount * int64(r.Quantity)
		items = append(items, cartusecase.CartItemDTO{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			PriceAmount: r.PriceAmount,
			Quantity:    r.Quantity,
			Subtotal:    subtotal,
		})
		total += subtotal
	}

	return &cartusecase.CartDTO{
		CustomerID:  customerID,
		Items:       items,
		TotalAmount: total,
	}, nil
}
