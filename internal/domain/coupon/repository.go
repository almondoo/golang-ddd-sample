package coupon

import "context"

// Repository はクーポン集約の永続化を抽象化するポート（インターフェース）である。
//
// ドメイン層がこのインターフェースを定義し、実装（GORM を使った具体的な
// SQL 操作等）はインフラストラクチャ層が提供する。依存の向きは
// インフラストラクチャ層 → ドメイン層 となり、依存性逆転の原則
// （Dependency Inversion Principle）を体現している。
//
// クーポンの検索キーに CouponID ではなく CouponCode を使っているのは、
// 「クーポンを利用する」という業務操作が常に顧客の入力したコードを起点に
// 行われるためである（顧客はクーポンの内部 ID を知らない）。
type Repository interface {
	// FindByCode は指定コードのクーポンを取得する。
	// 該当するクーポンが存在しない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByCode(ctx context.Context, code CouponCode) (*Coupon, error)
	// Save はクーポンを永続化する（新規作成・更新の両方を担う upsert）。
	Save(ctx context.Context, c *Coupon) error
}
