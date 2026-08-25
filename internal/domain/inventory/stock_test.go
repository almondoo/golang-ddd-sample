package inventory_test

import (
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// mustProductID はテスト用のヘルパーで、ProductID の生成失敗をテスト失敗として扱う。
func mustProductID(t *testing.T, s string) inventory.ProductID {
	t.Helper()
	id, err := inventory.NewProductID(s)
	if err != nil {
		t.Fatalf("failed to build product id fixture: %v", err)
	}
	return id
}

// TestNewStock は在庫の新規生成における妥当性検証をテーブル駆動で検証する。
func TestNewStock(t *testing.T) {
	tests := []struct {
		name        string
		quantity    int
		wantErr     bool
		wantRuleErr bool
	}{
		{name: "zero quantity is accepted", quantity: 0, wantErr: false},
		{name: "positive quantity is accepted", quantity: 10, wantErr: false},
		{name: "negative quantity is rejected", quantity: -1, wantErr: true, wantRuleErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productID := mustProductID(t, "product-1")
			s, err := inventory.NewStock(productID, tt.quantity)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantRuleErr && !shared.IsDomainRuleError(err) {
					t.Fatalf("expected domain rule error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s.ProductID() != productID {
				t.Errorf("ProductID() = %v, want %v", s.ProductID(), productID)
			}
			if s.Quantity() != tt.quantity {
				t.Errorf("Quantity() = %d, want %d", s.Quantity(), tt.quantity)
			}
			if s.Reserved() != 0 {
				t.Errorf("Reserved() = %d, want 0", s.Reserved())
			}
			// 新規生成した在庫の版数は 0 から始まる（「まだ一度も永続化されて
			// いない」ことを表す規約。実際に version=1 が採番されるのは
			// リポジトリの Save が INSERT を行った時点である）。
			if s.Version() != 0 {
				t.Errorf("Version() = %d, want 0", s.Version())
			}
		})
	}
}

// TestStock_Reserve は在庫引当（Reserve）に関する不変条件を検証する。
func TestStock_Reserve(t *testing.T) {
	t.Run("reserving within available stock succeeds", func(t *testing.T) {
		s, err := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := s.Reserve(4); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Reserved() != 4 {
			t.Errorf("Reserved() = %d, want 4", s.Reserved())
		}
		if s.Quantity() != 10 {
			t.Errorf("Quantity() = %d, want 10 (unchanged)", s.Quantity())
		}
		if s.Available() != 6 {
			t.Errorf("Available() = %d, want 6", s.Available())
		}
	})

	t.Run("reserving more than available is rejected", func(t *testing.T) {
		s, err := inventory.NewStock(mustProductID(t, "product-1"), 5)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := s.Reserve(3); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		// 残り Available は 2（5 - 3）なので、3 個の追加引当は拒否されるべき。
		err = s.Reserve(3)
		if err == nil {
			t.Fatal("expected error for insufficient stock, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Errorf("error = %v, want DomainRuleError", err)
		}
		if s.Reserved() != 3 {
			t.Errorf("Reserved() after rejected reserve = %d, want 3 (unchanged)", s.Reserved())
		}
	})

	t.Run("reserving zero is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		err := s.Reserve(0)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("Reserve(0) error = %v, want DomainRuleError", err)
		}
	})

	t.Run("reserving a negative amount is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		err := s.Reserve(-1)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("Reserve(-1) error = %v, want DomainRuleError", err)
		}
	})
}

// TestStock_Release は引当解除（Release）の境界値を検証する。
func TestStock_Release(t *testing.T) {
	t.Run("releasing within reserved amount succeeds", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(5); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		if err := s.Release(2); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Reserved() != 3 {
			t.Errorf("Reserved() = %d, want 3", s.Reserved())
		}
		if s.Quantity() != 10 {
			t.Errorf("Quantity() = %d, want 10 (unchanged)", s.Quantity())
		}
	})

	t.Run("releasing more than reserved is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(3); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		err := s.Release(4)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("Release(4) error = %v, want DomainRuleError", err)
		}
		if s.Reserved() != 3 {
			t.Errorf("Reserved() after rejected release = %d, want 3 (unchanged)", s.Reserved())
		}
	})

	t.Run("releasing zero is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(3); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		err := s.Release(0)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("Release(0) error = %v, want DomainRuleError", err)
		}
	})
}

