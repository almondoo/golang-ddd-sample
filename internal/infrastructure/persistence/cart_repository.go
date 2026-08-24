package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CartRepository は cart.Repository の GORM 実装である。
//
// アプリケーション層は cart.Repository（インターフェース）にのみ依存し、
// この構造体（具体的な GORM 実装）の存在を知らない。
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository は CartRepository を生成する。
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// コンパイル時に CartRepository が cart.Repository を満たすことを保証する。
var _ cart.Repository = (*CartRepository)(nil)

// FindByCustomerID は指定顧客のカートを cart_items から復元する。
//
// carts テーブル（ヘッダ行）が存在しないため、「カートが存在するか」は
// 「該当顧客の cart_items 行が 1 件以上あるか」で判定する。行が 0 件の
// 場合は shared.ErrNotFound を返す。呼び出し側（コマンドユースケース）は
// これを find-or-create のトリガーとして使う。
func (r *CartRepository) FindByCustomerID(ctx context.Context, id cart.CustomerID) (*cart.Cart, error) {
	db := DBFromContext(ctx, r.db)

	var models []CartItemModel
	if err := db.WithContext(ctx).Where("customer_id = ?", id.String()).Find(&models).Error; err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("cart for customer %s: %w", id, shared.ErrNotFound)
	}

	return cartFromModels(id, models)
}

// Save はカートの現在の状態を永続化する。
//
// 実装方針は「対象顧客の既存明細を全削除してから、現在の明細を全件挿入し
// 直す」という delete-all-then-insert である。集約の中身を差分検出して
// UPDATE/INSERT/DELETE を個別に振り分けるより実装がはるかに単純であり、
// 集約全体を 1 つの整合した単位として置き換えるという意味論とも一致する。
// カートの明細数には上限（最大 20 件）があるため、この方式でもパフォーマンス
// 上の問題は生じない規模である。より大きな集約や高頻度な更新が想定される
// 場合は差分更新への切り替えを検討すべきだが、本サンプルの学習目的の範囲では
// この単純さを優先する。
//
// この関数は tx.Manager の Do の中から呼ばれることを想定しており、
// DBFromContext がトランザクション用の *gorm.DB を返すため、削除と挿入は
// 1 つのトランザクションとしてアトミックに行われる。
func (r *CartRepository) Save(ctx context.Context, c *cart.Cart) error {
	db := DBFromContext(ctx, r.db)

	if err := db.WithContext(ctx).
		Where("customer_id = ?", c.CustomerID().String()).
		Delete(&CartItemModel{}).Error; err != nil {
		return err
	}

	models := cartItemModelsFromDomain(c)
	if len(models) == 0 {
		// 明細が 0 件（Clear 済み等）の場合、挿入するものがない。
		// 削除だけで「カートを空にする」という意図は達成できている。
		return nil
	}

	return db.WithContext(ctx).Create(&models).Error
}
