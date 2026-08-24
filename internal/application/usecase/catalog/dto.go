// Package catalog は catalog コンテキストのユースケース（アプリケーションサービス）を
// まとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - register_product.go / change_price.go … 書き込み系（コマンド）のユースケース。
//     catalog.Repository と tx.Manager に依存し、集約を読み込み・操作し・保存する。
//   - list_products.go / get_product.go     … 読み取り系（クエリ）のユースケース。
//     ProductQueryService にのみ依存し、ドメイン層を経由せず DTO を直接返す。
//
// かつては command / query という別パッケージに分割していたが、パッケージを
// 分けるほどの独立性はなく、むしろ「catalog コンテキストの入出力窓口」として
// 1 つのパッケージにまとまっている方が見通しやすいと判断し統合した。コマンドと
// クエリの区別はパッケージではなく、ファイル名と依存する型（Repository+tx.Manager
// か、QueryService のみか）で表現する。
package catalog

// ProductDTO は問い合わせ（Query）結果としてクライアントに返す
// データ転送オブジェクト（Data Transfer Object）である。
//
// ドメインの Product 集約をそのまま JSON にシリアライズせず、専用の DTO を
// 経由させているのは、「表示に必要な形」と「業務ルールを守るための内部
// 構造」を分離するためである。DTO はドメインの不変条件を一切持たない
// ただのデータの入れ物であり、Product の内部フィールド構成を変更しても
// API の応答形式を独立して保てる。
type ProductDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceAmount int64  `json:"priceAmount"`
	Currency    string `json:"currency"`
}
