package command

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
)

// ShipOrderInput は発送ユースケースの入力である。
type ShipOrderInput struct {
	OrderID string
}

// ShipOrderUseCase は注文を発送済みにし、状態を paid から shipped へ
// 遷移させるユースケースである（PayOrderUseCase と対になる構造）。
type ShipOrderUseCase struct {
	repo      order.Repository
	txManager tx.Manager
}

// NewShipOrderUseCase は ShipOrderUseCase を生成する。
func NewShipOrderUseCase(repo order.Repository, txManager tx.Manager) *ShipOrderUseCase {
	return &ShipOrderUseCase{repo: repo, txManager: txManager}
}

// Execute は発送ユースケースを実行する。
func (uc *ShipOrderUseCase) Execute(ctx context.Context, in ShipOrderInput) error {
	orderID, err := order.NewOrderID(in.OrderID)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		o, err := uc.repo.FindByID(ctx, orderID)
		if err != nil {
			return err
		}

		if err := o.Ship(); err != nil {
			return err
		}

		return uc.repo.Save(ctx, o)
	})
}
