package customer

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CustomerID は顧客を一意に識別する識別子である。
//
// string を直接使わず専用の型を定義しているのは order.OrderID や
// cart.CustomerID と同じ理由である。型を分けることで、例えば
// AddressID を誤って CustomerID の引数に渡すような取り違えを
// コンパイル時に検出できるようになる。
type CustomerID string

// NewCustomerID は既存の文字列から CustomerID を生成するコンストラクタである。
// 主に HTTP パスパラメータや DB から読み込んだ文字列を CustomerID に
// 変換する際に使う。空文字列は不正な識別子としてドメインルール違反にする。
func NewCustomerID(s string) (CustomerID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("customer: customer id must not be empty")
	}
	return CustomerID(s), nil
}

// GenerateCustomerID は新しい顧客を登録する際に使う ID 生成関数である。
// 採番ロジックを shared.NewID に委譲することで、将来 ID 生成方式を
// 変更する場合の変更点を shared パッケージに閉じ込める。
func GenerateCustomerID() CustomerID {
	return CustomerID(shared.NewID())
}

// String は CustomerID を文字列として取り出す。
func (id CustomerID) String() string {
	return string(id)
}

// AddressID は Customer 集約に属する住所（子エンティティ）を一意に識別する
// 識別子である。
//
// CartItem（cart コンテキスト）や OrderItem（order コンテキスト）とは異なり
// Address は独自の識別子を持つ。これは「特定の住所を名指しでデフォルトに
// 変更する」「特定の住所だけを削除する」といった、子エンティティ単位での
// 操作を外部（アプリケーション層）から要求できる必要があるためである。
type AddressID string

// NewAddressID は既存の文字列から AddressID を生成するコンストラクタである。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewAddressID(s string) (AddressID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("customer: address id must not be empty")
	}
	return AddressID(s), nil
}

// GenerateAddressID は Customer.AddAddress が新しい住所を追加する際に使う
// ID 生成関数である。
func GenerateAddressID() AddressID {
	return AddressID(shared.NewID())
}

// String は AddressID を文字列として取り出す。
func (id AddressID) String() string {
	return string(id)
}
