package order

import (
	"context"
	"errors"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	domainorder "github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// CancelOrderInput は注文取り消しユースケースの入力である。
type CancelOrderInput struct {
	OrderID string
}

// CancelOrderUseCase は注文を取り消し、状態を canceled へ遷移させる
// ユースケースである（PayOrderUseCase / ShipOrderUseCase と対になる構造）。
// pending・paid のいずれからも取り消せるという遷移可否の判断は
// Order.Cancel に委譲する。
//
// 【キャンセルで引当を戻す】
// 注文確定時（PlaceOrderUseCase）に inventory.Stock.Reserve で引き当てた
// 在庫は、注文が取り消された場合には解放されなければならない
// （さもないと「もう存在しない注文」のために在庫が永久に取り置かれたままになる）。
// これも order 集約単体では守れない、コンテキストをまたぐ整合性であるため
// アプリケーション層が担う。
type CancelOrderUseCase struct {
	orderRepo     domainorder.Repository
	inventoryRepo domaininventory.Repository
	txManager     tx.Manager
}

// NewCancelOrderUseCase は CancelOrderUseCase を生成する。
func NewCancelOrderUseCase(
	orderRepo domainorder.Repository,
	inventoryRepo domaininventory.Repository,
	txManager tx.Manager,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{orderRepo: orderRepo, inventoryRepo: inventoryRepo, txManager: txManager}
}

// Execute は注文取り消しユースケースを実行する。
func (uc *CancelOrderUseCase) Execute(ctx context.Context, in CancelOrderInput) error {
	orderID, err := domainorder.NewOrderID(in.OrderID)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		o, err := uc.orderRepo.FindByID(ctx, orderID)
		if err != nil {
			return err
		}

		if err := o.Cancel(); err != nil {
			return err
		}

		// 注文明細ごとに引当を解除する。実在庫（quantity）は変えず、
		// 引当済み数量（reserved）だけを戻す（inventory.Stock.Release を参照）。
		for _, item := range o.Items() {
			productID, err := domaininventory.NewProductID(item.ProductID())
			if err != nil {
				return err
			}
			stock, err := uc.inventoryRepo.FindByProductID(ctx, productID)
			if err != nil {
				if errors.Is(err, shared.ErrNotFound) {
					// 注文確定時点で在庫は必ず引き当て済みのはずであり、
					// ここで見つからないのはその不変条件が崩れた想定外の
					// 状態である。ship_order.go と同じ理由でドメインルール
					// 違反として表現する。
					return shared.NewDomainRuleError("order: stock for product %s not found", productID)
				}
				return err
			}
			if err := stock.Release(item.Quantity()); err != nil {
				return err
			}
			if err := uc.inventoryRepo.Save(ctx, stock); err != nil {
				return err
			}
		}

		return uc.orderRepo.Save(ctx, o)
	})
}
