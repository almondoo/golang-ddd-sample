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

func TestMoney_Subtract(t *testing.T) {
	a, err := NewMoney(300, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := NewMoney(100, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	diff, err := a.Subtract(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff.Amount() != 200 {
		t.Fatalf("expected amount 200, got %d", diff.Amount())
	}
}

func TestMoney_Subtract_CurrencyMismatch(t *testing.T) {
	jpy, err := NewMoney(100, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	usd, err := NewMoney(100, Currency("USD"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := jpy.Subtract(usd); err == nil {
		t.Fatal("expected error for currency mismatch, got nil")
	}
}

func TestMoney_Subtract_NegativeResultIsRejected(t *testing.T) {
	a, err := NewMoney(100, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := NewMoney(200, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := a.Subtract(b); err == nil {
		t.Fatal("expected error for negative result, got nil")
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

func TestNewMoney_MaxAmountIsAccepted(t *testing.T) {
	if _, err := NewMoney(maxAmount, JPY); err != nil {
		t.Fatalf("unexpected error at maxAmount: %v", err)
	}
}

func TestNewMoney_ExceedingMaxAmountIsRejected(t *testing.T) {
	if _, err := NewMoney(maxAmount+1, JPY); err == nil {
		t.Fatal("expected error for amount exceeding maxAmount, got nil")
	}
}

func TestMoney_Multiply_ExceedingMaxAmountIsRejected(t *testing.T) {
	// maxAmount(1兆円) × 99 = 9.9e13 は int64 の範囲には収まる
	// （オーバーフローはしない）ため、Multiply 内のオーバーフロー検査
	// （dead-man protection、公開 API 経由では到達しない）自体は発火しない。
	// しかしその計算結果は maxAmount を超えるため、Multiply の最後に呼ばれる
	// NewMoney の上限チェックによって正しく拒否される。この経路は
	// 「単価 × 数量」が現実的にあり得ない金額になるケースを、乗算そのものが
	// 安全でも上限チェックで確実に弾けることを示している。
	m, err := NewMoney(maxAmount, JPY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := m.Multiply(99); err == nil {
		t.Fatal("expected error for maxAmount * 99 exceeding maxAmount, got nil")
	}
}

func TestIsDomainRuleError(t *testing.T) {
	_, err := NewMoney(-1, JPY)
	if !IsDomainRuleError(err) {
		t.Fatal("expected IsDomainRuleError to be true for negative amount error")
	}
}
