package customer_test

import (
	"strings"
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/domain/customer"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// mustAddAddress はテスト用のヘルパーで、AddAddress の生成失敗をテスト失敗として扱う。
func mustAddAddress(t *testing.T, c *customer.Customer, postalCode, prefecture, city, line string) customer.AddressID {
	t.Helper()
	id, err := c.AddAddress(postalCode, prefecture, city, line)
	if err != nil {
		t.Fatalf("failed to add address fixture: %v", err)
	}
	return id
}

// TestNewCustomer は「顧客登録」というユースケースの入口にあたる
// コンストラクタが、業務ルール（不変条件）をどう守るかをテーブル駆動で検証する。
func TestNewCustomer(t *testing.T) {
	tests := []struct {
		name        string
		customer    string
		email       string
		wantErr     bool
		wantRuleErr bool
	}{
		{
			name:     "valid customer is created",
			customer: "山田太郎",
			email:    "taro@example.com",
			wantErr:  false,
		},
		{
			name:        "empty name is rejected",
			customer:    "",
			email:       "taro@example.com",
			wantErr:     true,
			wantRuleErr: true,
		},
		{
			name:        "name longer than 50 runes is rejected",
			customer:    strings.Repeat("あ", 51),
			email:       "taro@example.com",
			wantErr:     true,
			wantRuleErr: true,
		},
		{
			name:     "name of exactly 50 runes is accepted",
			customer: strings.Repeat("あ", 50),
			email:    "taro@example.com",
			wantErr:  false,
		},
		{
			name:        "empty email is rejected",
			customer:    "山田太郎",
			email:       "",
			wantErr:     true,
			wantRuleErr: true,
		},
		{
			name:        "email without @ is rejected",
			customer:    "山田太郎",
			email:       "taro.example.com",
			wantErr:     true,
			wantRuleErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := customer.NewCustomer(tt.customer, tt.email)
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
			if c.ID() == "" {
				t.Fatalf("expected generated ID, got empty")
			}
			if c.Name() != tt.customer {
				t.Fatalf("Name() = %q, want %q", c.Name(), tt.customer)
			}
			if c.Email() != tt.email {
				t.Fatalf("Email() = %q, want %q", c.Email(), tt.email)
			}
			if len(c.Addresses()) != 0 {
				t.Fatalf("expected no addresses on registration, got %d", len(c.Addresses()))
			}
		})
	}
}

// TestCustomer_AddAddress_FirstAddressBecomesDefault は、集約ルートが
// 「デフォルト住所は住所があれば必ず1つ」という不変条件を、子エンティティ
// （Address）に頼らずどう守るかを検証する。
func TestCustomer_AddAddress_FirstAddressBecomesDefault(t *testing.T) {
	c, err := customer.NewCustomer("山田太郎", "taro@example.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	firstID := mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")

	addrs := c.Addresses()
	if len(addrs) != 1 {
		t.Fatalf("len(Addresses()) = %d, want 1", len(addrs))
	}
	if addrs[0].ID() != firstID {
		t.Fatalf("Addresses()[0].ID() = %v, want %v", addrs[0].ID(), firstID)
	}
	if !addrs[0].IsDefault() {
		t.Fatalf("first address should automatically become default")
	}

	def, err := c.DefaultAddress()
	if err != nil {
		t.Fatalf("DefaultAddress() unexpected error: %v", err)
	}
	if def.ID() != firstID {
		t.Fatalf("DefaultAddress().ID() = %v, want %v", def.ID(), firstID)
	}

	// 2件目を追加してもデフォルトは変わらない。
	secondID := mustAddAddress(t, c, "530-0001", "大阪府", "大阪市", "2-2-2")
	def, err = c.DefaultAddress()
	if err != nil {
		t.Fatalf("DefaultAddress() unexpected error: %v", err)
	}
	if def.ID() != firstID {
		t.Fatalf("adding a second address should not change the default; DefaultAddress().ID() = %v, want %v", def.ID(), firstID)
	}
	if secondID == firstID {
		t.Fatalf("expected distinct address IDs")
	}
}

