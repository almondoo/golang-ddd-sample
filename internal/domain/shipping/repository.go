package shipping

import "context"

// Repository は Shipment 集約の永続化を抽象化するポート（インターフェース）である。
//
// order.Repository / cart.Repository と同様、ドメイン層がこの
// インターフェースを定義し、実装はインフラストラクチャ層（GORM を使った
// 具体的な SQL 操作）が提供する。依存の向きを反転させることで、
// ドメイン層は「どう保存されるか」を一切知らずに済む。
type Repository interface {
	// FindByID は id に対応する配送を取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id ShipmentID) (*Shipment, error)
	// FindByOrderID は orderID に対応する配送を取得する。
	// 1 つの注文に対して配送は高々 1 件という前提に立つ。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByOrderID(ctx context.Context, orderID OrderID) (*Shipment, error)
	// Save は配送を永続化する（新規作成・更新の両方を担う upsert）。
	Save(ctx context.Context, s *Shipment) error
}