// TestStock_ConsumeReserved は出荷時の消込（ConsumeReserved）が
// quantity と reserved の両方を正しく減らすことを検証する。
func TestStock_ConsumeReserved(t *testing.T) {
	t.Run("consuming within reserved amount decrements both quantity and reserved", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(6); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		if err := s.ConsumeReserved(4); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Quantity() != 6 {
			t.Errorf("Quantity() = %d, want 6", s.Quantity())
		}
		if s.Reserved() != 2 {
			t.Errorf("Reserved() = %d, want 2", s.Reserved())
		}
		if s.Available() != 4 {
			t.Errorf("Available() = %d, want 4", s.Available())
		}
	})

	t.Run("consuming more than reserved is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(3); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		err := s.ConsumeReserved(4)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("ConsumeReserved(4) error = %v, want DomainRuleError", err)
		}
		if s.Quantity() != 10 || s.Reserved() != 3 {
			t.Errorf("state changed after rejected consume: quantity=%d reserved=%d, want 10/3", s.Quantity(), s.Reserved())
		}
	})

	t.Run("consuming zero is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(3); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		err := s.ConsumeReserved(0)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("ConsumeReserved(0) error = %v, want DomainRuleError", err)
		}
	})
}

// TestStock_SetQuantity は実在庫数の更新に関する不変条件を検証する。
func TestStock_SetQuantity(t *testing.T) {
	t.Run("setting quantity above reserved succeeds", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(4); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		if err := s.SetQuantity(20); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Quantity() != 20 {
			t.Errorf("Quantity() = %d, want 20", s.Quantity())
		}
		if s.Reserved() != 4 {
			t.Errorf("Reserved() = %d, want 4 (unchanged)", s.Reserved())
		}
	})

	t.Run("setting quantity equal to reserved succeeds", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(4); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		if err := s.SetQuantity(4); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("setting quantity below reserved is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		if err := s.Reserve(4); err != nil {
			t.Fatalf("setup reserve failed: %v", err)
		}
		err := s.SetQuantity(3)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("SetQuantity(3) error = %v, want DomainRuleError", err)
		}
		if s.Quantity() != 10 {
			t.Errorf("Quantity() after rejected update = %d, want 10 (unchanged)", s.Quantity())
		}
	})

	t.Run("setting a negative quantity is rejected", func(t *testing.T) {
		s, _ := inventory.NewStock(mustProductID(t, "product-1"), 10)
		err := s.SetQuantity(-1)
		if err == nil || !shared.IsDomainRuleError(err) {
			t.Fatalf("SetQuantity(-1) error = %v, want DomainRuleError", err)
		}
	})
}

// TestReconstructStock は永続化層からの復元が検証をスキップして
// そのまま状態を再現することを確認する。version も DB から読み出した値を
// そのまま引き継ぐことを検証する（楽観ロックの前提となる挙動）。
func TestReconstructStock(t *testing.T) {
	productID := mustProductID(t, "product-1")
	s := inventory.ReconstructStock(productID, 10, 4, 3)
	if s.ProductID() != productID {
		t.Errorf("ProductID() = %v, want %v", s.ProductID(), productID)
	}
	if s.Quantity() != 10 {
		t.Errorf("Quantity() = %d, want 10", s.Quantity())
	}
	if s.Reserved() != 4 {
		t.Errorf("Reserved() = %d, want 4", s.Reserved())
	}
	if s.Available() != 6 {
		t.Errorf("Available() = %d, want 6", s.Available())
	}
	if s.Version() != 3 {
		t.Errorf("Version() = %d, want 3", s.Version())
	}
}
