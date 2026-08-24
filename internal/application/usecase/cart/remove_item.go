package cart

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincart "github.com/almondoo/golang-ddd-sample/internal/domain/cart"
)

// RemoveItemInput はカートからの商品削除ユースケースの入力である。
type RemoveItemInput struct {
	CustomerID string
	ProductID  string
}

// RemoveItemUseCase はカートから商品を削除するユースケースである。
type RemoveItemUseCase struct {
	repo domaincart.Repository
	tx   tx.Manager
}

// NewRemoveItemUseCase は RemoveItemUseCase を生成する。
func NewRemoveItemUseCase(repo domaincart.Repository, txManager tx.Manager) *RemoveItemUseCase {
	return &RemoveItemUseCase{repo: repo, tx: txManager}
}

// Execute は対象顧客のカートから指定商品を削除する。
//
// AddItem とは異なり find-or-create は行わない。カート自体が存在しない
// 状態から何かを削除しようとするのは呼び出し側の誤りである可能性が高く、
// FindByCustomerID が返す shared.ErrNotFound をそのまま呼び出し元へ
// 伝播させる。これはプレゼンテーション層で HTTP 404 に変換される。
func (u *RemoveItemUseCase) Execute(ctx context.Context, in RemoveItemInput) error {
	customerID, err := domaincart.NewCustomerID(in.CustomerID)
	if err != nil {
		return err
	}
	productID, err := domaincart.NewProductID(in.ProductID)
	if err != nil {
		return err
	}

	return u.tx.Do(ctx, func(ctx context.Context) error {
		c, err := u.repo.FindByCustomerID(ctx, customerID)
		if err != nil {
			return err
		}

		if err := c.RemoveItem(productID); err != nil {
			return err
		}

		return u.repo.Save(ctx, c)
	})
}
