package catalog

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincatalog "github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ChangePriceInput はユースケースの入力データである。
type ChangePriceInput struct {
	ProductID      string
	NewPriceAmount int64
}

// ChangePriceUseCase は既存商品の価格を変更するユースケースである。
type ChangePriceUseCase struct {
	repo      domaincatalog.Repository
	txManager tx.Manager
}

// NewChangePriceUseCase は ChangePriceUseCase を生成する。
func NewChangePriceUseCase(repo domaincatalog.Repository, txManager tx.Manager) *ChangePriceUseCase {
	return &ChangePriceUseCase{repo: repo, txManager: txManager}
}

// Execute は価格変更ユースケースを実行する。
//
// 「読み込み → ドメインロジックの実行 → 保存」という一連の流れを 1 つの
// トランザクションに包む。読み込みと保存の間で他のトランザクションが
// 割り込んでも問題ないかは、実際のプロダクションでは楽観ロック等の
// 追加対策を要する場合があるが、本サンプルでは学習用に単純な形に留めている。
func (uc *ChangePriceUseCase) Execute(ctx context.Context, input ChangePriceInput) error {
	productID, err := domaincatalog.NewProductID(input.ProductID)
	if err != nil {
		return err
	}

	newPrice, err := shared.NewMoney(input.NewPriceAmount, shared.JPY)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		product, err := uc.repo.FindByID(ctx, productID)
		if err != nil {
			return err
		}

		if err := product.ChangePrice(newPrice); err != nil {
			return err
		}

		return uc.repo.Save(ctx, product)
	})
}
