package cart

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CustomerID は顧客を一意に識別する識別子である。
//
// カートコンテキストは顧客そのもの（氏名や住所等）を管理しないが、
// 「誰のカートか」を表現するために顧客 ID だけを軽量に保持する。
// これも catalog.ProductID と同様、型を分けることで取り違えを
// コンパイル時に検出できるようにする定石である。
type CustomerID string

// NewCustomerID は既存の文字列から CustomerID を生成するコンストラクタである。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewCustomerID(s string) (CustomerID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("cart: customer id must not be empty")
	}
	return CustomerID(s), nil
}

// String は CustomerID を文字列として取り出す。
func (id CustomerID) String() string {
	return string(id)
}

// ProductID はカートが参照する商品を識別する ID である。
//
// 注意: これは catalog.ProductID とは別の型である。カートコンテキストは
// catalog パッケージを一切 import しない。値としては同じ文字列（UUID）を
// 共有するが、境界づけられたコンテキスト（Bounded Context）同士は互いの
// 内部モデルに依存すべきではない、という DDD の原則に従い、あえて
// 「ID という小さな型を重複させる」ことを選んでいる。コンテキスト間の
// 結合を避けるためのコストとしては、この程度の重複はごく小さい。
type ProductID string

// NewProductID は既存の文字列から ProductID を生成するコンストラクタである。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewProductID(s string) (ProductID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("cart: product id must not be empty")
	}
	return ProductID(s), nil
}

// String は ProductID を文字列として取り出す。
func (id ProductID) String() string {
	return string(id)
}
