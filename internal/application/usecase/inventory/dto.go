// Package inventory は inventory コンテキストのユースケース（アプリケーションサービス）を
// まとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - set_stock.go … 書き込み系（コマンド）のユースケース。
//     domaininventory.Repository と tx.Manager に依存し、集約を読み込み・操作し・保存する。
//   - get_stock.go … 読み取り系（クエリ）のユースケース。
//     StockQueryService にのみ依存し、ドメイン層を経由せず DTO を直接返す。
//
// コマンドとクエリの区別はパッケージではなく、ファイル名と依存する型
// （Repository+tx.Manager か、QueryService のみか）で表現する（cart パッケージと同じ方針）。
package inventory

// StockDTO は在庫を画面表示・API 応答用に平坦化したデータ転送オブジェクトである。
//
// ドメインの Stock 集約とは別に DTO を用意しているのは、CQRS の考え方に
// 基づく。Available（引当可能残数）はドメイン層では計算プロパティだが、
// クエリ側では読み取り専用の値としてそのまま返す。
type StockDTO struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Reserved  int    `json:"reserved"`
	Available int    `json:"available"`
}
