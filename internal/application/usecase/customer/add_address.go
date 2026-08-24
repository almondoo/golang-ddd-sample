package customer

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
)

// AddAddressInput は顧客への住所追加ユースケースの入力である。
type AddAddressInput struct {
	CustomerID string
	PostalCode string
	Prefecture string
	City       string
	Line       string
}

// AddAddressOutput は顧客への住所追加ユースケースの出力である。
type AddAddressOutput struct {
	AddressID string
}

// AddAddressUseCase は既存顧客に配送先住所を追加するユースケースである。
type AddAddressUseCase struct {
	repo domaincustomer.Repository
	tx   tx.Manager
}

// NewAddAddressUseCase は AddAddressUseCase を生成する。
func NewAddAddressUseCase(repo domaincustomer.Repository, txManager tx.Manager) *AddAddressUseCase {
	return &AddAddressUseCase{repo: repo, tx: txManager}
}

// Execute は対象顧客を読み込み、住所を追加してから保存する。
//
// 「最初の住所は自動的にデフォルトになる」等の不変条件はすべて
// Customer.AddAddress（ドメイン層）の責務であり、ユースケースは
// 集約の読み込み・操作・保存というオーケストレーションに徹する。
func (u *AddAddressUseCase) Execute(ctx context.Context, in AddAddressInput) (AddAddressOutput, error) {
	customerID, err := domaincustomer.NewCustomerID(in.CustomerID)
	if err != nil {
		return AddAddressOutput{}, err
	}

	var out AddAddressOutput

	err = u.tx.Do(ctx, func(ctx context.Context) error {
		c, err := u.repo.FindByID(ctx, customerID)
		if err != nil {
			return err
		}

		addressID, err := c.AddAddress(in.PostalCode, in.Prefecture, in.City, in.Line)
		if err != nil {
			return err
		}

		if err := u.repo.Save(ctx, c); err != nil {
			return err
		}

		out.AddressID = addressID.String()
		return nil
	})
	if err != nil {
		return AddAddressOutput{}, err
	}

	return out, nil
}
