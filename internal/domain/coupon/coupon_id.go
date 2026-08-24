package coupon

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CouponID はクーポンを一意に識別する識別子である。
//
// string を直接使わず専用の型を定義しているのは catalog.ProductID や
// order.OrderID と同じ理由である。型を分けることで、例えば OrderID を
// 誤って CouponID の引数に渡すような取り違えをコンパイル時に検出できる。
type CouponID string

// NewCouponID は既存の文字列から CouponID を生成するコンストラクタである。
// 主に永続化層から読み込んだ文字列を CouponID に変換する際に使う。
// 空文字列は不正な識別子としてドメインルール違反にする。
func NewCouponID(s string) (CouponID, error) {
	if s == "" {
		return "", shared.NewDomainRuleError("coupon: coupon id must not be empty")
	}
	return CouponID(s), nil
}

// GenerateCouponID は新しいクーポンを発行する際に使う ID 生成関数である。
// 採番ロジック（現状は UUID）を shared.NewID に委譲することで、
// 将来 ID 生成方式を変更する場合の変更点を shared パッケージに閉じ込める。
func GenerateCouponID() CouponID {
	return CouponID(shared.NewID())
}

// String は CouponID を文字列として取り出す。
// 主に永続化層やレスポンス生成時に生の文字列表現が必要な場面で使う。
func (id CouponID) String() string {
	return string(id)
}
