package catalog

import "context"

// Repository は Product 集約の永続化を抽象化するポート（インターフェース）である。
//
// このインターフェースをドメイン層に置き、実装（GORM を使った具体的な
// SQL 発行など）をインフラストラクチャ層に置くのは、依存性逆転の原則
// （Dependency Inversion Principle）の適用である。ドメイン層・アプリケーション層は
// 「どう永続化されるか」を一切知らずに「永続化されている」という事実にだけ
// 依存できる。これにより、DB を PostgreSQL から別の実装に差し替える場合も
// ドメイン層のコードは一切変更が不要になる。
type Repository interface {
	// FindByID は id に対応する商品を取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByID(ctx context.Context, id ProductID) (*Product, error)
	// Save は商品を永続化する（新規作成・更新のいずれも担う upsert）。
	Save(ctx context.Context, p *Product) error
}
