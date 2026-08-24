package inventory

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ProductID は在庫が対象とする商品を識別する ID である。
//
// 注意（意図的な設計判断）: これは catalog.ProductID とは別の型であり、
// inventory コンテキストは catalog パッケージを一切 import しない。
// 値としては同じ文字列（UUID）を共有するが、境界づけられたコンテキスト
// （Bounded Context）は互いの内部モデルに依存すべきではない、という
// DDD の原則に従い、あえて「ID という小さな型を重複させる」ことを選んでいる
// （既存の cart/ids.go と同じ判断であり、コンテキストの自律性を保つための
// コストとしてはごく小さい）。
type ProductID string

// NewProductID は既存の文字列から ProductID を生成するコンストラクタである。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewProductID(s string) (ProductID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("inventory: product id must not be empty")
	}
	return ProductID(s), nil
}

// String は ProductID を文字列として取り出す。
func (id ProductID) String() string {
	return string(id)
}
