package persistence

import (
	"context"

	"gorm.io/gorm"
)

// dbContextKey は context に *gorm.DB を格納する際に使う非公開のキー型である。
//
// 標準ライブラリの慣習として、context.WithValue のキーには他パッケージの
// キーと衝突しない専用の型を使うべきとされている（string 等の組み込み型を
// キーに使うと、異なるパッケージ同士でキーが衝突するリスクがある）。
// この型を非公開にすることで、persistence パッケージの外からは
// このキーで値を上書きできないことを保証している。
type dbContextKey struct{}

// WithDB は ctx に *gorm.DB を紐付けた新しい context を返す。
//
// 主に TxManager.Do がトランザクション用の *gorm.DB を context に埋め込む
// ために使う。これにより、ユースケースの中で呼ばれるリポジトリは
// context 経由でトランザクション中の DB ハンドルを受け取れる。
func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

// DBFromContext は ctx に紐付けられた *gorm.DB を取り出す。
//
// トランザクション中でなければ ctx には DB が埋め込まれていないため、
// その場合は fallback（通常はリポジトリが保持しているコネクションプール用の
// *gorm.DB）を返す。これにより、リポジトリの実装は「トランザクション中か
// どうか」を意識せず、常に DBFromContext(ctx, r.db) を呼ぶだけでよくなる。
func DBFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if db, ok := ctx.Value(dbContextKey{}).(*gorm.DB); ok && db != nil {
		return db
	}
	return fallback
}
