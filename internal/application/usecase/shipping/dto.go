// Package shipping は shipping コンテキストのユースケース（アプリケーション
// サービス）をまとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - deliver_shipment.go … 書き込み系（コマンド）のユースケース。
//     domainshipping.Repository と tx.Manager に依存し、集約を
//     読み込み・操作し・保存する。
//   - get_shipment.go     … 読み取り系（クエリ）のユースケース。
//     ShipmentQueryService にのみ依存し、ドメイン層を経由せず DTO を
//     直接返す。
//
// order・cart コンテキストと同じく、コマンドとクエリの区別はパッケージでは
// なくファイル名と依存する型（Repository+tx.Manager か、QueryService のみ
// か）で表現する方針を踏襲している。
package shipping

// ShipmentDTO は配送を画面表示用に平坦化したデータ転送オブジェクトである。
//
// ドメインの Shipment 集約とは別に DTO を用意しているのは、
// order.OrderDTO / cart.CartDTO と同じ理由（CQRS の考え方）による。
type ShipmentDTO struct {
	ID      string `json:"id"`
	OrderID string `json:"orderId"`
	Address string `json:"address"`
	Status  string `json:"status"`
}
