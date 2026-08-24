package cart

import (
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// newTestCart はテスト用に空のカートを 1 つ用意するヘルパーである。
func newTestCart(t *testing.T) *Cart {
	t.Helper()
	customerID, err := NewCustomerID("customer-1")
	if err != nil {
		t.Fatalf("NewCustomerID() unexpected error: %v", err)
	}
	return NewCart(customerID)
}

// TestCart_AddItem_AddsNewLine は、新規商品の追加が明細として 1 件増える
// ことを確認する（マージが発生しない基本ケース）。
func TestCart_AddItem_AddsNewLine(t *testing.T) {
	c := newTestCart(t)
	productID, err := NewProductID("product-1")
	if err != nil {
		t.Fatalf("NewProductID() unexpected error: %v", err)
	}

	if err := c.AddItem(productID, 3); err != nil {
		t.Fatalf("AddItem() unexpected error: %v", err)
	}

	items := c.Items()
	if len(items) != 1 {
		t.Fatalf("len(Items()) = %d, want 1", len(items))
	}
	if items[0].ProductID() != productID {
		t.Errorf("ProductID() = %v, want %v", items[0].ProductID(), productID)
	}
	if items[0].Quantity() != 3 {
		t.Errorf("Quantity() = %d, want 3", items[0].Quantity())
	}
}

// TestCart_AddItem_MergesSameProduct は、同一商品を複数回追加した場合に
// 明細が増えるのではなく数量がマージされることを確認する。
func TestCart_AddItem_MergesSameProduct(t *testing.T) {
	c := newTestCart(t)
	productID, _ := NewProductID("product-1")

	if err := c.AddItem(productID, 3); err != nil {
		t.Fatalf("AddItem() 1st call unexpected error: %v", err)
	}
	if err := c.AddItem(productID, 5); err != nil {
		t.Fatalf("AddItem() 2nd call unexpected error: %v", err)
	}

	items := c.Items()
	if len(items) != 1 {
		t.Fatalf("len(Items()) = %d, want 1 (merged)", len(items))
	}
	if items[0].Quantity() != 8 {
		t.Errorf("Quantity() = %d, want 8", items[0].Quantity())
	}
}

// TestCart_AddItem_QuantityBounds は、数量の下限・上限に関する不変条件を検証する。
func TestCart_AddItem_QuantityBounds(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		wantErr  bool
	}{
		{name: "below minimum", quantity: 0, wantErr: true},
		{name: "negative", quantity: -1, wantErr: true},
		{name: "minimum allowed", quantity: 1, wantErr: false},
		{name: "maximum allowed", quantity: 99, wantErr: false},
		{name: "above maximum", quantity: 100, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCart(t)
			productID, _ := NewProductID("product-1")

			err := c.AddItem(productID, tt.quantity)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("AddItem(%d) expected error, got nil", tt.quantity)
				}
				if !shared.IsDomainRuleError(err) {
					t.Errorf("AddItem(%d) error = %v, want DomainRuleError", tt.quantity, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("AddItem(%d) unexpected error: %v", tt.quantity, err)
			}
		})
	}
}

// TestCart_AddItem_MergedQuantityExceedsCeiling は、個々の追加は上限内でも
// マージ後の合計が上限を超える場合に拒否されることを確認する。
func TestCart_AddItem_MergedQuantityExceedsCeiling(t *testing.T) {
	c := newTestCart(t)
	productID, _ := NewProductID("product-1")

	if err := c.AddItem(productID, 60); err != nil {
		t.Fatalf("AddItem() 1st call unexpected error: %v", err)
	}
	err := c.AddItem(productID, 60)
	if err == nil {
		t.Fatal("AddItem() expected error for merged quantity exceeding ceiling, got nil")
	}
	if !shared.IsDomainRuleError(err) {
		t.Errorf("AddItem() error = %v, want DomainRuleError", err)
	}

	// マージが拒否された場合、元の数量は変更されずに保たれているべきである。
	items := c.Items()
	if items[0].Quantity() != 60 {
		t.Errorf("Quantity() after rejected merge = %d, want 60 (unchanged)", items[0].Quantity())
	}
}

// TestCart_AddItem_DistinctItemCeiling は、明細数（商品種別数）の上限を検証する。
func TestCart_AddItem_DistinctItemCeiling(t *testing.T) {
	c := newTestCart(t)

	for i := 0; i < maxDistinctItems; i++ {
		productID, _ := NewProductID(shared.NewID())
		if err := c.AddItem(productID, 1); err != nil {
			t.Fatalf("AddItem() call %d unexpected error: %v", i, err)
		}
	}

	// 21 件目（新しい商品）は上限超過として拒否されるべきである。
	extraProductID, _ := NewProductID(shared.NewID())
	err := c.AddItem(extraProductID, 1)
	if err == nil {
		t.Fatal("AddItem() expected error for exceeding distinct item ceiling, got nil")
	}
	if !shared.IsDomainRuleError(err) {
		t.Errorf("AddItem() error = %v, want DomainRuleError", err)
	}
	if len(c.Items()) != maxDistinctItems {
		t.Errorf("len(Items()) = %d, want %d (unchanged)", len(c.Items()), maxDistinctItems)
	}
}

// TestCart_RemoveItem_RemovesExistingLine は、存在する明細の削除が
// 正しく機能することを確認する。
func TestCart_RemoveItem_RemovesExistingLine(t *testing.T) {
	c := newTestCart(t)
	productID, _ := NewProductID("product-1")
	if err := c.AddItem(productID, 1); err != nil {
		t.Fatalf("AddItem() unexpected error: %v", err)
	}

	if err := c.RemoveItem(productID); err != nil {
		t.Fatalf("RemoveItem() unexpected error: %v", err)
	}
	if !c.IsEmpty() {
		t.Errorf("IsEmpty() = false, want true after removing the only item")
	}
}

// TestCart_RemoveItem_MissingProduct は、カートに存在しない商品を
// 削除しようとした場合にドメインルール違反となることを確認する。
func TestCart_RemoveItem_MissingProduct(t *testing.T) {
	c := newTestCart(t)
	productID, _ := NewProductID("product-not-in-cart")

	err := c.RemoveItem(productID)
	if err == nil {
		t.Fatal("RemoveItem() expected error for missing product, got nil")
	}
	if !shared.IsDomainRuleError(err) {
		t.Errorf("RemoveItem() error = %v, want DomainRuleError", err)
	}
}

// TestCart_Clear_And_IsEmpty は、Clear が全明細を取り除き
// IsEmpty がその状態を正しく反映することを確認する。
func TestCart_Clear_And_IsEmpty(t *testing.T) {
	c := newTestCart(t)
	if !c.IsEmpty() {
		t.Fatalf("IsEmpty() = false, want true for a fresh cart")
	}

	productID, _ := NewProductID("product-1")
	if err := c.AddItem(productID, 2); err != nil {
		t.Fatalf("AddItem() unexpected error: %v", err)
	}
	if c.IsEmpty() {
		t.Fatalf("IsEmpty() = true, want false after adding an item")
	}

	c.Clear()
	if !c.IsEmpty() {
		t.Errorf("IsEmpty() = false, want true after Clear()")
	}
	if len(c.Items()) != 0 {
		t.Errorf("len(Items()) = %d, want 0 after Clear()", len(c.Items()))
	}
}
