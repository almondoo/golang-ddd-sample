package persistence

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB は PostgreSQL に接続する *gorm.DB を生成する。
//
// この関数はインフラストラクチャ層にのみ存在し、ドメイン層・アプリケーション層は
// GORM や PostgreSQL の存在を一切知らない。これが「依存性のルール」であり、
// DB のような技術的詳細は最も外側の層に閉じ込め、内側の層（ドメイン）は
// 純粋なビジネスロジックだけに集中できるようにする。
//
// dsn は "host=... user=... password=... dbname=... port=... sslmode=..." の
// ような Postgres 接続文字列を想定している。
func NewDB(dsn string) (*gorm.DB, error) {
	// TranslateError: true は、ドライバ固有のエラー（*pgconn.PgError 等）を
	// GORM 共通のセンチネルエラー（gorm.ErrDuplicatedKey 等）に変換させる
	// 設定である。StockRepository.Save は一意制約違反の検出に
	// errors.Is(err, gorm.ErrDuplicatedKey) を使っており、これが有効で
	// ないとドライバ固有エラーがそのまま返ってしまい判定できない。
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	return db, nil
}
