package order

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domainorder "github.com/almondoo/golang-ddd-sample/internal/domain/order"
)

// PayOrderInput は入金確認ユースケースの入力である。
type PayOrderInput struct {
	OrderID string
}

// PayOrderUseCase は注文の入金を確認し、状態を pending から paid へ
// 遷移させるユースケースである。
//
// 「読み込み → ドメインロジックの実行（Order.Pay） → 保存」という
// register_product / change_price と同じ形の薄いユースケースである。
// 状態遷移が許されるかどうかの判断は一切ここに書かず、すべて
// Order.Pay に委譲する（ドメインロジックをアプリケーション層に
// 漏らさないための一貫した方針）。
type PayOrderUseCase struct {
	repo      domainorder.Repository
	txManager tx.Manager
}

// NewPayOrderUseCase は PayOrderUseCase を生成する。
func NewPayOrderUseCase(repo domainorder.Repository, txManager tx.Manager) *PayOrderUseCase {
	return &PayOrderUseCase{repo: repo, txManager: txManager}
}

// Execute は入金確認ユースケースを実行する。
func (uc *PayOrderUseCase) Execute(ctx context.Context, in PayOrderInput) error {
	orderID, err := domainorder.NewOrderID(in.OrderID)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		o, err := uc.repo.FindByID(ctx, orderID)
		if err != nil {
			return err
		}

		if err := o.Pay(); err != nil {
			return err
		}

		return uc.repo.Save(ctx, o)
	})
}
