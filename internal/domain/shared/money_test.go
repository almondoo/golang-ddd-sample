package shared

import "testing"

// Money 値オブジェクトの基本的な不変条件・振る舞いを検証するテスト。
// 網羅的なテストは test-verifier の責務だが、契約境界の最低限の動作は
// ここで確認しておく。

func TestNewMoney_NegativeAmountIsRejected(t *testing.T) {
	if _, err := NewMoney(-1, JPY); err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
}

func TestMoney_Add(t *testing.T) {
	a, err := NewMoney(100, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := NewMoney(200, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.Amount() != 300 {
		t.Fatalf("expected amount 300, got %d", sum.Amount())
	}
}

func TestMoney_Add_CurrencyMismatch(t *testing.T) {
	jpy, err := NewMoney(100, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	usd, err := NewMoney(100, Currency("USD"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := jpy.Add(usd); err == nil {
		t.Fatal("expected error for currency mismatch, got nil")
	}
}

func TestMoney_Multiply(t *testing.T) {
	m, err := NewMoney(100, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := m.Multiply(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Amount() != 300 {
		t.Fatalf("expected amount 300, got %d", result.Amount())
	}

	if _, err := m.Multiply(-1); err == nil {
		t.Fatal("expected error for negative multiplier, got nil")
	}
}

func TestMoney_IsZero(t *testing.T) {
	zero, err := NewMoney(0, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !zero.IsZero() {
		t.Fatal("expected IsZero to be true")
	}

	nonZero, err := NewMoney(1, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nonZero.IsZero() {
		t.Fatal("expected IsZero to be false")
	}
}

func TestMoney_Equals(t *testing.T) {
	a, _ := NewMoney(100, JPY)
	b, _ := NewMoney(100, JPY)
	c, _ := NewMoney(200, JPY)

	if !a.Equals(b) {
		t.Fatal("expected a to equal b")
	}
	if a.Equals(c) {
		t.Fatal("expected a to not equal c")
	}
}

func TestIsDomainRuleError(t *testing.T) {
	_, err := NewMoney(-1, JPY)
	if !IsDomainRuleError(err) {
		t.Fatal("expected IsDomainRuleError to be true for negative amount error")
	}
}
