package coupon

import (
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// Coupon はクーポンを表す集約ルート（Aggregate Root）である。
//
// 割引方式（DiscountType）に応じて amount / ratePercent のどちらか一方
// だけが意味を持つ、という構造になっている。amount 型では amount
// フィールドのみが有効で ratePercent は使われず、rate 型ではその逆になる。
// このような「一方の値だけが有効」という不変条件を守るために、
// NewAmountCoupon / NewRateCoupon という 2 つの専用コンストラクタを用意し、
// 「amount 型なのに ratePercent が設定されている」といった不正な組み合わせを
// そもそも型として作れないようにしている（詳細は各コンストラクタのコメントを参照）。
type Coupon struct {
	id           CouponID
	code         CouponCode
	discountType DiscountType
	amount       shared.Money // amount 型のみ有効
	ratePercent  int          // rate 型のみ有効（1〜100）
	expiresAt    time.Time
	usageLimit   int
	usedCount    int
}

// NewAmountCoupon は「固定金額を割り引く」クーポンを新規発行するコンストラクタである。
// amount がゼロ円のクーポンは「割引が存在しない」ことと同義であり、
// クーポンとして発行する意味がないためドメインルール違反として拒否する。
func NewAmountCoupon(code CouponCode, amount shared.Money, expiresAt time.Time, usageLimit int) (*Coupon, error) {
	if amount.IsZero() {
		return nil, shared.NewDomainRuleError("coupon: amount discount must not be zero")
	}
	if err := validateUsageLimit(usageLimit); err != nil {
		return nil, err
	}
	return &Coupon{
		id:           GenerateCouponID(),
		code:         code,
		discountType: DiscountTypeAmount,
		amount:       amount,
		ratePercent:  0,
		expiresAt:    expiresAt,
		usageLimit:   usageLimit,
		usedCount:    0,
	}, nil
}

// NewRateCoupon は「合計金額に対する割合（%）で割り引く」クーポンを
// 新規発行するコンストラクタである。ratePercent は 1〜100 の範囲でなければ
// ならない（0% は割引が存在しないことと同義であり、100% を超える割引率は
// 業務的に意味を持たないため）。
func NewRateCoupon(code CouponCode, ratePercent int, expiresAt time.Time, usageLimit int) (*Coupon, error) {
	if ratePercent < 1 || ratePercent > 100 {
		return nil, shared.NewDomainRuleError("coupon: rate percent must be between 1 and 100, got %d", ratePercent)
	}
	if err := validateUsageLimit(usageLimit); err != nil {
		return nil, err
	}
	return &Coupon{
		id:           GenerateCouponID(),
		code:         code,
		discountType: DiscountTypeRate,
		amount:       shared.Money{}, // rate 型では未使用のためゼロ値のまま
		ratePercent:  ratePercent,
		expiresAt:    expiresAt,
		usageLimit:   usageLimit,
		usedCount:    0,
	}, nil
}

// validateUsageLimit は利用回数上限の業務ルールを検証する。
// 0 回以下の上限は「一度も使えないクーポン」を意味し、発行する意味がない。
func validateUsageLimit(usageLimit int) error {
	if usageLimit < 1 {
		return shared.NewDomainRuleError("coupon: usage limit must be at least 1, got %d", usageLimit)
	}
	return nil
}

// ReconstructCoupon は永続化層から読み込んだデータをもとに Coupon を再構築する。
//
// NewAmountCoupon / NewRateCoupon との違いは Product/Cart 等と同様である。
// 「これから新しく発行するクーポン」に対する業務ルール（金額・割合の
// 妥当性チェック）を課す入口が New* 系コンストラクタであるのに対し、
// DB から読み出す値は過去にそのチェックを通過済みのデータであるため、
// 再度同じ検証を強制しない「素通し」の再構築専用コンストラクタとして分離している。
// リポジトリ実装（infrastructure 層）からのみ呼ばれることを想定している。
func ReconstructCoupon(
	id CouponID,
	code CouponCode,
	discountType DiscountType,
	amount shared.Money,
	ratePercent int,
	expiresAt time.Time,
	usageLimit, usedCount int,
) *Coupon {
	return &Coupon{
		id:           id,
		code:         code,
		discountType: discountType,
		amount:       amount,
		ratePercent:  ratePercent,
		expiresAt:    expiresAt,
		usageLimit:   usageLimit,
		usedCount:    usedCount,
	}
}

// ID はクーポンの識別子を返す。
func (c *Coupon) ID() CouponID {
	return c.id
}

// Code はクーポンコードを返す。
func (c *Coupon) Code() CouponCode {
	return c.code
}

// Type は割引方式を返す。
func (c *Coupon) Type() DiscountType {
	return c.discountType
}

// Amount は固定金額割引の金額を返す（amount 型でのみ意味を持つ）。
func (c *Coupon) Amount() shared.Money {
	return c.amount
}

// RatePercent は割合割引のパーセンテージを返す（rate 型でのみ意味を持つ）。
func (c *Coupon) RatePercent() int {
	return c.ratePercent
}

// ExpiresAt は有効期限を返す。
func (c *Coupon) ExpiresAt() time.Time {
	return c.expiresAt
}

// UsageLimit は利用回数の上限を返す。
func (c *Coupon) UsageLimit() int {
	return c.usageLimit
}

// UsedCount は現在までの利用回数を返す。
func (c *Coupon) UsedCount() int {
	return c.usedCount
}

// IsExpired は now 時点でこのクーポンが期限切れかどうかを返す。
//
// ドメイン層で time.Now() を直接呼ばず、now を引数として受け取る設計に
// しているのは「時刻もまた入力の一つである」という考え方による。もし
// 内部で time.Now() を呼んでしまうと、期限切れの境界値（有効期限ちょうど、
// 1 秒前など）をテストする際にテスト実行時刻に依存してしまい、再現性の
// あるテストが書けなくなる。呼び出し側（ユースケース層）が now を注入する
// ことで、ドメイン層は完全に決定的（deterministic）に振る舞える。
func (c *Coupon) IsExpired(now time.Time) bool {
	return now.After(c.expiresAt)
}

// Use はクーポンを 1 回消費する。
//
// 期限切れの場合、または既に利用回数上限に達している場合はドメインルール
// 違反として拒否する。呼び出し元（注文確定ユースケース想定）は、
// このメソッドが成功した後にリポジトリ経由でクーポンの状態を保存する
// 責任を負う。
func (c *Coupon) Use(now time.Time) error {
	if c.IsExpired(now) {
		return shared.NewDomainRuleError("coupon: coupon %s is expired", c.code.String())
	}
	if c.usedCount >= c.usageLimit {
		return shared.NewDomainRuleError("coupon: coupon %s has reached its usage limit", c.code.String())
	}
	c.usedCount++
	return nil
}

// Refund はクーポンの消費を 1 回分取り消す（注文キャンセル時に呼ばれる
// ことを想定）。注文キャンセルで利用実績を戻す。在庫の Release と対になる
// 操作である。
//
// usedCount が 0 の状態で Refund を呼ぶことは「一度も消費していないのに
// 取り消そうとする」という想定外の状態であり、ドメインルール違反として
// 拒否する。
func (c *Coupon) Refund() error {
	if c.usedCount <= 0 {
		return shared.NewDomainRuleError("coupon: coupon %s has no usage to refund", c.code.String())
	}
	c.usedCount--
	return nil
}

// DiscountFor は注文合計金額 total に対する割引額を計算する。
//
// amount 型では「割引額は合計金額を超えない」という業務ルールを守るため、
// min(amount, total) を返す。合計を超える割引を認めてしまうと「お釣りが
// 出る」ことになり、業務的に不自然だからである。
// rate 型では total × ratePercent / 100 を計算する。int64 の整数除算に
// より 1 円未満の端数は切り捨てられる。四捨五入や切り上げなど他の丸め方式
// もあり得るが、本サンプルでは実装の単純さを優先し切り捨てを採用している。
func (c *Coupon) DiscountFor(total shared.Money) (shared.Money, error) {
	switch c.discountType {
	case DiscountTypeAmount:
		if c.amount.Currency() != total.Currency() {
			return shared.Money{}, shared.NewDomainRuleError(
				"coupon: currency mismatch: %s vs %s", c.amount.Currency(), total.Currency())
		}
		discount := c.amount.Amount()
		if discount > total.Amount() {
			// 割引額が合計金額を超える場合は合計金額でキャップする。
			discount = total.Amount()
		}
		return shared.NewMoney(discount, total.Currency())
	case DiscountTypeRate:
		discount := total.Amount() * int64(c.ratePercent) / 100
		return shared.NewMoney(discount, total.Currency())
	default:
		return shared.Money{}, shared.NewDomainRuleError("coupon: unknown discount type %q", c.discountType)
	}
}
