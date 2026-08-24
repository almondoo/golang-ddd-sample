package catalog_test

import (
	"strings"
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

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

// TestNewProduct は「商品登録」というユースケースの入口にあたる
// コンストラクタが、業務ルール（不変条件）をどう守るかをテーブル駆動で検証する。
func TestNewProduct(t *testing.T) {
	price := mustMoney(t, 1000)

	tests := []struct {
		name        string
		productName string
		description string
		price       shared.Money
		wantErr     bool
		wantRuleErr bool
	}{
		{
			name:        "valid product is created",
			productName: "サンプル商品",
			description: "説明文",
			price:       price,
			wantErr:     false,
		},
		{
			name:        "empty name is rejected",
			productName: "",
			description: "説明文",
			price:       price,
			wantErr:     true,
			wantRuleErr: true,
		},
		{
			name:        "name longer than 100 runes is rejected",
			productName: strings.Repeat("あ", 101),
			description: "説明文",
			price:       price,
			wantErr:     true,
			wantRuleErr: true,
		},
		{
			name:        "name of exactly 100 runes is accepted",
			productName: strings.Repeat("あ", 100),
			description: "説明文",
			price:       price,
			wantErr:     false,
		},
		{
			name:        "description longer than 1000 runes is rejected",
			productName: "サンプル商品",
			description: strings.Repeat("あ", 1001),
			price:       price,
			wantErr:     true,
			wantRuleErr: true,
		},
		{
			name:        "zero price is accepted (not negative)",
			productName: "サンプル商品",
			description: "説明文",
			price:       mustMoney(t, 0),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := catalog.NewProduct(tt.productName, tt.description, tt.price)
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
			if p.ID() == "" {
				t.Fatalf("expected generated ID, got empty")
			}
			if p.Name() != tt.productName {
				t.Fatalf("Name() = %q, want %q", p.Name(), tt.productName)
			}
			if p.Description() != tt.description {
				t.Fatalf("Description() = %q, want %q", p.Description(), tt.description)
			}
			if !p.Price().Equals(tt.price) {
				t.Fatalf("Price() = %v, want %v", p.Price(), tt.price)
			}
		})
	}
}

// TestProduct_ChangePrice は「価格変更」という業務操作のルールを検証する。
func TestProduct_ChangePrice(t *testing.T) {
	t.Run("changing to a different price succeeds", func(t *testing.T) {
		p, err := catalog.NewProduct("商品", "説明", mustMoney(t, 1000))
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		newPrice := mustMoney(t, 2000)
		if err := p.ChangePrice(newPrice); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !p.Price().Equals(newPrice) {
			t.Fatalf("Price() = %v, want %v", p.Price(), newPrice)
		}
	})

	t.Run("changing to the same price is rejected as a domain rule violation", func(t *testing.T) {
		price := mustMoney(t, 1000)
		p, err := catalog.NewProduct("商品", "説明", price)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		err = p.ChangePrice(price)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		if !p.Price().Equals(price) {
			t.Fatalf("price should remain unchanged after rejected update, got %v", p.Price())
		}
	})
}
