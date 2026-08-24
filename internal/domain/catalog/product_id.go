package catalog

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ProductID は商品を一意に識別する識別子（識別子オブジェクト）である。
//
// string を直接使わず専用の型を定義しているのは、「型システムでミスを防ぐ」
// という DDD の定石である。例えば CartID を誤って ProductID の引数に渡す
// ような取り違えは、コンパイル時にエラーとして検出できるようになる。
type ProductID string

// NewProductID は既存の文字列から ProductID を生成するコンストラクタである。
// 主に HTTP パスパラメータや DB から読み込んだ文字列を ProductID に
// 変換する際に使う。空文字列は不正な識別子としてドメインルール違反にする。
func NewProductID(s string) (ProductID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("catalog: product id must not be empty")
	}
	return ProductID(s), nil
}

// GenerateProductID は新しい商品を登録する際に使う ID 生成関数である。
// 採番ロジック（現状は UUID）を shared.NewID に委譲することで、
// 将来 ID 生成方式を変更する場合の変更点を shared パッケージに閉じ込める。
func GenerateProductID() ProductID {
	return ProductID(shared.NewID())
}

// String は ProductID を文字列として取り出す。
// 主に永続化層やレスポンス生成時に生の文字列表現が必要な場面で使う。
func (id ProductID) String() string {
	return string(id)
}
