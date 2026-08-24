package customer

import "context"

// CustomerQueryService は顧客の参照専用の問い合わせを表すポートである。
//
// domaincustomer.Repository（ドメイン層）が「集約を読み書きする」ための
// インターフェースであるのに対し、こちらは「画面表示に必要な形へ加工済みの
// データを返す」ためのインターフェースである。両者を分けているのは
// cart.CartQueryService / order.OrderQueryService と同じ CQRS
// （コマンドとクエリの責務分離）の考え方に基づく。
type CustomerQueryService interface {
	// FindByID は id に対応する顧客を、登録済みの住所を含む DTO として取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id string) (*CustomerDTO, error)
}
