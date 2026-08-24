package shipping

import "context"

// GetShipmentUseCase は単一配送の詳細を取得するユースケースである。
// GetOrderUseCase / GetCartUseCase と同じく、トランザクション境界を
// 持たない薄いユースケースである（読み取り専用の操作は書き込みの
// 整合性を必要としないため）。
type GetShipmentUseCase struct {
	queryService ShipmentQueryService
}

// NewGetShipmentUseCase は GetShipmentUseCase を生成する。
func NewGetShipmentUseCase(queryService ShipmentQueryService) *GetShipmentUseCase {
	return &GetShipmentUseCase{queryService: queryService}
}

// Execute は id に対応する配送の詳細を返す。
func (uc *GetShipmentUseCase) Execute(ctx context.Context, id string) (*ShipmentDTO, error) {
	return uc.queryService.FindByID(ctx, id)
}
