package persistence

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/cart"
)

// CartItemModel は cart_items テーブルに対応する GORM 用の構造体である。
//
// 注意: carts テーブルは存在しない。カート集約は「顧客 ID + 明細の集まり」
// でしかなく、明細を持たない空のカートはドメイン上「まだ何も持っていない
// カート」であって、ヘッダ行として永続化しておく価値のある独自の状態を
// 持たない（作成日時等のメタデータも今のところ不要）。そのため、
// カートという集約全体を「1 顧客につき 0〜N 件の cart_items 行」という
// 形だけで表現し、ヘッダ行（carts テーブル）を持たずに済ませている。
// これは学習用の割り切りであり、将来「カート作成日時を記録したい」等の
// 要件が生まれた場合は carts テーブルを追加する設計変更が必要になる。
//
// 複合主キー（CustomerID, ProductID）を使っているのは、
// 「同一顧客・同一商品の明細は 1 行しか存在しない」という不変条件（マージ）
// をアプリケーション側だけでなく DB スキーマのレベルでも保証するためである。
type CartItemModel struct {
	CustomerID string `gorm:"primaryKey"`
	ProductID  string `gorm:"primaryKey"`
	Quantity   int
}

// TableName は GORM に対して物理テーブル名を明示する。
func (CartItemModel) TableName() string {
	return "cart_items"
}

// cartFromModels は cart_items テーブルの行群からカート集約を復元する。
// DB に保存されている値はすでに過去にドメイン層の検証を通過済みという
// 前提のもと、ReconstructCart / ReconstructCartItem（検証を行わない
// 再構築コンストラクタ）を使う。
func cartFromModels(customerID cart.CustomerID, models []CartItemModel) (*cart.Cart, error) {
	items := make([]cart.CartItem, 0, len(models))
	for _, m := range models {
		productID, err := cart.NewProductID(m.ProductID)
		if err != nil {
			return nil, err
		}
		items = append(items, cart.ReconstructCartItem(productID, m.Quantity))
	}
	return cart.ReconstructCart(customerID, items), nil
}

// cartItemModelsFromDomain はカート集約から永続化モデルの一覧を組み立てる。
func cartItemModelsFromDomain(c *cart.Cart) []CartItemModel {
	items := c.Items()
	models := make([]CartItemModel, 0, len(items))
	for _, item := range items {
		models = append(models, CartItemModel{
			CustomerID: c.CustomerID().String(),
			ProductID:  item.ProductID().String(),
			Quantity:   item.Quantity(),
		})
	}
	return models
}
