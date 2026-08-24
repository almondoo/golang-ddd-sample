package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	customerusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/customer"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CustomerQuery は customerusecase.CustomerQueryService の GORM 実装である。
//
// customer.Repository を経由せず gorm.DB へ直接クエリを発行しているのは、
// cart_query_service.go / order_query_service.go と同じ CQRS の実践である。
// customers・customer_addresses は他コンテキストのテーブルとの JOIN を
// 必要としない点は order_query_service.go（スナップショット設計）に近いが、
// こちらは「集約内の子エンティティ（住所）テーブルとの JOIN」のみで
// DTO を組み立てられる、というのがこのクエリの特徴である。
type CustomerQuery struct {
	db *gorm.DB
}

// NewCustomerQuery は CustomerQuery を生成する。
func NewCustomerQuery(db *gorm.DB) *CustomerQuery {
	return &CustomerQuery{db: db}
}

// コンパイル時に CustomerQuery が customerusecase.CustomerQueryService を満たすことを保証する。
var _ customerusecase.CustomerQueryService = (*CustomerQuery)(nil)

// FindByID は id に対応する顧客を、住所を含む DTO として返す。
func (q *CustomerQuery) FindByID(ctx context.Context, id string) (*customerusecase.CustomerDTO, error) {
	db := DBFromContext(ctx, q.db)

	var model CustomerModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("customer %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	var addressModels []CustomerAddressModel
	if err := db.WithContext(ctx).Where("customer_id = ?", id).Find(&addressModels).Error; err != nil {
		return nil, err
	}

	addresses := make([]customerusecase.AddressDTO, 0, len(addressModels))
	for _, am := range addressModels {
		addresses = append(addresses, customerusecase.AddressDTO{
			ID:         am.ID,
			PostalCode: am.PostalCode,
			Prefecture: am.Prefecture,
			City:       am.City,
			Line:       am.Line,
			IsDefault:  am.IsDefault,
		})
	}

	return &customerusecase.CustomerDTO{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Addresses: addresses,
	}, nil
}
