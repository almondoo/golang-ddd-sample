package order

import "context"

// Repository は Order 集約の永続化を抽象化するポート（インターフェース）である。
//
// cart.Repository / catalog.Repository と同様、ドメイン層がこの
// インターフェースを定義し、実装はインフラストラクチャ層（GORM を使った
// 具体的な SQL 操作）が提供する。依存の向きを反転させることで、
// ドメイン層は「どう保存されるか」を一切知らずに済む。
type Repository interface {
	// FindByID は id に対応する注文を取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id OrderID) (*Order, error)
	// Save は注文を永続化する（新規作成・更新の両方を担う upsert）。
	// 明細を含む集約全体を単位として保存する。
	Save(ctx context.Context, o *Order) error
}
