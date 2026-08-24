package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/application/usecase/order"
	domaincoupon "github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	domainorder "github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// このファイルは CancelOrderUseCase（orderRepo / inventoryRepo / couponRepo +
// tx.Manager に依存するユースケース）を検証する。place_order_test.go で
// 定義済みの手書きフェイク（fakeOrderRepository 等）・共通フィクスチャ
// （testCustomerID, mustMoney, farFutureExpiry）を同じパッケージ内で再利用し、
// 重複するフェイク定義は行わない。

const cancelTestProductID = "product-1"

// cancelOrderFixture は CancelOrderUseCase 用のテストフィクスチャである。
type cancelOrderFixture struct {
	t *testing.T

	orderRepo     *fakeOrderRepository
	inventoryRepo *fakeInventoryRepository
	couponRepo    *fakeCouponRepository
	txManager     *fakeTxManager

	useCase *order.CancelOrderUseCase
}

func newCancelOrderFixture(t *testing.T) *cancelOrderFixture {
	t.Helper()
	f := &cancelOrderFixture{
		t:             t,
		orderRepo:     newFakeOrderRepository(),
		inventoryRepo: newFakeInventoryRepository(),
		couponRepo:    newFakeCouponRepository(),
		txManager:     &fakeTxManager{},
	}
	f.useCase = order.NewCancelOrderUseCase(f.orderRepo, f.inventoryRepo, f.couponRepo, f.txManager)
	return f
}

// addStock は cancelTestProductID の在庫を quantity・reserved であらかじめ
// 登録する（PlaceOrderUseCase による引当が既に完了している状態を再現する）。
func (f *cancelOrderFixture) addStock(quantity, reserved int) domaininventory.ProductID {
	f.t.Helper()
	pid, err := domaininventory.NewProductID(cancelTestProductID)
	if err != nil {
		f.t.Fatalf("failed to build inventory product id fixture: %v", err)
	}
	f.inventoryRepo.stocks[pid] = domaininventory.ReconstructStock(pid, quantity, reserved)
	return pid
}

// newOrder は取り消し対象の注文をフィクスチャとして組み立て、
// orderRepo に登録してから返す。withCoupon が true の場合は
// クーポンを1回消費した状態で注文に適用し、couponRepo にも登録する。
func (f *cancelOrderFixture) newOrder(quantity int, withCoupon bool) (*domainorder.Order, domaincoupon.CouponCode) {
	f.t.Helper()

	customerID, err := domainorder.NewCustomerID(testCustomerID)
	if err != nil {
		f.t.Fatalf("failed to build customer id fixture: %v", err)
	}
	price := mustMoney(f.t, 1000)
	item, err := domainorder.NewOrderItem(cancelTestProductID, "商品A", price, quantity)
	if err != nil {
		f.t.Fatalf("failed to build order item fixture: %v", err)
	}
	o, err := domainorder.NewOrder(customerID, []domainorder.OrderItem{item}, time.Now())
	if err != nil {
		f.t.Fatalf("failed to build order fixture: %v", err)
	}

	var code domaincoupon.CouponCode
	if withCoupon {
		code, err = domaincoupon.NewCouponCode("SUMMER10")
		if err != nil {
			f.t.Fatalf("failed to build coupon code fixture: %v", err)
		}
		cp, err := domaincoupon.NewAmountCoupon(code, mustMoney(f.t, 200), farFutureExpiry, 5)
		if err != nil {
			f.t.Fatalf("failed to build coupon fixture: %v", err)
		}
		// 注文確定時に PlaceOrderUseCase が行うのと同じ手順で、
		// あらかじめ1回消費させた状態を再現する。
		if err := cp.Use(time.Now()); err != nil {
			f.t.Fatalf("failed to pre-consume coupon fixture: %v", err)
		}
		f.couponRepo.coupons[code] = cp

		total, err := o.TotalAmount()
		if err != nil {
			f.t.Fatalf("TotalAmount() unexpected error: %v", err)
		}
		discount, err := cp.DiscountFor(total)
		if err != nil {
			f.t.Fatalf("DiscountFor() unexpected error: %v", err)
		}
		if err := o.ApplyDiscount(code.String(), discount); err != nil {
			f.t.Fatalf("ApplyDiscount() unexpected error: %v", err)
		}
	}

	f.orderRepo.orders[o.ID()] = o
	return o, code
}

// ---- 正常系 ----