// TestCustomer_AddAddress_ValidationAndLimit は必須項目チェックと
// 住所件数の上限（5件）を検証する。
func TestCustomer_AddAddress_ValidationAndLimit(t *testing.T) {
	t.Run("each field is required", func(t *testing.T) {
		tests := []struct {
			name       string
			postalCode string
			prefecture string
			city       string
			line       string
		}{
			{name: "empty postal code", postalCode: "", prefecture: "東京都", city: "千代田区", line: "1-1-1"},
			{name: "empty prefecture", postalCode: "100-0001", prefecture: "", city: "千代田区", line: "1-1-1"},
			{name: "empty city", postalCode: "100-0001", prefecture: "東京都", city: "", line: "1-1-1"},
			{name: "empty line", postalCode: "100-0001", prefecture: "東京都", city: "千代田区", line: ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c, err := customer.NewCustomer("山田太郎", "taro@example.com")
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				_, err = c.AddAddress(tt.postalCode, tt.prefecture, tt.city, tt.line)
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !shared.IsDomainRuleError(err) {
					t.Fatalf("expected domain rule error, got %v", err)
				}
			})
		}
	})

	t.Run("6th address is rejected as exceeding the ceiling", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		for i := 0; i < 5; i++ {
			mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")
		}
		if len(c.Addresses()) != 5 {
			t.Fatalf("len(Addresses()) = %d, want 5", len(c.Addresses()))
		}

		_, err = c.AddAddress("100-0001", "東京都", "千代田区", "1-1-1")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		if len(c.Addresses()) != 5 {
			t.Fatalf("address count should remain 5 after rejected addition, got %d", len(c.Addresses()))
		}
	})
}

// TestCustomer_ChangeDefaultAddress は「デフォルトはちょうど1つ」という
// 不変条件が、変更操作の前後で常に成立し続けることを検証する。
func TestCustomer_ChangeDefaultAddress(t *testing.T) {
	t.Run("changing default flips old default off and new one on", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		firstID := mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")
		secondID := mustAddAddress(t, c, "530-0001", "大阪府", "大阪市", "2-2-2")

		if err := c.ChangeDefaultAddress(secondID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var defaultCount int
		for _, a := range c.Addresses() {
			if a.IsDefault() {
				defaultCount++
			}
			if a.ID() == firstID && a.IsDefault() {
				t.Fatalf("old default address should no longer be default")
			}
			if a.ID() == secondID && !a.IsDefault() {
				t.Fatalf("new default address should be default")
			}
		}
		if defaultCount != 1 {
			t.Fatalf("exactly one address should be default, got %d", defaultCount)
		}
	})

	t.Run("unknown address id is rejected", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")

		unknownID, err := customer.NewAddressID("unknown-address")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		err = c.ChangeDefaultAddress(unknownID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}

// TestCustomer_RemoveAddress は「デフォルト住所は他に住所が残る場合は
// 削除できないが、最後の1件は削除できる」という設計判断を検証する。
func TestCustomer_RemoveAddress(t *testing.T) {
	t.Run("removing the default address while others remain is rejected", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		firstID := mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")
		mustAddAddress(t, c, "530-0001", "大阪府", "大阪市", "2-2-2")

		err = c.RemoveAddress(firstID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		if len(c.Addresses()) != 2 {
			t.Fatalf("address count should remain 2 after rejected removal, got %d", len(c.Addresses()))
		}
	})

	t.Run("removing a non-default address while others remain succeeds", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		firstID := mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")
		secondID := mustAddAddress(t, c, "530-0001", "大阪府", "大阪市", "2-2-2")

		if err := c.RemoveAddress(secondID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		addrs := c.Addresses()
		if len(addrs) != 1 {
			t.Fatalf("len(Addresses()) = %d, want 1", len(addrs))
		}
		if addrs[0].ID() != firstID {
			t.Fatalf("remaining address ID = %v, want %v", addrs[0].ID(), firstID)
		}
	})

	t.Run("removing the last remaining address is allowed", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		onlyID := mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")

		if err := c.RemoveAddress(onlyID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(c.Addresses()) != 0 {
			t.Fatalf("expected no addresses remaining, got %d", len(c.Addresses()))
		}
	})

	t.Run("unknown address id is rejected", func(t *testing.T) {
		c, err := customer.NewCustomer("山田太郎", "taro@example.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		mustAddAddress(t, c, "100-0001", "東京都", "千代田区", "1-1-1")

		unknownID, err := customer.NewAddressID("unknown-address")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		err = c.RemoveAddress(unknownID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}

// TestCustomer_DefaultAddress_NoAddresses は、住所が1件もない顧客に
// DefaultAddress を呼んだ場合の挙動を検証する。
func TestCustomer_DefaultAddress_NoAddresses(t *testing.T) {
	c, err := customer.NewCustomer("山田太郎", "taro@example.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = c.DefaultAddress()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !shared.IsDomainRuleError(err) {
		t.Fatalf("expected domain rule error, got %v", err)
	}
}
