package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	domainorder "github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
	domainshipping "github.com/almondoo/golang-ddd-sample/internal/domain/shipping"
)

// ShipOrderInput は発送ユースケースの入力である。
type ShipOrderInput struct {
	OrderID string
}

// ShipOrderUseCase は注文を発送済みにし、状態を paid から shipped へ
// 遷移させるユースケースである。
//
// 【発送 = Shipment 生成 + 在庫消込】
// 発送は単に Order の状態を進めるだけでなく、業務的には「配送を開始する
// （shipping.Shipment を生成する）」ことと「取り置きしていた在庫を実際に
// 払い出す（inventory.Stock.ConsumeReserved）」ことの両方を意味する。
// これらは互いに異なるコンテキストの集約であるため、order ドメインパッケージ
// 自身が両方の不変条件を守ることはできない。PlaceOrderUseCase と同様、
// この複数コンテキストにまたがるオーケストレーションはアプリケーション層の
// 責務である。
type ShipOrderUseCase struct {
	orderRepo     domainorder.Repository
	customerRepo  domaincustomer.Repository
	shippingRepo  domainshipping.Repository
	inventoryRepo domaininventory.Repository
	txManager     tx.Manager
}

// NewShipOrderUseCase は ShipOrderUseCase を生成する。
func NewShipOrderUseCase(
	orderRepo domainorder.Repository,
	customerRepo domaincustomer.Repository,
	shippingRepo domainshipping.Repository,
	inventoryRepo domaininventory.Repository,
	txManager tx.Manager,
) *ShipOrderUseCase {
	return &ShipOrderUseCase{
		orderRepo:     orderRepo,
		customerRepo:  customerRepo,
		shippingRepo:  shippingRepo,
		inventoryRepo: inventoryRepo,
		txManager:     txManager,
	}
}

// Execute は発送ユースケースを実行する。
//
// 処理の流れ:
//  1. 注文を読み込み、Order.Ship() で paid から shipped へ遷移させる
//  2. 注文した顧客のデフォルト配送先住所を取得し、1 行の文字列に整形する
//  3. shipping.Shipment を生成し、発送指示と同時に出荷済みにする
//  4. 注文明細ごとに、注文確定時に引き当てた在庫を消込む（実在庫を減らす）
//  5. 変更後の Order を保存する
func (uc *ShipOrderUseCase) Execute(ctx context.Context, in ShipOrderInput) error {
	orderID, err := domainorder.NewOrderID(in.OrderID)
	if err != nil {
		return err
	}

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		o, err := uc.orderRepo.FindByID(ctx, orderID)
		if err != nil {
			return err
		}

		if err := o.Ship(); err != nil {
			return err
		}

		// 2. 配送先住所の取得。
		customerID, err := domaincustomer.NewCustomerID(o.CustomerID().String())
		if err != nil {
			return err
		}
		cust, err := uc.customerRepo.FindByID(ctx, customerID)
		if err != nil {
			return err
		}
		addr, err := cust.DefaultAddress()
		if err != nil {
			// DefaultAddress は「住所が 1 件も登録されていない」場合に
			// ドメインルール違反を返すが、発送という業務の文脈でより
			// 意味の伝わるメッセージに置き換える。
			return fmt.Errorf("order: shipping address is not registered: %w", err)
		}
		addressLine := fmt.Sprintf("%s %s %s %s", addr.PostalCode(), addr.Prefecture(), addr.City(), addr.Line())

		// 3. Shipment の生成。発送指示と同時に出荷済みにする
		// （本サンプルでは「準備中」の期間を別ユースケースとして分離せず、
		// 発送指示の時点で出荷済みまで一気に進める単純化を行っている）。
		shipOrderID, err := domainshipping.NewOrderID(o.ID().String())
		if err != nil {
			return err
		}
		s, err := domainshipping.NewShipment(shipOrderID, addressLine)
		if err != nil {
			return err
		}
		if err := s.MarkShipped(); err != nil {
			return err
		}
		if err := uc.shippingRepo.Save(ctx, s); err != nil {
			return err
		}

		// 4. 在庫の消込。注文確定時点（PlaceOrderUseCase）で在庫はすでに
		// 引き当て済みであるという不変条件のもと、ここでは
		// ConsumeReserved のみを呼ぶ。万一その不変条件が崩れて在庫行その
		// ものが存在しない場合は、想定外の状態としてドメインルール違反を返す。
		for _, item := range o.Items() {
			productID, err := domaininventory.NewProductID(item.ProductID())
			if err != nil {
				return err
			}
			stock, err := uc.inventoryRepo.FindByProductID(ctx, productID)
			if err != nil {
				if errors.Is(err, shared.ErrNotFound) {
					// 注文確定時点（PlaceOrderUseCase）で在庫は必ず引き当て済み
					// のはずであり、ここで見つからないのはその不変条件が崩れた
					// 想定外の状態である。とはいえシステムエラー（500）ではなく、
					// 業務ルール違反として表現する。
					return shared.NewDomainRuleError("order: stock for product %s not found", productID)
				}
				return err
			}
			if err := stock.ConsumeReserved(item.Quantity()); err != nil {
				return err
			}
			if err := uc.inventoryRepo.Save(ctx, stock); err != nil {
				return err
			}
		}

		return uc.orderRepo.Save(ctx, o)
	})
}
