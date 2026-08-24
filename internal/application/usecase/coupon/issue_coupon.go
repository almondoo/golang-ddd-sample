package coupon

import (
	"context"
	"errors"
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincoupon "github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// IssueCouponInput はクーポン発行ユースケースの入力である。
//
// アプリケーション層は HTTP のリクエストボディをそのまま受け取らず、
// この専用の入力型に一度変換してから使う。Amount / RatePercent は
// どちらか一方だけが DiscountType に応じて意味を持つ（両方を常に受け取り、
// 未使用側は無視する）。
type IssueCouponInput struct {
	Code         string
	DiscountType string
	Amount       int64
	RatePercent  int
	ExpiresAt    time.Time
	UsageLimit   int
}

// IssueCouponOutput はクーポン発行ユースケースの出力である。
type IssueCouponOutput struct {
	CouponID string
}

// IssueCouponUseCase は新規クーポンを発行するユースケースである。
type IssueCouponUseCase struct {
	repo domaincoupon.Repository
	tx   tx.Manager
}

// NewIssueCouponUseCase は IssueCouponUseCase を生成する。
func NewIssueCouponUseCase(repo domaincoupon.Repository, txManager tx.Manager) *IssueCouponUseCase {
	return &IssueCouponUseCase{repo: repo, tx: txManager}
}

// Execute はクーポン発行ユースケースを実行する。
//
// 「1 ユースケース = 1 トランザクション」という原則に従い、ここで
// tx.Do を呼んでトランザクション境界を明示する。
func (u *IssueCouponUseCase) Execute(ctx context.Context, in IssueCouponInput) (IssueCouponOutput, error) {
	code, err := domaincoupon.NewCouponCode(in.Code)
	if err != nil {
		return IssueCouponOutput{}, err
	}

	discountType, err := domaincoupon.NewDiscountType(in.DiscountType)
	if err != nil {
		return IssueCouponOutput{}, err
	}

	// 割引方式ごとに異なる専用コンストラクタへ分岐する。「amount 型なのに
	// ratePercent が設定されている」といった不正な組み合わせをそもそも
	// 型として作れないようにするための、ドメイン層の設計をそのまま踏襲する。
	var c *domaincoupon.Coupon
	switch discountType {
	case domaincoupon.DiscountTypeAmount:
		amount, merr := shared.NewMoney(in.Amount, shared.JPY)
		if merr != nil {
			return IssueCouponOutput{}, merr
		}
		c, err = domaincoupon.NewAmountCoupon(code, amount, in.ExpiresAt, in.UsageLimit)
	case domaincoupon.DiscountTypeRate:
		c, err = domaincoupon.NewRateCoupon(code, in.RatePercent, in.ExpiresAt, in.UsageLimit)
	default:
		// NewDiscountType の検証をすでに通過しているため、通常はここへ到達しない。
		return IssueCouponOutput{}, shared.NewDomainRuleError("coupon: unknown discount type %q", discountType)
	}
	if err != nil {
		return IssueCouponOutput{}, err
	}

	err = u.tx.Do(ctx, func(ctx context.Context) error {
		// 一意性は DB 制約（uniqueIndex）でも守るが、分かりやすいエラーを
		// 返すために先に既存コードの有無を確認する。
		if _, ferr := u.repo.FindByCode(ctx, code); ferr == nil {
			return shared.NewDomainRuleError("coupon: coupon code %s already exists", code.String())
		} else if !errors.Is(ferr, shared.ErrNotFound) {
			return ferr
		}

		return u.repo.Save(ctx, c)
	})
	if err != nil {
		return IssueCouponOutput{}, err
	}

	return IssueCouponOutput{CouponID: c.ID().String()}, nil
}
