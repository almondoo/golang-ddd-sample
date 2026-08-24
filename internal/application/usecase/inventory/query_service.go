package inventory

import "context"

// StockQueryService は在庫の参照専用の問い合わせを表すポートである。
//
// domaininventory.Repository（ドメイン層）が「集約を読み書きする」ための
// インターフェースであるのに対し、こちらは「画面表示・API 応答に必要な
// 形へ加工済みのデータを返す」ためのインターフェースである。両者を
// 分けているのは CQRS（コマンドとクエリの責務分離）の考え方に基づく
// （cart.CartQueryService と同じ設計）。
type StockQueryService interface {
	// FindByProductID は指定商品の在庫を DTO として取得する。
	//
	// 対象商品の在庫レコードがまだ 1 度も作られていない場合でも、エラーには
	// せず、指定した productID を持つゼロ値の StockDTO（quantity/reserved/available
	// がすべて 0）を返す。これは「在庫がまだ登録されていない」ことと
	// 「在庫は登録されているが数量が 0」であることを、参照側（利用者）から
	// 見た意味としては区別する価値がないためである（未登録は「在庫0」と
	// 同じ意味に倒す）。これを区別してエラーとして扱うと、呼び出し側が
	// 毎回「404 かもしれないので特別扱いする」という不要な分岐を
	// 持たされることになる（cart.CartQueryService と同じ判断）。
	FindByProductID(ctx context.Context, productID string) (*StockDTO, error)
}
