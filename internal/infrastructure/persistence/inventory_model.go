package persistence

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
)

// StockModel は stocks テーブルに対応する GORM 用の構造体である。
//
// ドメインの Stock 集約とは意図的に別の型として定義している。GORM の
// タグはあくまで永続化の都合であり、ドメイン層の型にこれを持ち込むと
// 「ドメインが ORM の存在を知っている」ことになり、依存性のルールに
// 反する（catalog.ProductModel と同じ判断）。
type StockModel struct {
	ProductID string `gorm:"primaryKey"`
	Quantity  int
	Reserved  int
	// Version は楽観ロック用の版数カラムである。Save 時に
	// WHERE version = ? の条件付き更新に使い、更新のたびに +1 する。
	Version int `gorm:"not null;default:1"`
}

// TableName は GORM に対して物理テーブル名を明示する。
func (StockModel) TableName() string {
	return "stocks"
}

// toDomain は永続化モデルからドメイン集約を復元する。
// DB から読み出した値は過去にドメイン層のバリデーションを通過済みという
// 前提のもと、ReconstructStock（検証を行わない再構築コンストラクタ）を使う。
func (m StockModel) toDomain() (*inventory.Stock, error) {
	productID, err := inventory.NewProductID(m.ProductID)
	if err != nil {
		return nil, err
	}
	return inventory.ReconstructStock(productID, m.Quantity, m.Reserved, m.Version), nil
}

// stockModelFromDomain はドメイン集約から永続化モデルを組み立てる。
func stockModelFromDomain(s *inventory.Stock) StockModel {
	return StockModel{
		ProductID: s.ProductID().String(),
		Quantity:  s.Quantity(),
		Reserved:  s.Reserved(),
		Version:   s.Version(),
	}
}
