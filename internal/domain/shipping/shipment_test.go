package shipping_test

import (
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shipping"
)

// mustOrderID はテスト用のヘルパーで、OrderID の生成失敗をテスト失敗として扱う。
func mustOrderID(t *testing.T) shipping.OrderID {
	t.Helper()
	id, err := shipping.NewOrderID("order-1")
	if err != nil {
		t.Fatalf("failed to build order id fixture: %v", err)
	}
	return id
}

// TestNewShipment は「配送を生成する」というユースケースの入口にあたる
// コンストラクタが、業務ルール（不変条件）をどう守るかを検証する。
func TestNewShipment(t *testing.T) {
	t.Run("shipment with address is created as preparing", func(t *testing.T) {
		orderID := mustOrderID(t)

		s, err := shipping.NewShipment(orderID, "東京都千代田区1-1-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.ID() == "" {
			t.Fatalf("expected generated ID, got empty")
		}
		if s.OrderID() != orderID {
			t.Fatalf("OrderID() = %v, want %v", s.OrderID(), orderID)
		}
		if s.Address() != "東京都千代田区1-1-1" {
			t.Fatalf("Address() = %v, want %v", s.Address(), "東京都千代田区1-1-1")
		}
		if s.Status() != shipping.StatusPreparing {
			t.Fatalf("Status() = %v, want %v", s.Status(), shipping.StatusPreparing)
		}
	})

	t.Run("shipment without address is rejected as a domain rule violation", func(t *testing.T) {
		orderID := mustOrderID(t)

		_, err := shipping.NewShipment(orderID, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}

// TestShipment_StateTransitions は、状態機械の遷移マトリクスを網羅的に検証する。
// 3 状態（preparing/shipped/delivered）× 2 操作（MarkShipped/MarkDelivered）の
// 全 6 通りについて、許可される遷移とドメインルール違反になる遷移の
// 両方を明示的にテーブル化する（order.Order の状態遷移テストと同じ方針）。
func TestShipment_StateTransitions(t *testing.T) {
	type transition struct {
		op         string
		apply      func(s *shipping.Shipment) error
		wantStatus shipping.Status
	}

	markShipped := transition{op: "MarkShipped", apply: (*shipping.Shipment).MarkShipped, wantStatus: shipping.StatusShipped}
	markDelivered := transition{op: "MarkDelivered", apply: (*shipping.Shipment).MarkDelivered, wantStatus: shipping.StatusDelivered}

	tests := []struct {
		from    shipping.Status
		trans   transition
		wantErr bool
	}{
		// preparing からの遷移
		{from: shipping.StatusPreparing, trans: markShipped, wantErr: false},
		{from: shipping.StatusPreparing, trans: markDelivered, wantErr: true},

		// shipped からの遷移
		{from: shipping.StatusShipped, trans: markShipped, wantErr: true},
		{from: shipping.StatusShipped, trans: markDelivered, wantErr: false},

		// delivered からの遷移（すべて拒否される終端状態）
		{from: shipping.StatusDelivered, trans: markShipped, wantErr: true},
		{from: shipping.StatusDelivered, trans: markDelivered, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"/"+tt.trans.op, func(t *testing.T) {
			orderID := mustOrderID(t)
			s := shipping.ReconstructShipment(shipping.ShipmentID("shipment-1"), orderID, "東京都千代田区1-1-1", tt.from)

			err := tt.trans.apply(s)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s() from %s: expected error, got nil", tt.trans.op, tt.from)
				}
				if !shared.IsDomainRuleError(err) {
					t.Fatalf("%s() from %s: expected domain rule error, got %v", tt.trans.op, tt.from, err)
				}
				if s.Status() != tt.from {
					t.Fatalf("%s() from %s: status changed to %s despite rejection", tt.trans.op, tt.from, s.Status())
				}
				return
			}
			if err != nil {
				t.Fatalf("%s() from %s: unexpected error: %v", tt.trans.op, tt.from, err)
			}
			if s.Status() != tt.trans.wantStatus {
				t.Fatalf("%s() from %s: Status() = %v, want %v", tt.trans.op, tt.from, s.Status(), tt.trans.wantStatus)
			}
		})
	}
}
