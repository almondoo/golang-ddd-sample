package coupon

import "github.com/almondoo/golang-ddd-sample/internal/domain/shared"

// DiscountType はクーポンの割引方式を表す値オブジェクトである。
//
// string を基底型にしているのは order.Status と同じ理由で、DB へそのまま
// 文字列として保存でき、ログ・デバッグ出力でも人間が読める形になるためである。
type DiscountType string

const (
	// DiscountTypeAmount は「固定金額を割り引く」方式を表す。
	DiscountTypeAmount DiscountType = "amount"
	// DiscountTypeRate は「合計金額に対する割合（%）で割り引く」方式を表す。
	DiscountTypeRate DiscountType = "rate"
)

// NewDiscountType は文字列から DiscountType を生成するコンストラクタである。
//
// 主に永続化層が DB に保存された文字列を Coupon 集約へ復元する際や、
// プレゼンテーション層からの入力を検証する際に使う。あらかじめ定義された
// 2 つの方式のいずれかであることを検証する。
func NewDiscountType(s string) (DiscountType, error) {
	switch DiscountType(s) {
	case DiscountTypeAmount, DiscountTypeRate:
		return DiscountType(s), nil
	default:
		return "", shared.NewDomainRuleError("coupon: unknown discount type %q", s)
	}
}

// String は DiscountType を文字列として取り出す。
func (t DiscountType) String() string {
	return string(t)
}