func TestCancelOrderUseCase_Execute_WithCoupon_RefundsCoupon(t *testing.T) {
	f := newCancelOrderFixture(t)
	pid := f.addStock(5, 2)
	o, code := f.newOrder(2, true)

	if err := f.useCase.Execute(context.Background(), order.CancelOrderInput{OrderID: o.ID().String()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := f.orderRepo.orders[o.ID()].Status(); got != domainorder.StatusCanceled {
		t.Fatalf("Status() = %v, want %v", got, domainorder.StatusCanceled)
	}

	// 在庫の引当が解除されていること。
	if got := f.inventoryRepo.stocks[pid].Reserved(); got != 0 {
		t.Fatalf("stock Reserved() = %d, want 0", got)
	}

	// クーポンの利用実績が戻っていること（再利用可能になっている）。
	savedCoupon, ok := f.couponRepo.coupons[code]
	if !ok {
		t.Fatalf("coupon was not found")
	}
	if got := savedCoupon.UsedCount(); got != 0 {
		t.Fatalf("UsedCount() = %d, want 0", got)
	}

	// 状態検証だけでは Save が実際に呼ばれたことを証明できないため、
	// 呼び出し記録も併せて検証する（フェイクがポインタを共有するため）。
	if len(f.inventoryRepo.saved) != 1 {
		t.Fatalf("inventoryRepo.Save called %d times, want 1", len(f.inventoryRepo.saved))
	}
	if len(f.couponRepo.saved) != 1 {
		t.Fatalf("couponRepo.Save called %d times, want 1", len(f.couponRepo.saved))
	}
	if len(f.orderRepo.saved) != 1 {
		t.Fatalf("orderRepo.Save called %d times, want 1", len(f.orderRepo.saved))
	}
}

func TestCancelOrderUseCase_Execute_WithoutCoupon_CouponRepoUntouched(t *testing.T) {
	f := newCancelOrderFixture(t)
	pid := f.addStock(5, 2)
	o, _ := f.newOrder(2, false)

	if err := f.useCase.Execute(context.Background(), order.CancelOrderInput{OrderID: o.ID().String()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := f.orderRepo.orders[o.ID()].Status(); got != domainorder.StatusCanceled {
		t.Fatalf("Status() = %v, want %v", got, domainorder.StatusCanceled)
	}
	if got := f.inventoryRepo.stocks[pid].Reserved(); got != 0 {
		t.Fatalf("stock Reserved() = %d, want 0", got)
	}

	// クーポン未適用の注文なので、クーポンリポジトリには一切触れないはず。
	if f.couponRepo.findCalls != 0 {
		t.Fatalf("couponRepo.FindByCode called %d times, want 0", f.couponRepo.findCalls)
	}
	if len(f.couponRepo.saved) != 0 {
		t.Fatalf("couponRepo.Save called %d times, want 0", len(f.couponRepo.saved))
	}
}

// ---- 異常系(ドメインルール違反) ----

func TestCancelOrderUseCase_Execute_AlreadyShipped_RuleErrorAndNoSideEffects(t *testing.T) {
	f := newCancelOrderFixture(t)
	pid := f.addStock(5, 2)
	o, code := f.newOrder(2, true)

	// 発送済みまで状態を進める（pending --Pay--> paid --Ship--> shipped）。
	if err := o.Pay(); err != nil {
		t.Fatalf("Pay() unexpected error: %v", err)
	}
	if err := o.Ship(); err != nil {
		t.Fatalf("Ship() unexpected error: %v", err)
	}

	err := f.useCase.Execute(context.Background(), order.CancelOrderInput{OrderID: o.ID().String()})
	if err == nil {
		t.Fatal("expected error for canceling a shipped order, got nil")
	}
	if !shared.IsDomainRuleError(err) {
		t.Fatalf("expected domain rule error, got %v", err)
	}

	// Order.Cancel が状態遷移の可否チェックで先に弾くため、在庫の引当解除
	// もクーポンの払い戻しも一切実行されていないはずである
	// （フェイクの呼び出し記録が空であることで証明する）。
	if len(f.inventoryRepo.saved) != 0 {
		t.Fatalf("inventoryRepo.Save must not be called, but was called %d times", len(f.inventoryRepo.saved))
	}
	if got := f.inventoryRepo.stocks[pid].Reserved(); got != 2 {
		t.Fatalf("stock Reserved() = %d, want unchanged 2", got)
	}
	if len(f.couponRepo.saved) != 0 {
		t.Fatalf("couponRepo.Save must not be called, but was called %d times", len(f.couponRepo.saved))
	}
	if got := f.couponRepo.coupons[code].UsedCount(); got != 1 {
		t.Fatalf("coupon UsedCount() = %d, want unchanged 1", got)
	}
	if len(f.orderRepo.saved) != 0 {
		t.Fatalf("orderRepo.Save must not be called, but was called %d times", len(f.orderRepo.saved))
	}
}
