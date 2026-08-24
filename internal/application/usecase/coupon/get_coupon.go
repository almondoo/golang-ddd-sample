package coupon

import "context"

// GetCouponUseCase はクーポンの参照を行うユースケースである。
//
// コマンド側のユースケース（IssueCouponUseCase）が domaincoupon.Repository と
// tx.Manager に依存するのに対し、こちらは CouponQueryService 1 つにしか
// 依存しない。参照系はトランザクション境界を意識する必要がなく
// （単一の読み取りクエリで完結する）、書き込み系より薄いユースケースに
// なるのが自然である。
type GetCouponUseCase struct {
	queryService CouponQueryService
}

// NewGetCouponUseCase は GetCouponUseCase を生成する。
func NewGetCouponUseCase(queryService CouponQueryService) *GetCouponUseCase {
	return &GetCouponUseCase{queryService: queryService}
}

// Execute は指定コードのクーポンを DTO として取得する。
func (u *GetCouponUseCase) Execute(ctx context.Context, code string) (*CouponDTO, error) {
	return u.queryService.FindByCode(ctx, code)
}
