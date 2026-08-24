package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/application/catalog/query"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ProductQuery は query.ProductQueryService の GORM 実装である。
//
// 型名を ProductQueryService としなかったのは、query パッケージの
// インターフェース名（query.ProductQueryService）とこの実装の型名が
// 同名になり紛らわしくなるのを避けるためである。command 側の
// ProductRepository が「catalog.Repository の実装」であるのと対比して、
// こちらは「query.ProductQueryService の実装」であることを var _ の
// アサーションで明示している。
//
// catalog.Repository を経由せず gorm.DB へ直接クエリを発行しているのは、
// query_service.go に書いた通り、読み取り専用の問い合わせにドメイン集約の
// 組み立てコストを払わせないための軽量 CQRS の実践である。
type ProductQuery struct {
	db *gorm.DB
}

// NewProductQuery は ProductQuery を生成する。
func NewProductQuery(db *gorm.DB) *ProductQuery {
	return &ProductQuery{db: db}
}

// コンパイル時に ProductQuery が query.ProductQueryService を満たすことを保証する。
var _ query.ProductQueryService = (*ProductQuery)(nil)

// List は登録済みの商品を一覧で返す。
func (q *ProductQuery) List(ctx context.Context) ([]query.ProductDTO, error) {
	db := DBFromContext(ctx, q.db)

	var models []ProductModel
	if err := db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	dtos := make([]query.ProductDTO, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, toProductDTO(m))
	}
	return dtos, nil
}

// FindByID は id に対応する商品を返す。
func (q *ProductQuery) FindByID(ctx context.Context, id string) (*query.ProductDTO, error) {
	db := DBFromContext(ctx, q.db)

	var model ProductModel
	if err := db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	dto := toProductDTO(model)
	return &dto, nil
}

// toProductDTO は永続化モデルを問い合わせ用 DTO に変換する。
// ドメイン集約を経由しないため、ここでの変換はバリデーションを伴わない
// 単純なフィールドのコピーである。
func toProductDTO(m ProductModel) query.ProductDTO {
	return query.ProductDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		PriceAmount: m.PriceAmount,
		Currency:    m.Currency,
	}
}
