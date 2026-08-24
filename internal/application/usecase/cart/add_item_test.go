package cart_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/application/usecase/cart"
	domaincart "github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// このファイルは AddItemUseCase の「find-or-create」分岐(既存カートへの
// マージ / 新規カートの暗黙的な生成)を、手書きのフェイクだけで検証する。
// place_order_test.go と同じ理由でモックライブラリは使わない。

// ---- fakeCartRepository ----

// fakeCartRepository は domaincart.Repository の手書きフェイクである。
// findErr を差し替えることで「ErrNotFound 以外のエラー」も注入できる
// ようにしている点が、find-or-create 分岐のテストに必要な唯一の工夫である。
type fakeCartRepository struct {
	carts     map[domaincart.CustomerID]*domaincart.Cart
	findErr   error
	saveErr   error
	saved     []*domaincart.Cart
	saveCalls int
}

func newFakeCartRepository() *fakeCartRepository {
	return &fakeCartRepository{carts: map[domaincart.CustomerID]*domaincart.Cart{}}
}

func (r *fakeCartRepository) FindByCustomerID(_ context.Context, id domaincart.CustomerID) (*domaincart.Cart, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	c, ok := r.carts[id]
	if !ok {
		return nil, fmt.Errorf("cart not found: customer=%s: %w", id, shared.ErrNotFound)
	}
	return c, nil
}

func (r *fakeCartRepository) Save(_ context.Context, c *domaincart.Cart) error {
	r.saveCalls++
	if r.saveErr != nil {
		return r.saveErr
	}
	r.carts[c.CustomerID()] = c
	r.saved = append(r.saved, c)
	return nil
}

// ---- fakeTxManager ----

// fakeTxManager は fn(ctx) を素通しするだけの tx.Manager フェイクである。
// 本物のロールバック検証はしない — トランザクション境界が1回で呼ばれる
// ことだけを確認する(place_order_test.go の fakeTxManager と同じ方針)。
type fakeTxManager struct {
	calls int
}

func (m *fakeTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	m.calls++
	return fn(ctx)
}

var _ tx.Manager = (*fakeTxManager)(nil)

const testCustomerID = "customer-1"

func TestAddItemUseCase_Execute(t *testing.T) {
	t.Run("adding to an existing cart merges the quantity", func(t *testing.T) {
		repo := newFakeCartRepository()
		txManager := &fakeTxManager{}
		customerID, err := domaincart.NewCustomerID(testCustomerID)
		if err != nil {
			t.Fatalf("failed to build customer id fixture: %v", err)
		}
		productID, err := domaincart.NewProductID("product-1")
		if err != nil {
			t.Fatalf("failed to build product id fixture: %v", err)
		}

		existing := domaincart.NewCart(customerID)
		if err := existing.AddItem(productID, 2); err != nil {
			t.Fatalf("failed to seed existing cart fixture: %v", err)
		}
		repo.carts[customerID] = existing

		useCase := cart.NewAddItemUseCase(repo, txManager)
		err = useCase.Execute(context.Background(), cart.AddItemInput{
			CustomerID: testCustomerID,
			ProductID:  "product-1",
			Quantity:   3,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if txManager.calls != 1 {
			t.Fatalf("txManager.Do called %d times, want 1", txManager.calls)
		}
		if repo.saveCalls != 1 {
			t.Fatalf("Save called %d times, want 1", repo.saveCalls)
		}

		saved, ok := repo.carts[customerID]
		if !ok {
			t.Fatalf("cart was not saved")
		}
		items := saved.Items()
		if len(items) != 1 {
			t.Fatalf("len(Items()) = %d, want 1 (merged into the existing line item)", len(items))
		}
		if items[0].Quantity() != 5 { // 既存2 + 追加3
			t.Fatalf("Quantity() = %d, want 5", items[0].Quantity())
		}
	})

	t.Run("cart does not exist yet: a new cart is created and saved (find-or-create)", func(t *testing.T) {
		repo := newFakeCartRepository()
		txManager := &fakeTxManager{}
		useCase := cart.NewAddItemUseCase(repo, txManager)

		err := useCase.Execute(context.Background(), cart.AddItemInput{
			CustomerID: testCustomerID,
			ProductID:  "product-1",
			Quantity:   1,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		customerID, _ := domaincart.NewCustomerID(testCustomerID)
		saved, ok := repo.carts[customerID]
		if !ok {
			t.Fatalf("expected a new cart to be created and saved")
		}
		items := saved.Items()
		if len(items) != 1 || items[0].Quantity() != 1 {
			t.Fatalf("unexpected items on newly created cart: %+v", items)
		}
	})

	t.Run("FindByCustomerID error other than ErrNotFound propagates without being swallowed", func(t *testing.T) {
		repo := newFakeCartRepository()
		injectedErr := errors.New("db: connection lost")
		repo.findErr = injectedErr
		txManager := &fakeTxManager{}
		useCase := cart.NewAddItemUseCase(repo, txManager)

		err := useCase.Execute(context.Background(), cart.AddItemInput{
			CustomerID: testCustomerID,
			ProductID:  "product-1",
			Quantity:   1,
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, injectedErr) {
			t.Fatalf("expected error to wrap the injected error, got %v", err)
		}
		// find-or-create のフォールバックは shared.ErrNotFound のときだけ
		// 発動すべきであり、それ以外のインフラエラーは「カートが無いだけ」
		// として握りつぶさずそのまま呼び出し元へ返さなければならない。
		if repo.saveCalls != 0 {
			t.Fatalf("Save must not be called when FindByCustomerID fails with a non-NotFound error, called %d times", repo.saveCalls)
		}
	})

	t.Run("domain rule violation (zero quantity) is propagated and Save is not called", func(t *testing.T) {
		repo := newFakeCartRepository()
		txManager := &fakeTxManager{}
		useCase := cart.NewAddItemUseCase(repo, txManager)

		err := useCase.Execute(context.Background(), cart.AddItemInput{
			CustomerID: testCustomerID,
			ProductID:  "product-1",
			Quantity:   0,
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		if repo.saveCalls != 0 {
			t.Fatalf("Save must not be called when AddItem is rejected, called %d times", repo.saveCalls)
		}
	})
}
