package coupon

import (
	"testing"
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// Coupon 集約・値オブジェクトの基本的な不変条件・振る舞いを検証するテスト。
// 網羅的なテストは test-verifier の責務だが、契約境界の最低限の動作は
// ここで確認しておく。

// TestNewCouponCode_FormatMatrix はクーポンコードの形式ルール
// （4〜20文字、大文字英数字とハイフンのみ）を表形式で検証する。
func TestNewCouponCode_FormatMatrix(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "minimum length (4 chars)", code: "ABCD", wantErr: false},
		{name: "maximum length (20 chars)", code: "ABCDEFGHIJKLMNOPQRST", wantErr: false},
		{name: "digits allowed", code: "AB12", wantErr: false},
		{name: "hyphen allowed", code: "AB-CD", wantErr: false},
		{name: "too short (3 chars)", code: "ABC", wantErr: true},
		{name: "too long (21 chars)", code: "ABCDEFGHIJKLMNOPQRSTU", wantErr: true},
		{name: "lowercase not allowed", code: "abcd", wantErr: true},
		{name: "underscore not allowed", code: "AB_CD", wantErr: true},
		{name: "empty string", code: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := NewCouponCode(tt.code)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewCouponCode(%q) expected error, got nil", tt.code)
				}
				if !shared.IsDomainRuleError(err) {
					t.Errorf("NewCouponCode(%q) error = %v, want DomainRuleError", tt.code, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewCouponCode(%q) unexpected error: %v", tt.code, err)
			}
			if code.String() != tt.code {
				t.Errorf("String() = %q, want %q", code.String(), tt.code)
			}
		})
	}
}

// newTestCode はテスト用に妥当なクーポンコードを 1 つ用意するヘルパーである。
func newTestCode(t *testing.T) CouponCode {
	t.Helper()
	code, err := NewCouponCode("SUMMER10")
	if err != nil {
		t.Fatalf("NewCouponCode() unexpected error: %v", err)
	}
	return code
}

// TestNewAmountCoupon_Validations は amount 型クーポンのコンストラクタが
// 守る不変条件（金額はゼロ不可、利用回数上限は 1 以上）を検証する。
func TestNewAmountCoupon_Validations(t *testing.T) {
	code := newTestCode(t)
	future := time.Now().Add(24 * time.Hour)

	t.Run("valid amount coupon", func(t *testing.T) {
		amount, _ := shared.NewMoney(500, shared.JPY)
		c, err := NewAmountCoupon(code, amount, future, 1)
		if err != nil {
			t.Fatalf("NewAmountCoupon() unexpected error: %v", err)
		}
		if c.Type() != DiscountTypeAmount {
			t.Errorf("Type() = %v, want %v", c.Type(), DiscountTypeAmount)
		}
		if c.Amount().Amount() != 500 {
			t.Errorf("Amount().Amount() = %d, want 500", c.Amount().Amount())
		}
	})

	t.Run("zero amount is rejected", func(t *testing.T) {
		zero, _ := shared.NewMoney(0, shared.JPY)
		_, err := NewAmountCoupon(code, zero, future, 1)
		if err == nil {
			t.Fatal("NewAmountCoupon() expected error for zero amount, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Errorf("NewAmountCoupon() error = %v, want DomainRuleError", err)
		}
	})

	t.Run("usage limit below 1 is rejected", func(t *testing.T) {
		amount, _ := shared.NewMoney(500, shared.JPY)
		_, err := NewAmountCoupon(code, amount, future, 0)
		if err == nil {
			t.Fatal("NewAmountCoupon() expected error for usage limit 0, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Errorf("NewAmountCoupon() error = %v, want DomainRuleError", err)
		}
	})
}

// TestNewRateCoupon_Validations は rate 型クーポンのコンストラクタが
// 守る不変条件（割合は 1〜100、利用回数上限は 1 以上）を検証する。
func TestNewRateCoupon_Validations(t *testing.T) {
	code := newTestCode(t)
	future := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		ratePercent int
		usageLimit  int
		wantErr     bool
	}{
		{name: "minimum rate (1%)", ratePercent: 1, usageLimit: 1, wantErr: false},
		{name: "maximum rate (100%)", ratePercent: 100, usageLimit: 1, wantErr: false},
		{name: "rate below minimum (0%)", ratePercent: 0, usageLimit: 1, wantErr: true},
		{name: "rate above maximum (101%)", ratePercent: 101, usageLimit: 1, wantErr: true},
		{name: "usage limit below 1", ratePercent: 10, usageLimit: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewRateCoupon(code, tt.ratePercent, future, tt.usageLimit)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewRateCoupon() expected error, got nil")
				}
				if !shared.IsDomainRuleError(err) {
					t.Errorf("NewRateCoupon() error = %v, want DomainRuleError", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewRateCoupon() unexpected error: %v", err)
			}
			if c.RatePercent() != tt.ratePercent {
				t.Errorf("RatePercent() = %d, want %d", c.RatePercent(), tt.ratePercent)
			}
		})
	}
}

