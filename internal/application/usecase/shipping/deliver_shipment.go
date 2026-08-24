package shipping

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domainshipping "github.com/almondoo/golang-ddd-sample/internal/domain/shipping"
)

// DeliverShipmentInput は配達完了ユースケースの入力である。
type DeliverShipmentInput struct {
	ShipmentID string
}

// DeliverShipmentUseCase は配送を配達完了にし、状態を shipped から
// delivered へ遷移させるユースケースである（order.ShipOrderUseCase と
// 対になる構造）。
type DeliverShipmentUseCase struct {
	repo      domainshipping.Repository
	txManager tx.Manager
}

// NewDeliverShipmentUseCase は DeliverShipmentUseCase を生成する。
func NewDeliverShipmentUseCase(repo domainshipping.Repository, txManager tx.Manager) *DeliverShipmentUseCase {
	return &DeliverShipmentUseCase{repo: repo, txManager: txManager}
}

// Execute は配達完了ユースケースを実行する。
func (uc *DeliverShipmentUseCase) Execute(ctx context.Context, in DeliverShipmentInput) error {
	shipmentID, err := domainshipping.NewShipmentID(in.ShipmentID)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		s, err := uc.repo.FindByID(ctx, shipmentID)
		if err != nil {
			return err
		}

		if err := s.MarkDelivered(); err != nil {
			return err
		}

		return uc.repo.Save(ctx, s)
	})
}
