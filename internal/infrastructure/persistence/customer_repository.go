package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/customer"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CustomerRepository は customer.Repository の GORM 実装である。
//
// アプリケーション層は customer.Repository（インターフェース）にのみ依存し、
// この構造体（具体的な GORM 実装）の存在を知らない。
type CustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository は CustomerRepository を生成する。
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// コンパイル時に CustomerRepository が customer.Repository を満たすことを保証する。
var _ customer.Repository = (*CustomerRepository)(nil)

// FindByID は id に対応する顧客を customers + customer_addresses から復元する。
func (r *CustomerRepository) FindByID(ctx context.Context, id customer.CustomerID) (*customer.Customer, error) {
	db := DBFromContext(ctx, r.db)

	var model CustomerModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("customer %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	var addressModels []CustomerAddressModel
	if err := db.WithContext(ctx).Where("customer_id = ?", id.String()).Find(&addressModels).Error; err != nil {
		return nil, err
	}

	return customerFromModels(model, addressModels)
}

// Save は顧客の現在の状態を永続化する（新規作成・更新の両方を担う upsert）。
//
// 【集約全体を単位とした永続化】
// customers 行は GORM の Save（主キーが存在すれば UPDATE、なければ INSERT）で
// upsert する。customer_addresses は cart_items / order_items と同じ
// 「対象顧客の既存住所を全削除してから、現在の住所を全件挿入し直す」という
// delete-all-then-insert 方式を採る。Customer 集約は住所帳の更新頻度が
// 高くなく（最大 5 件という上限もある）、この単純な方式でもパフォーマンス上
// の問題は生じない規模である。
//
// この関数は tx.Manager の Do の中から呼ばれることを想定しており、
// DBFromContext がトランザクション用の *gorm.DB を返すため、customers への
// upsert と customer_addresses の削除・挿入は 1 つのトランザクションとして
// アトミックに行われる。この「集約全体を 1 トランザクションでまとめて
// 書き換える」という振る舞いこそが、Customer + Address という集約の境界を
// 永続化層で裏付けるものである。
func (r *CustomerRepository) Save(ctx context.Context, c *customer.Customer) error {
	db := DBFromContext(ctx, r.db)

	model := customerModelFromDomain(c)
	if err := db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	if err := db.WithContext(ctx).
		Where("customer_id = ?", c.ID().String()).
		Delete(&CustomerAddressModel{}).Error; err != nil {
		return err
	}

	addressModels := customerAddressModelsFromDomain(c)
	if len(addressModels) == 0 {
		// 住所が 0 件（RemoveAddress で最後の 1 件を削除した等）の場合、
		// 挿入するものがない。削除だけで「住所帳を空にする」という意図は
		// 達成できている。
		return nil
	}

	return db.WithContext(ctx).Create(&addressModels).Error
}
