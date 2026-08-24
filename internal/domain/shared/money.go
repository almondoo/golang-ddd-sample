package shared

import "fmt"

// Currency は通貨単位を表す値オブジェクトである。
// string を基底型にすることで JPY, USD のような通貨コードをそのまま表現できる。
type Currency string

// JPY は日本円を表す。
// 本サンプルでは通貨は日本円のみを扱う想定だが、将来の多通貨対応に備えて
// Money 側は Currency を汎用的なフィールドとして保持している。
const JPY Currency = "JPY"

// Money は金額を表す値オブジェクト（Value Object）である。
//
// DDD において「値オブジェクト」とは、識別子を持たず値そのもので同一性が
// 決まるオブジェクトを指す。Money は amount と currency の組み合わせで
// 同一性が決まり、一度生成した後は不変（immutable）である。
//
// フィールドを非公開（unexported）にしているのは、不変条件（金額は負数を
// 許さない、通貨単位は必ず設定されている等）を NewMoney コンストラクタと
// メソッド経由でのみ変更させ、外部から不正な状態を作れないようにするため。
// これは「自己防衛的（self-validating）な値オブジェクト」という DDD の
// 定石パターンである。
type Money struct {
	amount   int64
	currency Currency
}

// NewMoney は Money を生成するコンストラクタである。
// 負の金額はドメインルール違反として拒否する。
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, NewDomainRuleError("money: amount must not be negative, got %d", amount)
	}
	return Money{amount: amount, currency: currency}, nil
}

// Amount は金額（最小通貨単位、例: 円）を返す。
func (m Money) Amount() int64 {
	return m.amount
}

// Currency は通貨単位を返す。
func (m Money) Currency() Currency {
	return m.currency
}

// Add は 2 つの Money を加算した新しい Money を返す。
// 値オブジェクトは不変なので、既存の m を書き換えるのではなく
// 新しいインスタンスを返す点に注意する。
// 通貨単位が異なる Money 同士の加算は意味を持たないため、ドメインルール
// 違反として扱う（例: 100円 + 1ドル は定義できない）。
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, NewDomainRuleError("money: currency mismatch: %s vs %s", m.currency, other.currency)
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

// Multiply は Money を n 倍した新しい Money を返す。
// 例えばカート内の「単価 × 数量」の計算に使うことを想定している。
func (m Money) Multiply(n int) (Money, error) {
	if n < 0 {
		return Money{}, NewDomainRuleError("money: multiplier must not be negative, got %d", n)
	}
	return NewMoney(m.amount*int64(n), m.currency)
}

// IsZero は金額がゼロかどうかを返す。
func (m Money) IsZero() bool {
	return m.amount == 0
}

// Equals は値オブジェクトとしての等価性を判定する。
// 値オブジェクトの同一性は識別子ではなく保持する値そのもので決まるため、
// フィールドを直接比較する。
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String は Money を人間が読める形式で表示する（デバッグ・ログ用途）。
func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount, m.currency)
}
