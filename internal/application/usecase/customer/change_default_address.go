package customer

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
)

// ChangeDefaultAddressInput はデフォルト住所変更ユースケースの入力である。
type ChangeDefaultAddressInput struct {
	CustomerID string
	AddressID  string
}

// ChangeDefaultAddressUseCase は顧客のデフォルト住所を変更するユースケースである。
type ChangeDefaultAddressUseCase struct {
	repo domaincustomer.Repository
	tx   tx.Manager
}

// NewChangeDefaultAddressUseCase は ChangeDefaultAddressUseCase を生成する。
func NewChangeDefaultAddressUseCase(repo domaincustomer.Repository, txManager tx.Manager) *ChangeDefaultAddressUseCase {
	return &ChangeDefaultAddressUseCase{repo: repo, tx: txManager}
}

// Execute は対象顧客を読み込み、指定された住所を新しいデフォルトに変更してから
// 保存する。
//
// 「既存のデフォルトを降ろし、新しいデフォルトを立てる」という 2 ステップを
// 1 つの不可分な操作として扱うのは Customer.ChangeDefaultAddress
// （ドメイン層）の責務である。ユースケースはその呼び出しを 1 つの
// トランザクションで包むだけでよい。
func (u *ChangeDefaultAddressUseCase) Execute(ctx context.Context, in ChangeDefaultAddressInput) error {
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

		if err := c.ChangeDefaultAddress(addressID); err != nil {
			return err
		}

		return u.repo.Save(ctx, c)
	})
}
