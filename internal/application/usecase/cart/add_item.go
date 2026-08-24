package cart

import (
	"context"
	"errors"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincart "github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// AddItemInput はカートへの商品追加ユースケースの入力である。
// プレゼンテーション層から渡ってくる生の文字列をそのまま受け取り、
// ユースケース内部でドメインの値オブジェクト（CustomerID/ProductID）に
// 変換する。入力検証をユースケースに閉じ込め、プレゼンテーション層は
// 単なる転送に徹する。
type AddItemInput struct {
	CustomerID string
	ProductID  string
	Quantity   int
}

// AddItemUseCase はカートに商品を追加するユースケースである。
type AddItemUseCase struct {
	repo domaincart.Repository
	tx   tx.Manager
}

// NewAddItemUseCase は AddItemUseCase を生成する。
func NewAddItemUseCase(repo domaincart.Repository, txManager tx.Manager) *AddItemUseCase {
	return &AddItemUseCase{repo: repo, tx: txManager}
}

// Execute は入力を検証したうえで、対象顧客のカートに商品を追加する。
//
// find-or-create（見つからなければ新規作成する）方針を採っている点に注意する。
// カートは「注文もカート作成も明示的に行わず、最初に商品を入れた瞬間に
// 暗黙的に生まれる」という業務上の性質を持つため、事前に空カートを
// 作成しておく別ユースケースを用意するよりも、追加操作自体に
// find-or-create の責任を持たせたほうが自然である。
func (u *AddItemUseCase) Execute(ctx context.Context, in AddItemInput) error {
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
			if errors.Is(err, shared.ErrNotFound) {
				// 初回追加であればカートはまだ存在しない。これはエラーではなく
				// 「これから作る」という正常系の分岐として扱う。
				c = domaincart.NewCart(customerID)
			} else {
				return err
			}
		}

		if err := c.AddItem(productID, in.Quantity); err != nil {
			return err
		}

		return u.repo.Save(ctx, c)
	})
}
