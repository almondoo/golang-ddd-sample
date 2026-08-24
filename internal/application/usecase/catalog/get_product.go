package catalog

import "context"

// GetProductUseCase は単一商品の詳細を取得するユースケースである。
type GetProductUseCase struct {
	queryService ProductQueryService
}

// NewGetProductUseCase は GetProductUseCase を生成する。
func NewGetProductUseCase(queryService ProductQueryService) *GetProductUseCase {
	return &GetProductUseCase{queryService: queryService}
}

// Execute は id に対応する商品の詳細を返す。
func (uc *GetProductUseCase) Execute(ctx context.Context, id string) (*ProductDTO, error) {
	return uc.queryService.FindByID(ctx, id)
}
