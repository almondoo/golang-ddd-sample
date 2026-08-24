package persistence

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ProductModel は products テーブルに対応する GORM 用の構造体である。
//
// ドメインの Product 集約とは意図的に別の型として定義している。
// GORM のタグ（gorm:"primaryKey" 等）はあくまで永続化の都合であり、
// ドメイン層の型にこれを持ち込むと「ドメインが ORM の存在を知っている」
// ことになり、依存性のルール（ドメインは技術的詳細に依存しない）に反する。
// この 2 つの型の間を変換する責務は、このファイルの toDomain / productModelFromDomain が担う。
//
// なお、created_at / updated_at のような監査列はあえて持たない。
// このサンプルの永続化は「ドメインモデルから GORM モデルを組み立て直して
// Save する」方式のため、created_at のような書き込み時に保持できない
// 監査列は持たない。GORM の Save は全カラム UPDATE を行うので、値を
// 復元しない列はゼロ値で上書きされてしまう（実運用では fetch-then-update や
// Omit、マイグレーション管理と併せて設計する）。
type ProductModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	PriceAmount int64
	Currency    string
}

// TableName は GORM に対して物理テーブル名を明示する。
// GORM のデフォルト命名規則（構造体名の複数形スネークケース）に頼らず
// 明示することで、将来構造体名を変更してもテーブル名は変わらないようにする。
func (ProductModel) TableName() string {
	return "products"
}

// toDomain は永続化モデルからドメイン集約を復元する。
// DB から読み出した値は過去にドメイン層のバリデーションを通過済みという
// 前提のもと、ReconstructProduct（検証を行わない再構築コンストラクタ）を使う。
func (m ProductModel) toDomain() (*catalog.Product, error) {
	id, err := catalog.NewProductID(m.ID)
	if err != nil {
		return nil, err
	}
	price, err := shared.NewMoney(m.PriceAmount, shared.Currency(m.Currency))
	if err != nil {
		return nil, err
	}
	return catalog.ReconstructProduct(id, m.Name, m.Description, price), nil
}

// productModelFromDomain はドメイン集約から永続化モデルを組み立てる。
func productModelFromDomain(p *catalog.Product) ProductModel {
	return ProductModel{
		ID:          p.ID().String(),
		Name:        p.Name(),
		Description: p.Description(),
		PriceAmount: p.Price().Amount(),
		Currency:    string(p.Price().Currency()),
	}
}
