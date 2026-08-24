package coupon

import (
	"regexp"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// couponCodePattern はクーポンコードとして許可される文字列パターンである。
// 4〜20文字の範囲で、大文字英字（A-Z）・数字（0-9）・ハイフン（-）のみを許可する。
var couponCodePattern = regexp.MustCompile(`^[A-Z0-9-]{4,20}$`)

// CouponCode はクーポンコードを表す値オブジェクト（Value Object）である。
//
// 「コードは 4〜20 文字で、大文字英数字とハイフンのみから成る」という
// 形式ルールを、この値オブジェクトに閉じ込めている。もしこのルールを
// 単なる string のバリデーション関数として各所に散らばらせてしまうと、
// チェックを忘れた箇所から不正な形式のコードが紛れ込む余地が生まれる。
// CouponCode 型として一度生成できれば、それだけで「正しい形式のコード
// である」ことが型システムによって保証され、以降の呼び出し側は
// 再検証を意識しなくてよくなる。
type CouponCode struct {
	value string
}

// NewCouponCode は文字列から CouponCode を生成するコンストラクタである。
// 形式ルールに違反する文字列はドメインルール違反として拒否する。
func NewCouponCode(s string) (CouponCode, error) {
	if !couponCodePattern.MatchString(s) {
		return CouponCode{}, shared.NewDomainRuleError(
			"coupon: code must be 4-20 characters of uppercase letters, digits, or hyphens, got %q", s)
	}
	return CouponCode{value: s}, nil
}

// String はクーポンコードを文字列として取り出す。
func (c CouponCode) String() string {
	return c.value
}
