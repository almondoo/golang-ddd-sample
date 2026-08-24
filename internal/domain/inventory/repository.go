package inventory

import "context"

// Repository は Stock 集約の永続化を抽象化するポート（インターフェース）である。
//
// ドメイン層がこのインターフェースを定義し、実装（GORM を使った具体的な
// SQL 操作等）はインフラストラクチャ層が提供する。依存の向きは
// インフラストラクチャ層 → ドメイン層 となり、依存性逆転の原則
// （Dependency Inversion Principle）を体現する（catalog.Repository / cart.Repository
// と同じ設計）。
type Repository interface {
	// FindByProductID は指定商品の在庫を取得する。
	// 該当する在庫が存在しない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByProductID(ctx context.Context, id ProductID) (*Stock, error)
	// Save は在庫の現在の状態を永続化する（新規作成・更新の両方を担う upsert）。
	Save(ctx context.Context, s *Stock) error
}
