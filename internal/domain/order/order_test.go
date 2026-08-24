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

		o, err := order.NewOrder(customerID, []order.OrderItem{item})
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
	})

	t.Run("order without items is rejected as a domain rule violation", func(t *testing.T) {
		customerID := mustCustomerID(t)

		_, err := order.NewOrder(customerID, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}

// TestNewOrder_RecordsOrderPlacedEvent は、NewOrder が OrderPlaced イベントを
// ちょうど 1 件記録すること、および PullEvents が「一度取り出したら空になる」
// という二重配信防止の契約を守っていることを確認する。
func TestNewOrder_RecordsOrderPlacedEvent(t *testing.T) {
	customerID := mustCustomerID(t)
	item := mustOrderItem(t, "product-1", "商品A", 1000, 1)

	o, err := order.NewOrder(customerID, []order.OrderItem{item})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := o.PullEvents()
	if len(events) != 1 {
		t.Fatalf("len(PullEvents()) = %d, want 1", len(events))
	}
	if events[0].EventName() != order.OrderPlacedEventName {
		t.Fatalf("EventName() = %q, want %q", events[0].EventName(), order.OrderPlacedEventName)
	}
	placed, ok := events[0].(order.OrderPlaced)
	if !ok {
		t.Fatalf("event type = %T, want order.OrderPlaced", events[0])
	}
	if placed.OrderID() != o.ID() {
		t.Errorf("OrderID() = %v, want %v", placed.OrderID(), o.ID())
	}
	if placed.CustomerID() != customerID {
		t.Errorf("CustomerID() = %v, want %v", placed.CustomerID(), customerID)
	}

	// 2 回目の PullEvents は空でなければならない（二重配信防止）。
	second := o.PullEvents()
	if len(second) != 0 {
		t.Fatalf("second PullEvents() = %d events, want 0", len(second))
	}
}

// TestOrder_ReconstructOrder_DoesNotRecordEvent は、DB からの復元では
// OrderPlaced イベントが再記録されないことを確認する。すでに確定済みの
// 過去の注文を読み込むたびにイベントが再発火すると、カートが誤って
// 何度も空にされる等の意図しない副作用が起きてしまうためである。
func TestOrder_ReconstructOrder_DoesNotRecordEvent(t *testing.T) {
	customerID := mustCustomerID(t)
	item := mustOrderItem(t, "product-1", "商品A", 1000, 1)

	o := order.ReconstructOrder(order.OrderID("order-1"), customerID, []order.OrderItem{item}, order.StatusPending, fixedPlacedAt)

	events := o.PullEvents()
	if len(events) != 0 {
		t.Fatalf("PullEvents() after ReconstructOrder = %d events, want 0", len(events))
	}
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
			o := order.ReconstructOrder(order.OrderID("order-1"), customerID, []order.OrderItem{item}, tt.from, fixedPlacedAt)

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

	o, err := order.NewOrder(customerID, items)
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