// TestCoupon_Use はクーポン消費（Use）の 3 パターン
// （期限切れ／利用回数上限到達／正常消費）を検証する。
func TestCoupon_Use(t *testing.T) {
	code := newTestCode(t)
	amount, _ := shared.NewMoney(500, shared.JPY)
	now := time.Now()

	t.Run("expired coupon is rejected", func(t *testing.T) {
		past := now.Add(-1 * time.Hour)
		c, err := NewAmountCoupon(code, amount, past, 5)
		if err != nil {
			t.Fatalf("NewAmountCoupon() unexpected error: %v", err)
		}

		if err := c.Use(now); err == nil {
			t.Fatal("Use() expected error for expired coupon, got nil")
		} else if !shared.IsDomainRuleError(err) {
			t.Errorf("Use() error = %v, want DomainRuleError", err)
		}
	})

	t.Run("usage limit reached is rejected", func(t *testing.T) {
		future := now.Add(24 * time.Hour)
		c, err := NewAmountCoupon(code, amount, future, 1)
		if err != nil {
			t.Fatalf("NewAmountCoupon() unexpected error: %v", err)
		}

		if err := c.Use(now); err != nil {
			t.Fatalf("Use() 1st call unexpected error: %v", err)
		}
		if err := c.Use(now); err == nil {
			t.Fatal("Use() 2nd call expected error (limit reached), got nil")
		} else if !shared.IsDomainRuleError(err) {
			t.Errorf("Use() error = %v, want DomainRuleError", err)
		}
	})

	t.Run("successful use increments used count", func(t *testing.T) {
		future := now.Add(24 * time.Hour)
		c, err := NewAmountCoupon(code, amount, future, 3)
		if err != nil {
			t.Fatalf("NewAmountCoupon() unexpected error: %v", err)
		}

		if err := c.Use(now); err != nil {
			t.Fatalf("Use() unexpected error: %v", err)
		}
		if c.UsedCount() != 1 {
			t.Errorf("UsedCount() = %d, want 1", c.UsedCount())
		}
	})
}

// TestCoupon_DiscountFor_Amount は amount 型クーポンの割引計算が
// 「割引額は合計金額を超えない」というルールを守ることを検証する。
func TestCoupon_DiscountFor_Amount(t *testing.T) {
	code := newTestCode(t)
	future := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name         string
		couponAmount int64
		total        int64
		wantDiscount int64
	}{
		{name: "amount less than total", couponAmount: 300, total: 1000, wantDiscount: 300},
		{name: "amount exceeds total is capped", couponAmount: 500, total: 300, wantDiscount: 300},
		{name: "amount equals total", couponAmount: 1000, total: 1000, wantDiscount: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, _ := shared.NewMoney(tt.couponAmount, shared.JPY)
			c, err := NewAmountCoupon(code, amount, future, 1)
			if err != nil {
				t.Fatalf("NewAmountCoupon() unexpected error: %v", err)
			}

			total, _ := shared.NewMoney(tt.total, shared.JPY)
			discount, err := c.DiscountFor(total)
			if err != nil {
				t.Fatalf("DiscountFor() unexpected error: %v", err)
			}
			if discount.Amount() != tt.wantDiscount {
				t.Errorf("DiscountFor().Amount() = %d, want %d", discount.Amount(), tt.wantDiscount)
			}
		})
	}
}

// TestCoupon_DiscountFor_Rate は rate 型クーポンの割引計算が
// int64 の切り捨て除算により端数を切り捨てることを検証する。
func TestCoupon_DiscountFor_Rate(t *testing.T) {
	code := newTestCode(t)
	future := time.Now().Add(24 * time.Hour)

	// 10% of 101 = 10.1 → 10（1円未満切り捨て）となることを確認する。
	c, err := NewRateCoupon(code, 10, future, 1)
	if err != nil {
		t.Fatalf("NewRateCoupon() unexpected error: %v", err)
	}

	total, _ := shared.NewMoney(101, shared.JPY)
	discount, err := c.DiscountFor(total)
	if err != nil {
		t.Fatalf("DiscountFor() unexpected error: %v", err)
	}
	if discount.Amount() != 10 {
		t.Errorf("DiscountFor().Amount() = %d, want 10 (truncated)", discount.Amount())
	}
}
