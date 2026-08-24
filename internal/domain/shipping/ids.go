package shipping

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ShipmentID は配送を一意に識別する識別子である。
//
// string を直接使わず専用の型を定義しているのは order.OrderID や
// catalog.ProductID と同じ理由である。型を分けることで、例えば
// OrderID を誤って ShipmentID の引数に渡すような取り違えを
// コンパイル時に検出できるようになる。
type ShipmentID string

// NewShipmentID は既存の文字列から ShipmentID を生成するコンストラクタである。
// 主に HTTP パスパラメータや DB から読み込んだ文字列を ShipmentID に
// 変換する際に使う。空文字列は不正な識別子としてドメインルール違反にする。
func NewShipmentID(s string) (ShipmentID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("shipping: shipment id must not be empty")
	}
	return ShipmentID(s), nil
}

// GenerateShipmentID は新しい配送を生成する際に使う ID 生成関数である。
// 採番ロジック（現状は UUID）を shared.NewID に委譲することで、
// 将来 ID 生成方式を変更する場合の変更点を shared パッケージに閉じ込める。
func GenerateShipmentID() ShipmentID {
	return ShipmentID(shared.NewID())
}

// String は ShipmentID を文字列として取り出す。
// 主に永続化層やレスポンス生成時に生の文字列表現が必要な場面で使う。
func (id ShipmentID) String() string {
	return string(id)
}

// OrderID は配送対象の注文を識別する ID である。
//
// 注意: これは order.OrderID とは別の型である。shipping コンテキストは
// order パッケージを一切 import しない。値としては同じ文字列（UUID）を
// 共有するが、境界づけられたコンテキスト（Bounded Context）同士は互いの
// 内部モデルに依存すべきではない、という DDD の原則（コンテキストの
// 自律性）に従い、あえて「ID という小さな型を重複させる」ことを選んでいる
// （order.CustomerID が cart.CustomerID と別型である理由と同じ）。
// この重複はコンパイル時の型チェックという利益に対して十分小さいコストである。
type OrderID string

// NewOrderID は既存の文字列から OrderID を生成するコンストラクタである。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewOrderID(s string) (OrderID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("shipping: order id must not be empty")
	}
	return OrderID(s), nil
}

// String は OrderID を文字列として取り出す。
func (id OrderID) String() string {
	return string(id)
}
