package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
)

// TxManager は tx.Manager の GORM 実装である。
//
// アプリケーション層は tx.Manager インターフェースにのみ依存しており、
// この構造体（GORM 実装）の存在を知らない。依存関係は
// アプリケーション層 → tx.Manager（インターフェース）← TxManager（実装）
// という向きになっており、これは依存性逆転の原則（DIP）そのものである。
type TxManager struct {
	db *gorm.DB
}

// NewTxManager は TxManager を生成する。
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// コンパイル時に TxManager が tx.Manager を満たすことを保証するアサーション。
// 実行時ではなくビルド時にインターフェース実装漏れを検知できる。
var _ tx.Manager = (*TxManager)(nil)

// Do は fn を 1 つの DB トランザクションの中で実行する。
//
// GORM の db.Transaction はコールバックが nil を返せば COMMIT、
// エラーを返す（あるいは panic する）と ROLLBACK するという契約を持つ。
// ここではトランザクション用の *gorm.DB を WithDB で ctx に埋め込んでから
// fn を呼び出すことで、fn の中で呼ばれるリポジトリが
// persistence.DBFromContext(ctx, ...) 経由で自動的に同じトランザクションに
// 参加できるようにしている。
//
// これにより、ユースケース自身は「どのリポジトリ呼び出しがトランザクションに
// 含まれるか」を個別に配線する必要がなく、fn の中で呼んだものはすべて
// 自動的に同一トランザクションに含まれる、という単純なメンタルモデルで
// 実装できる。
func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.Transaction(func(txDB *gorm.DB) error {
		txCtx := WithDB(ctx, txDB)
		return fn(txCtx)
	})
}
