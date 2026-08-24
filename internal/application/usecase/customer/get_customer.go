package customer

import "context"

// GetCustomerUseCase は顧客の参照を行うユースケースである。
//
// コマンド側のユースケース（AddAddressUseCase 等）が
// domaincustomer.Repository と tx.Manager に依存するのに対し、こちらは
// CustomerQueryService 1 つにしか依存しない。参照系はトランザクション
// 境界を意識する必要がなく（単一の読み取りクエリで完結する）、
// 書き込み系より薄いユースケースになるのが自然である
// （cart.GetCartUseCase / order.GetOrderUseCase と同じ設計）。
type GetCustomerUseCase struct {
	queryService CustomerQueryService
}

// NewGetCustomerUseCase は GetCustomerUseCase を生成する。
func NewGetCustomerUseCase(queryService CustomerQueryService) *GetCustomerUseCase {
	return &GetCustomerUseCase{queryService: queryService}
}

// Execute は指定顧客を、登録済みの住所を含む DTO として取得する。
func (u *GetCustomerUseCase) Execute(ctx context.Context, id string) (*CustomerDTO, error) {
	return u.queryService.FindByID(ctx, id)
}
