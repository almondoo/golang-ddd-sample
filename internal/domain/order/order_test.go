package order_test

import (
	"testing"
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// fixedPlacedAt はテスト全体で使い回す固定の注文確定時刻である。
// PlacedAt の具体的な値そのものは状態遷移テストの関心事ではないため、
// テストごとに time.Now() を呼ぶより固定値のほうが意図が明確になる。
var fixedPlacedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// mustMoney はテスト用のヘルパーで、Money の生成失敗をテスト失敗として扱う。
// 金額の妥当性は shared パッケージ側のテストで既に検証済みなので、
// ここでは「正しく生成できる前提の値」を手早く作るためだけに使う。
func mustMoney(t *testing.T, amount int64) shared.Money {
	t.Helper()
	m, err := shared.NewMoney(amount, shared.JPY)
	if err != nil {
		t.Fatalf("failed to build money fixture: %v", err)
	}
	return m
}

// mustOrderItem はテスト用のヘルパーで、OrderItem の生成失敗をテスト失敗として扱う。
func mustOrderItem(t *testing.T, productID, productName string, unitPriceAmount int64, quantity int) order.OrderItem {
	t.Helper()
	item, err := order.NewOrderItem(productID, productName, mustMoney(t, unitPriceAmount), quantity)
	if err != nil {
		t.Fatalf("failed to build order item fixture: %v", err)
	}
	return item
}

// mustCustomerID はテスト用のヘルパーである。
func mustCustomerID(t *testing.T) order.CustomerID {
	t.Helper()
	id, err := order.NewCustomerID("customer-1")
	if err != nil {
		t.Fatalf("failed to build customer id fixture: %v", err)
	}
	return id
}

// TestNewOrder は「注文確定」というユースケースの入口にあたる
// コンストラクタが、業務ルール（不変条件）をどう守るかを検証する。
func TestNewOrder(t *testing.T) {
	t.Run("order with items is created as pending", func(t *testing.T) {
		customerID := mustCustomerID(t)
		item := mustOrderItem(t, "product-1", "商品A", 1000, 2)

		o, err := order.NewOrder(customerID, []order.OrderItem{item}, fixedPlacedAt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if o.ID() == "" {
			t.Fatalf("expected generated ID, got empty")
		}
		if o.CustomerID() != customerID {
			t.Fatalf("CustomerID() = %v, want %v", o.CustomerID(), customerID)
		}
		if o.Status() != order.StatusPending {
			t.Fatalf("Status() = %v, want %v", o.Status(), order.StatusPending)
		}
		if len(o.Items()) != 1 {
			t.Fatalf("len(Items()) = %d, want 1", len(o.Items()))
		}
		// PlacedAt は NewOrder が time.Now() ではなく引数 now をそのまま
		// 採用することを検証する（ドメイン層の決定性の担保）。
		if !o.PlacedAt().Equal(fixedPlacedAt) {
			t.Fatalf("PlacedAt() = %v, want %v", o.PlacedAt(), fixedPlacedAt)
		}
	})

	t.Run("order without items is rejected as a domain rule violation", func(t *testing.T) {
		customerID := mustCustomerID(t)

		_, err := order.NewOrder(customerID, nil, fixedPlacedAt)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}

// TestOrder_StateTransitions は、状態機械の遷移マトリクスを網羅的に検証する。
// 4 状態（pending/paid/shipped/canceled）× 3 操作（Pay/Ship/Cancel）の
// 全 12 通りについて、許可される遷移とドメインルール違反になる遷移の
// 両方を明示的にテーブル化する。
func TestOrder_StateTransitions(t *testing.T) {
	type transition struct {
		op         string
		apply      func(o *order.Order) error
		wantStatus order.Status
	}

	pay := transition{op: "Pay", apply: (*order.Order).Pay, wantStatus: order.StatusPaid}
	ship := transition{op: "Ship", apply: (*order.Order).Ship, wantStatus: order.StatusShipped}
	cancel := transition{op: "Cancel", apply: (*order.Order).Cancel, wantStatus: order.StatusCanceled}

	tests := []struct {
		from    order.Status
		trans   transition
		wantErr bool
	}{
		// pending からの遷移
		{from: order.StatusPending, trans: pay, wantErr: false},
		{from: order.StatusPending, trans: ship, wantErr: true},
		{from: order.StatusPending, trans: cancel, wantErr: false},

		// paid からの遷移
		{from: order.StatusPaid, trans: pay, wantErr: true},
		{from: order.StatusPaid, trans: ship, wantErr: false},
		{from: order.StatusPaid, trans: cancel, wantErr: false},

		// shipped からの遷移（すべて拒否される終端に近い状態）
		{from: order.StatusShipped, trans: pay, wantErr: true},
		{from: order.StatusShipped, trans: ship, wantErr: true},
		{from: order.StatusShipped, trans: cancel, wantErr: true},

		// canceled からの遷移（すべて拒否される終端状態）
		{from: order.StatusCanceled, trans: pay, wantErr: true},
		{from: order.StatusCanceled, trans: ship, wantErr: true},
		{from: order.StatusCanceled, trans: cancel, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"/"+tt.trans.op, func(t *testing.T) {
			customerID := mustCustomerID(t)
			item := mustOrderItem(t, "product-1", "商品A", 1000, 1)
			o := order.ReconstructOrder(order.OrderID("order-1"), customerID, []order.OrderItem{item}, tt.from, fixedPlacedAt, "", shared.Money{})

			err := tt.trans.apply(o)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s() from %s: expected error, got nil", tt.trans.op, tt.from)
				}
				if !shared.IsDomainRuleError(err) {
					t.Fatalf("%s() from %s: expected domain rule error, got %v", tt.trans.op, tt.from, err)
				}
				if o.Status() != tt.from {
					t.Fatalf("%s() from %s: status changed to %s despite rejection", tt.trans.op, tt.from, o.Status())
				}
				return
			}
			if err != nil {
				t.Fatalf("%s() from %s: unexpected error: %v", tt.trans.op, tt.from, err)
			}
			if o.Status() != tt.trans.wantStatus {
				t.Fatalf("%s() from %s: Status() = %v, want %v", tt.trans.op, tt.from, o.Status(), tt.trans.wantStatus)
			}
		})
	}
}

// TestOrder_TotalAmount は、明細の小計を正しく合算することを確認する。
func TestOrder_TotalAmount(t *testing.T) {
	customerID := mustCustomerID(t)
	items := []order.OrderItem{
		mustOrderItem(t, "product-1", "商品A", 1000, 2), // 2000
		mustOrderItem(t, "product-2", "商品B", 500, 3),  // 1500
	}

	o, err := order.NewOrder(customerID, items, fixedPlacedAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total, err := o.TotalAmount()
	if err != nil {
		t.Fatalf("TotalAmount() unexpected error: %v", err)
	}
	if total.Amount() != 3500 {
		t.Fatalf("TotalAmount().Amount() = %d, want 3500", total.Amount())
	}
	if total.Currency() != shared.JPY {
		t.Fatalf("TotalAmount().Currency() = %v, want %v", total.Currency(), shared.JPY)
	}
}

// TestOrder_PayableAmount_WithoutDiscount は、クーポン未適用時の
// PayableAmount が TotalAmount と一致することを確認する。
func TestOrder_PayableAmount_WithoutDiscount(t *testing.T) {
	customerID := mustCustomerID(t)
	item := mustOrderItem(t, "product-1", "商品A", 1000, 2) // 2000

	o, err := order.NewOrder(customerID, []order.OrderItem{item}, fixedPlacedAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payable, err := o.PayableAmount()
	if err != nil {
		t.Fatalf("PayableAmount() unexpected error: %v", err)
	}
	if payable.Amount() != 2000 {
		t.Fatalf("PayableAmount().Amount() = %d, want 2000", payable.Amount())
	}
}

// TestOrder_ApplyDiscount は、割引適用の可否を Order 集約自身が守る
// 不変条件（pending 限定・二重適用禁止・合計超過禁止）を検証する。
func TestOrder_ApplyDiscount(t *testing.T) {
	t.Run("succeeds on a pending order and updates PayableAmount", func(t *testing.T) {
		customerID := mustCustomerID(t)
		item := mustOrderItem(t, "product-1", "商品A", 1000, 2) // 2000
		o, err := order.NewOrder(customerID, []order.OrderItem{item}, fixedPlacedAt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := o.ApplyDiscount("SUMMER-500", mustMoney(t, 500)); err != nil {
			t.Fatalf("ApplyDiscount() unexpected error: %v", err)
		}
		if o.CouponCode() != "SUMMER-500" {
			t.Fatalf("CouponCode() = %q, want %q", o.CouponCode(), "SUMMER-500")
		}
		if o.DiscountAmount().Amount() != 500 {
			t.Fatalf("DiscountAmount().Amount() = %d, want 500", o.DiscountAmount().Amount())
		}

		payable, err := o.PayableAmount()
		if err != nil {
			t.Fatalf("PayableAmount() unexpected error: %v", err)
		}
		if payable.Amount() != 1500 {
			t.Fatalf("PayableAmount().Amount() = %d, want 1500", payable.Amount())
		}
	})

	t.Run("rejected on a non-pending order", func(t *testing.T) {
		customerID := mustCustomerID(t)
		item := mustOrderItem(t, "product-1", "商品A", 1000, 2)
		o := order.ReconstructOrder(order.OrderID("order-1"), customerID, []order.OrderItem{item}, order.StatusPaid, fixedPlacedAt, "", shared.Money{})

		err := o.ApplyDiscount("SUMMER-500", mustMoney(t, 500))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("rejected when a coupon is already applied", func(t *testing.T) {
		customerID := mustCustomerID(t)
		item := mustOrderItem(t, "product-1", "商品A", 1000, 2)
		o, err := order.NewOrder(customerID, []order.OrderItem{item}, fixedPlacedAt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := o.ApplyDiscount("FIRST-100", mustMoney(t, 100)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = o.ApplyDiscount("SECOND-100", mustMoney(t, 100))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		// 最初に適用したクーポンの状態が保持されたままであることも確認する。
		if o.CouponCode() != "FIRST-100" {
			t.Fatalf("CouponCode() = %q, want %q", o.CouponCode(), "FIRST-100")
		}
	})

	t.Run("rejected when discount exceeds total amount", func(t *testing.T) {
		customerID := mustCustomerID(t)
		item := mustOrderItem(t, "product-1", "商品A", 1000, 1) // 1000
		o, err := order.NewOrder(customerID, []order.OrderItem{item}, fixedPlacedAt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = o.ApplyDiscount("TOO-MUCH", mustMoney(t, 1001))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}
