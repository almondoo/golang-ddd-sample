package shipping

import "context"

// ShipmentQueryService は配送の参照専用の問い合わせを表すポートである。
//
// domainshipping.Repository（ドメイン層）が「集約を読み書きする」ための
// インターフェースであるのに対し、こちらは「画面表示に必要な形へ加工済みの
// データを返す」ためのインターフェースである。order.OrderQueryService と
// 同じ軽量 CQRS の適用例であり、実装（infrastructure 層）は Shipment
// 集約を経由せず、GORM で shipments テーブルを直接読んで DTO を組み立てる。
type ShipmentQueryService interface {
	// FindByID は id に対応する配送を DTO として取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id string) (*ShipmentDTO, error)
}
