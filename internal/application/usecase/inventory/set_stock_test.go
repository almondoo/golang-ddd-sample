package inventory_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/application/usecase/inventory"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// このファイルは SetStockUseCase の「find-or-create」分岐(既存在庫の更新 /
// 未登録商品の在庫を新規作成)を、手書きのフェイクだけで検証する。
// cart.AddItemUseCase のテストと同じ形の分岐であり、フェイクの作り方も
// 意図的に揃えている(add_item_test.go / place_order_test.go を参照)。

// ---- fakeInventoryRepository ----

type fakeInventoryRepository struct {
	stocks    map[domaininventory.ProductID]*domaininventory.Stock
	findErr   error
	saveErr   error
	saved     []*domaininventory.Stock
	saveCalls int
}

func newFakeInventoryRepository() *fakeInventoryRepository {
	return &fakeInventoryRepository{stocks: map[domaininventory.ProductID]*domaininventory.Stock{}}
}

func (r *fakeInventoryRepository) FindByProductID(_ context.Context, id domaininventory.ProductID) (*domaininventory.Stock, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	s, ok := r.stocks[id]
	if !ok {
		return nil, fmt.Errorf("stock not found: product=%s: %w", id, shared.ErrNotFound)
	}
	return s, nil
}

func (r *fakeInventoryRepository) Save(_ context.Context, s *domaininventory.Stock) error {
	r.saveCalls++
	if r.saveErr != nil {
		return r.saveErr
	}
	r.stocks[s.ProductID()] = s
	r.saved = append(r.saved, s)
	return nil
}

// ---- fakeTxManager ----

// fakeTxManager は fn(ctx) を素通しするだけの tx.Manager フェイクである。
// 本物のロールバック検証はしない — トランザクション境界が1回で呼ばれる
// ことだけを確認する(他パッケージの fakeTxManager と同じ方針)。
type fakeTxManager struct {
	calls int
}

func (m *fakeTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	m.calls++
	return fn(ctx)
}

var _ tx.Manager = (*fakeTxManager)(nil)

const testProductID = "product-1"

func TestSetStockUseCase_Execute(t *testing.T) {
	t.Run("updates the quantity of existing stock via SetQuantity", func(t *testing.T) {
		repo := newFakeInventoryRepository()
		txManager := &fakeTxManager{}
		productID, err := domaininventory.NewProductID(testProductID)
		if err != nil {
			t.Fatalf("failed to build product id fixture: %v", err)
		}
		existing, err := domaininventory.NewStock(productID, 10)
		if err != nil {
			t.Fatalf("failed to build stock fixture: %v", err)
		}
		repo.stocks[productID] = existing

		useCase := inventory.NewSetStockUseCase(repo, txManager)
		err = useCase.Execute(context.Background(), inventory.SetStockInput{
			ProductID: testProductID,
			Quantity:  25,
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
		saved, ok := repo.stocks[productID]
		if !ok {
			t.Fatalf("stock was not saved")
		}
		if saved.Quantity() != 25 {
			t.Fatalf("Quantity() = %d, want 25", saved.Quantity())
		}
	})

	t.Run("stock is not registered yet: a new stock is created (find-or-create)", func(t *testing.T) {
		repo := newFakeInventoryRepository()
		txManager := &fakeTxManager{}
		useCase := inventory.NewSetStockUseCase(repo, txManager)

		err := useCase.Execute(context.Background(), inventory.SetStockInput{
			ProductID: testProductID,
			Quantity:  15,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		productID, _ := domaininventory.NewProductID(testProductID)
		saved, ok := repo.stocks[productID]
		if !ok {
			t.Fatalf("expected a new stock to be created and saved")
		}
		if saved.Quantity() != 15 {
			t.Fatalf("Quantity() = %d, want 15", saved.Quantity())
		}
		if saved.Reserved() != 0 {
			t.Fatalf("Reserved() = %d, want 0 for a freshly created stock", saved.Reserved())
		}
	})

	t.Run("setting quantity below reserved is a domain rule violation and Save is not called", func(t *testing.T) {
		repo := newFakeInventoryRepository()
		txManager := &fakeTxManager{}
		productID, err := domaininventory.NewProductID(testProductID)
		if err != nil {
			t.Fatalf("failed to build product id fixture: %v", err)
		}
		existing, err := domaininventory.NewStock(productID, 10)
		if err != nil {
			t.Fatalf("failed to build stock fixture: %v", err)
		}
		if err := existing.Reserve(7); err != nil {
			t.Fatalf("failed to seed reserved quantity fixture: %v", err)
		}
		repo.stocks[productID] = existing

		useCase := inventory.NewSetStockUseCase(repo, txManager)
		err = useCase.Execute(context.Background(), inventory.SetStockInput{
			ProductID: testProductID,
			Quantity:  5, // reserved(7) を下回る
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		if repo.saveCalls != 0 {
			t.Fatalf("Save must not be called when SetQuantity is rejected, called %d times", repo.saveCalls)
		}
		// 元の在庫状態がそのまま残っていることも確認する(部分更新されていない)。
		if repo.stocks[productID].Quantity() != 10 {
			t.Fatalf("Quantity() = %d, want unchanged 10", repo.stocks[productID].Quantity())
		}
	})

	t.Run("FindByProductID error propagates", func(t *testing.T) {
		repo := newFakeInventoryRepository()
		injectedErr := errors.New("db: connection lost")
		repo.findErr = injectedErr
		txManager := &fakeTxManager{}
		useCase := inventory.NewSetStockUseCase(repo, txManager)

		err := useCase.Execute(context.Background(), inventory.SetStockInput{
			ProductID: testProductID,
			Quantity:  1,
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, injectedErr) {
			t.Fatalf("expected error to wrap the injected error, got %v", err)
		}
		if repo.saveCalls != 0 {
			t.Fatalf("Save must not be called when FindByProductID fails, called %d times", repo.saveCalls)
		}
	})
}
