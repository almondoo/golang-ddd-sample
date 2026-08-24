package command

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
)

// CancelOrderInput は注文取り消しユースケースの入力である。
type CancelOrderInput struct {
	OrderID string
}

// CancelOrderUseCase は注文を取り消し、状態を canceled へ遷移させる
// ユースケースである（PayOrderUseCase / ShipOrderUseCase と対になる構造）。
// pending・paid のいずれからも取り消せるという遷移可否の判断は
// Order.Cancel に委譲する。
type CancelOrderUseCase struct {
	repo      order.Repository
	txManager tx.Manager
}

// NewCancelOrderUseCase は CancelOrderUseCase を生成する。
func NewCancelOrderUseCase(repo order.Repository, txManager tx.Manager) *CancelOrderUseCase {
	return &CancelOrderUseCase{repo: repo, txManager: txManager}
}

// Execute は注文取り消しユースケースを実行する。
func (uc *CancelOrderUseCase) Execute(ctx context.Context, in CancelOrderInput) error {
	orderID, err := order.NewOrderID(in.OrderID)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		o, err := uc.repo.FindByID(ctx, orderID)
		if err != nil {
			return err
		}

		if err := o.Cancel(); err != nil {
			return err
		}

		return uc.repo.Save(ctx, o)
	})
}
