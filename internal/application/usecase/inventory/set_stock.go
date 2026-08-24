package inventory

import (
	"context"
	"errors"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// SetStockInput は在庫数設定ユースケースの入力である。
// プレゼンテーション層から渡ってくる生の文字列・数値をそのまま受け取り、
// ユースケース内部でドメインの値オブジェクト（ProductID）に変換する。
type SetStockInput struct {
	ProductID string
	Quantity  int
}

// SetStockUseCase は商品の在庫数を設定するユースケースである。
type SetStockUseCase struct {
	repo domaininventory.Repository
	tx   tx.Manager
}

// NewSetStockUseCase は SetStockUseCase を生成する。
func NewSetStockUseCase(repo domaininventory.Repository, txManager tx.Manager) *SetStockUseCase {
	return &SetStockUseCase{repo: repo, tx: txManager}
}

// Execute は指定商品の在庫数を設定する。
//
// find-or-create（見つからなければ新規作成する）方針を採っている点に注意する。
// 在庫は「商品登録とは別に、在庫担当者が最初に数量を入力した瞬間に
// 暗黙的に生まれる」という業務上の性質を持ちうるため、事前に在庫レコードを
// 作成しておく別ユースケースを用意するよりも、設定操作自体に
// find-or-create の責任を持たせたほうが自然である（cart.AddItemUseCase と同じ判断）。
func (u *SetStockUseCase) Execute(ctx context.Context, in SetStockInput) error {
	productID, err := domaininventory.NewProductID(in.ProductID)
	if err != nil {
		return err
	}

	return u.tx.Do(ctx, func(ctx context.Context) error {
		s, err := u.repo.FindByProductID(ctx, productID)
		if err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				// 初回設定であれば在庫はまだ存在しない。これはエラーではなく
				// 「これから作る」という正常系の分岐として扱う。
				s, err = domaininventory.NewStock(productID, in.Quantity)
				if err != nil {
					return err
				}
				return u.repo.Save(ctx, s)
			}
			return err
		}

		if err := s.SetQuantity(in.Quantity); err != nil {
			return err
		}

		return u.repo.Save(ctx, s)
	})
}
