package inventory

import "context"

// GetStockUseCase は在庫の参照を行うユースケースである。
//
// コマンド側のユースケース（SetStockUseCase）が domaininventory.Repository と
// tx.Manager に依存するのに対し、こちらは StockQueryService 1 つにしか
// 依存しない。参照系はトランザクション境界を意識する必要がなく
// （単一の読み取りクエリで完結する）、書き込み系より薄いユースケースに
// なるのが自然である（cart.GetCartUseCase と同じ設計）。
type GetStockUseCase struct {
	queryService StockQueryService
}

// NewGetStockUseCase は GetStockUseCase を生成する。
func NewGetStockUseCase(queryService StockQueryService) *GetStockUseCase {
	return &GetStockUseCase{queryService: queryService}
}

// Execute は指定商品の在庫を DTO として取得する。
func (u *GetStockUseCase) Execute(ctx context.Context, productID string) (*StockDTO, error) {
	return u.queryService.FindByProductID(ctx, productID)
}
