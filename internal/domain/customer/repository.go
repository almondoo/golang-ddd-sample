package customer

import "context"

// Repository は Customer 集約の永続化を抽象化するポート（インターフェース）
// である。
//
// cart.Repository / order.Repository と同様、ドメイン層がこの
// インターフェースを定義し、実装はインフラストラクチャ層（GORM を使った
// 具体的な SQL 操作）が提供する。依存の向きを反転させることで、
// ドメイン層は「どう保存されるか」を一切知らずに済む。
type Repository interface {
	// FindByID は id に対応する顧客を、登録済みの住所（Address）すべてを
	// 含めた集約全体として取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id CustomerID) (*Customer, error)
	// Save は顧客を永続化する（新規作成・更新の両方を担う upsert）。
	// 住所を含む集約全体を単位として保存する。
	Save(ctx context.Context, c *Customer) error
}
