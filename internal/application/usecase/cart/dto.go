// Package cart は cart コンテキストのユースケース（アプリケーションサービス）を
// まとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - add_item.go / remove_item.go … 書き込み系（コマンド）のユースケース。
//     cart.Repository と tx.Manager に依存し、集約を読み込み・操作し・保存する。
//   - get_cart.go                  … 読み取り系（クエリ）のユースケース。
//     CartQueryService にのみ依存し、ドメイン層を経由せず DTO を直接返す。
//
// かつては command / query という別パッケージに分割していたが、パッケージを
// 分けるほどの独立性はなく、コマンドとクエリの区別はパッケージではなく
// ファイル名と依存する型（Repository+tx.Manager か、QueryService のみか）で
// 表現する方針に統合した。
package cart

// CartDTO はカートを画面表示用に平坦化したデータ転送オブジェクトである。
//
// ドメインの Cart 集約とは別に DTO を用意しているのは、
// クエリ側（読み取り）とコマンド側（書き込み）で最適な形が異なるためである
// （CQRS の考え方）。ドメインモデルは価格を一切知らないが、画面には
// 商品名や金額を表示したい。この落差を埋めるのがクエリ側の責務であり、
// internal/infrastructure/persistence/cart_query_service.go の実装は
// cart_items テーブルと catalog コンテキストが所有する products テーブルを
// SQL レベルで直接 JOIN して求める。これはドメイン層を経由しないため、
// 書き込み側（Cart 集約が catalog.Product への参照を持たず ID のみで
// 疎結合を保つ設計）とは異なり、読み取り側には catalog の物理スキーマ
// （products のテーブル名・カラム名）への結合が残る。catalog 側が
// テーブル・カラムを変更しても、この JOIN はコンパイルエラーにはならず、
// 実行時に初めて壊れる（コンテキスト間の疎結合を破らない、とまでは
// 言えない点に注意。詳細は internal/domain/cart/README.md を参照）。
type CartDTO struct {
	CustomerID  string        `json:"customerId"`
	Items       []CartItemDTO `json:"items"`
	TotalAmount int64         `json:"totalAmount"`
}

// CartItemDTO はカート内の 1 明細を画面表示用に表したものである。
type CartItemDTO struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	PriceAmount int64  `json:"priceAmount"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}
