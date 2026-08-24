package query

import "context"

// GetCartUseCase はカートの参照を行うユースケースである。
//
// コマンド側のユースケース（AddItemUseCase 等）が cart.Repository と
// tx.Manager に依存するのに対し、こちらは CartQueryService 1 つにしか
// 依存しない。参照系はトランザクション境界を意識する必要がなく
// （単一の読み取りクエリで完結する）、書き込み系より薄いユースケースに
// なるのが自然である。
type GetCartUseCase struct {
	queryService CartQueryService
}

// NewGetCartUseCase は GetCartUseCase を生成する。
func NewGetCartUseCase(queryService CartQueryService) *GetCartUseCase {
	return &GetCartUseCase{queryService: queryService}
}

// Execute は指定顧客のカートを DTO として取得する。
func (u *GetCartUseCase) Execute(ctx context.Context, customerID string) (*CartDTO, error) {
	return u.queryService.FindByCustomerID(ctx, customerID)
}
