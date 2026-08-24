package catalog

import "context"

// ProductQueryService は商品情報の問い合わせ（読み取り専用の参照）を
// 抽象化するポートである。
//
// これは軽量な CQRS（コマンドクエリ責務分離）の適用例である。書き込み側
// （register_product.go / change_price.go）は domaincatalog.Repository 経由で
// ドメイン集約を読み書きし、業務ルールを守りながら状態を変更する。一方、
// 読み取り側はドメイン集約を経由せず、DB から必要な列だけを直接 DTO に
// 詰めて返す。
//
// なぜ読み取りにまでドメイン層を通さないのか？
// 一覧表示のような「ただ表示するだけ」の処理でドメイン集約を都度
// 組み立てるのはオーバーヘッドであり、集約が持つ業務ルールも読み取り
// には不要である。読み取り専用の関心事（ページング、ソート、JOIN 等）を
// ドメイン層に持ち込まずに済むという実利もある。
type ProductQueryService interface {
	// List は登録済みの商品を一覧で返す。
	List(ctx context.Context) ([]ProductDTO, error)
	// FindByID は id に対応する商品を返す。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id string) (*ProductDTO, error)
}
