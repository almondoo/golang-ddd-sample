package query

import "context"

// GetOrderUseCase は単一注文の詳細を取得するユースケースである。
// GetProductUseCase / GetCartUseCase と同じく、トランザクション境界を
// 持たない薄いユースケースである（読み取り専用の操作は書き込みの
// 整合性を必要としないため）。
type GetOrderUseCase struct {
	queryService OrderQueryService
}

// NewGetOrderUseCase は GetOrderUseCase を生成する。
func NewGetOrderUseCase(queryService OrderQueryService) *GetOrderUseCase {
	return &GetOrderUseCase{queryService: queryService}
}

// Execute は id に対応する注文の詳細を返す。
func (uc *GetOrderUseCase) Execute(ctx context.Context, id string) (*OrderDTO, error) {
	return uc.queryService.FindByID(ctx, id)
}
