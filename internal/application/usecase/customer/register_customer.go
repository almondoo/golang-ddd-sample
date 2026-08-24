package customer

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
)

// RegisterCustomerInput は顧客登録ユースケースの入力である。
// プレゼンテーション層から渡ってくる生の文字列をそのまま受け取り、
// ユースケース内部でドメインの検証（NewCustomer）に委ねる。
type RegisterCustomerInput struct {
	Name  string
	Email string
}

// RegisterCustomerOutput は顧客登録ユースケースの出力である。
type RegisterCustomerOutput struct {
	CustomerID string
}

// RegisterCustomerUseCase は新規顧客を登録するユースケースである。
type RegisterCustomerUseCase struct {
	repo domaincustomer.Repository
	tx   tx.Manager
}

// NewRegisterCustomerUseCase は RegisterCustomerUseCase を生成する。
func NewRegisterCustomerUseCase(repo domaincustomer.Repository, txManager tx.Manager) *RegisterCustomerUseCase {
	return &RegisterCustomerUseCase{repo: repo, tx: txManager}
}

// Execute は新規顧客を登録し、生成された顧客 ID を返す。
func (u *RegisterCustomerUseCase) Execute(ctx context.Context, in RegisterCustomerInput) (RegisterCustomerOutput, error) {
	var out RegisterCustomerOutput

	err := u.tx.Do(ctx, func(ctx context.Context) error {
		c, err := domaincustomer.NewCustomer(in.Name, in.Email)
		if err != nil {
			return err
		}

		if err := u.repo.Save(ctx, c); err != nil {
			return err
		}

		out.CustomerID = c.ID().String()
		return nil
	})
	if err != nil {
		return RegisterCustomerOutput{}, err
	}

	return out, nil
}
