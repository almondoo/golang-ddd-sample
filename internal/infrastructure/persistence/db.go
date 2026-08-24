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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
