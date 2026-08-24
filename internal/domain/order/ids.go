package order

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// OrderID は注文を一意に識別する識別子である。
//
// string を直接使わず専用の型を定義しているのは catalog.ProductID や
// cart.CustomerID と同じ理由である。型を分けることで、例えば
// CustomerID を誤って OrderID の引数に渡すような取り違えを
// コンパイル時に検出できるようになる。
type OrderID string

// NewOrderID は既存の文字列から OrderID を生成するコンストラクタである。
// 主に HTTP パスパラメータや DB から読み込んだ文字列を OrderID に
// 変換する際に使う。空文字列は不正な識別子としてドメインルール違反にする。
func NewOrderID(s string) (OrderID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("order: order id must not be empty")
	}
	return OrderID(s), nil
}

// GenerateOrderID は新しい注文を確定する際に使う ID 生成関数である。
// 採番ロジック（現状は UUID）を shared.NewID に委譲することで、
// 将来 ID 生成方式を変更する場合の変更点を shared パッケージに閉じ込める。
func GenerateOrderID() OrderID {
	return OrderID(shared.NewID())
}

// String は OrderID を文字列として取り出す。
// 主に永続化層やレスポンス生成時に生の文字列表現が必要な場面で使う。
func (id OrderID) String() string {
	return string(id)
}

// CustomerID は注文を行った顧客を識別する ID である。
//
// 注意: これは cart.CustomerID とは別の型である。order コンテキストは
// cart パッケージを一切 import しない。値としては同じ文字列（UUID）を
// 共有するが、境界づけられたコンテキスト（Bounded Context）同士は互いの
// 内部モデルに依存すべきではない、という DDD の原則（コンテキストの
// 自律性）に従い、あえて「ID という小さな型を重複させる」ことを選んでいる。
// この重複はコンパイル時の型チェックという利益に対して十分小さいコストである。
type CustomerID string

// NewCustomerID は既存の文字列から CustomerID を生成するコンストラクタである。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewCustomerID(s string) (CustomerID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("order: customer id must not be empty")
	}
	return CustomerID(s), nil
}

// String は CustomerID を文字列として取り出す。
func (id CustomerID) String() string {
	return string(id)
}
