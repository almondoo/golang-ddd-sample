package customer

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
)

// RemoveAddressInput は住所削除ユースケースの入力である。
type RemoveAddressInput struct {
	CustomerID string
	AddressID  string
}

// RemoveAddressUseCase は顧客の住所を削除するユースケースである。
type RemoveAddressUseCase struct {
	repo domaincustomer.Repository
	tx   tx.Manager
}

// NewRemoveAddressUseCase は RemoveAddressUseCase を生成する。
func NewRemoveAddressUseCase(repo domaincustomer.Repository, txManager tx.Manager) *RemoveAddressUseCase {
	return &RemoveAddressUseCase{repo: repo, tx: txManager}
}

// Execute は対象顧客を読み込み、指定された住所を削除してから保存する。
//
// 「デフォルト住所は他に住所が残る場合は削除できない」というルールは
// Customer.RemoveAddress（ドメイン層）が判断する。ユースケースはその
// エラーをそのまま呼び出し元へ伝播させるだけでよい。
func (u *RemoveAddressUseCase) Execute(ctx context.Context, in RemoveAddressInput) error {
	customerID, err := domaincustomer.NewCustomerID(in.CustomerID)
	if err != nil {
		return err
	}
	addressID, err := domaincustomer.NewAddressID(in.AddressID)
	if err != nil {
		return err
	}

	return u.tx.Do(ctx, func(ctx context.Context) error {
		c, err := u.repo.FindByID(ctx, customerID)
		if err != nil {
			return err
		}

		if err := c.RemoveAddress(addressID); err != nil {
			return err
		}

		return u.repo.Save(ctx, c)
	})
}
