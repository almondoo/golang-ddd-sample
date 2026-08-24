package catalog

import "context"

// ListProductsUseCase は商品一覧を取得するユースケースである。
//
// トランザクション境界（tx.Manager）を持たない点が command 側のユースケースと
// 異なる。読み取り専用の操作はデータの整合性を書き換えないため、
// トランザクションで囲む必然性がない（DB のデフォルトの読み取り一貫性で十分）。
type ListProductsUseCase struct {
	queryService ProductQueryService
}

// NewListProductsUseCase は ListProductsUseCase を生成する。
func NewListProductsUseCase(queryService ProductQueryService) *ListProductsUseCase {
	return &ListProductsUseCase{queryService: queryService}
}

// Execute は商品一覧を返す。
func (uc *ListProductsUseCase) Execute(ctx context.Context) ([]ProductDTO, error) {
	return uc.queryService.List(ctx)
}
