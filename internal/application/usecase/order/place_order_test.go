package order_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/application/usecase/order"
	domaincart "github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	domaincatalog "github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	domaincoupon "github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	domainorder "github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// このファイルは PlaceOrderUseCase（5リポジトリ + tx.Manager に依存する
// 本サンプル最複雑のユースケース）を、モックライブラリを使わず手書きの
// フェイクだけで検証する。各フェイクはポート（1〜2メソッドの interface）を
// map ベースで愚直に実装するだけであり、これ自体が
// 「ポートを切っておくとテストが書きやすい」という設計判断の見返りを
// 示す教材でもある（docs/specs/ddd-improvements.md 項目5）。

// ---- fakeCustomerRepository ----

// fakeCustomerRepository は domaincustomer.Repository の手書きフェイクである。
// map をバックエンドにするだけで、GORM や実 DB を一切必要としない。
type fakeCustomerRepository struct {
	customers map[domaincustomer.CustomerID]*domaincustomer.Customer
	findErr   error
	saveErr   error
	saved     []*domaincustomer.Customer
}

func newFakeCustomerRepository() *fakeCustomerRepository {
	return &fakeCustomerRepository{customers: map[domaincustomer.CustomerID]*domaincustomer.Customer{}}
}

func (r *fakeCustomerRepository) FindByID(_ context.Context, id domaincustomer.CustomerID) (*domaincustomer.Customer, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	c, ok := r.customers[id]
	if !ok {
		return nil, fmt.Errorf("customer not found: id=%s: %w", id, shared.ErrNotFound)
	}
	return c, nil
}

func (r *fakeCustomerRepository) Save(_ context.Context, c *domaincustomer.Customer) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.customers[c.ID()] = c
	r.saved = append(r.saved, c)
	return nil
}

// ---- fakeCartRepository ----

type fakeCartRepository struct {
	carts   map[domaincart.CustomerID]*domaincart.Cart
	findErr error
	saveErr error
	saved   []*domaincart.Cart
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
	if r.saveErr != nil {
		return r.saveErr
	}
	r.carts[c.CustomerID()] = c
	r.saved = append(r.saved, c)
	return nil
}

// ---- fakeCatalogRepository ----

type fakeCatalogRepository struct {
	products map[domaincatalog.ProductID]*domaincatalog.Product
	findErr  error
	saveErr  error
}

func newFakeCatalogRepository() *fakeCatalogRepository {
	return &fakeCatalogRepository{products: map[domaincatalog.ProductID]*domaincatalog.Product{}}
}

func (r *fakeCatalogRepository) FindByID(_ context.Context, id domaincatalog.ProductID) (*domaincatalog.Product, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	p, ok := r.products[id]
	if !ok {
		return nil, fmt.Errorf("product not found: id=%s: %w", id, shared.ErrNotFound)
	}
	return p, nil
}

func (r *fakeCatalogRepository) Save(_ context.Context, p *domaincatalog.Product) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.products[p.ID()] = p
	return nil
}

// ---- fakeInventoryRepository ----

type fakeInventoryRepository struct {
	stocks  map[domaininventory.ProductID]*domaininventory.Stock
	findErr error
	saveErr error
	saved   []*domaininventory.Stock
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
	if r.saveErr != nil {
		return r.saveErr
	}
	r.stocks[s.ProductID()] = s
	r.saved = append(r.saved, s)
	return nil
}

// ---- fakeCouponRepository ----

type fakeCouponRepository struct {
	coupons   map[domaincoupon.CouponCode]*domaincoupon.Coupon
	findErr   error
	saveErr   error
	saved     []*domaincoupon.Coupon
	findCalls int
}

func newFakeCouponRepository() *fakeCouponRepository {
	return &fakeCouponRepository{coupons: map[domaincoupon.CouponCode]*domaincoupon.Coupon{}}
}

func (r *fakeCouponRepository) FindByCode(_ context.Context, code domaincoupon.CouponCode) (*domaincoupon.Coupon, error) {
	r.findCalls++
	if r.findErr != nil {
		return nil, r.findErr
	}
	c, ok := r.coupons[code]
	if !ok {
		return nil, fmt.Errorf("coupon not found: code=%s: %w", code.String(), shared.ErrNotFound)
	}
	return c, nil
}

func (r *fakeCouponRepository) Save(_ context.Context, c *domaincoupon.Coupon) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.coupons[c.Code()] = c
	r.saved = append(r.saved, c)
	return nil
}

// ---- fakeOrderRepository ----

type fakeOrderRepository struct {
	orders  map[domainorder.OrderID]*domainorder.Order
	saveErr error
	saved   []*domainorder.Order
}

func newFakeOrderRepository() *fakeOrderRepository {
	return &fakeOrderRepository{orders: map[domainorder.OrderID]*domainorder.Order{}}
}

func (r *fakeOrderRepository) FindByID(_ context.Context, id domainorder.OrderID) (*domainorder.Order, error) {
	o, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found: id=%s: %w", id, shared.ErrNotFound)
	}
	return o, nil
}

func (r *fakeOrderRepository) Save(_ context.Context, o *domainorder.Order) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.orders[o.ID()] = o
	r.saved = append(r.saved, o)
	return nil
}

// ---- fakeTxManager ----

// fakeTxManager は tx.Manager の手書きフェイクである。fn(ctx) をそのまま
// 呼び出すだけで、実際のコミット・ロールバックは一切行わない。
//
// 本物のロールバック検証（fn がエラーを返したときに DB への書き込みが
// 本当に取り消されるか）は、実 DB を使ったインフラ層の統合テストの責務
// であり、ここでは行わない。このフェイクで確認できるのは
// 「txManager.Do がユースケース実行につき 1 回だけ呼ばれる」＝
// トランザクション境界が意図通り 1 回であること、および fn が返した
// エラーがそのまま呼び出し元へ伝播することだけである。
type fakeTxManager struct {
	calls int
}

func (m *fakeTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	m.calls++
	return fn(ctx)
}

var _ tx.Manager = (*fakeTxManager)(nil)

// ---- テスト共通のフィクスチャ ----

const testCustomerID = "customer-1"

type fixture struct {
	t *testing.T

	customerRepo  *fakeCustomerRepository
	cartRepo      *fakeCartRepository
	catalogRepo   *fakeCatalogRepository
	inventoryRepo *fakeInventoryRepository
	couponRepo    *fakeCouponRepository
	orderRepo     *fakeOrderRepository
	txManager     *fakeTxManager

	useCase *order.PlaceOrderUseCase
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{
		t:             t,
		customerRepo:  newFakeCustomerRepository(),
		cartRepo:      newFakeCartRepository(),
		catalogRepo:   newFakeCatalogRepository(),
		inventoryRepo: newFakeInventoryRepository(),
		couponRepo:    newFakeCouponRepository(),
		orderRepo:     newFakeOrderRepository(),
		txManager:     &fakeTxManager{},
	}
	f.useCase = order.NewPlaceOrderUseCase(
		f.orderRepo, f.cartRepo, f.catalogRepo, f.customerRepo, f.inventoryRepo, f.couponRepo, f.txManager,
	)
	return f
}

// withCustomer は testCustomerID の顧客を実在させる。
func (f *fixture) withCustomer() *fixture {
	f.t.Helper()
	id, err := domaincustomer.NewCustomerID(testCustomerID)
	if err != nil {
		f.t.Fatalf("failed to build customer id fixture: %v", err)
	}
	f.customerRepo.customers[id] = domaincustomer.ReconstructCustomer(id, "Taro", "taro@example.com", nil)
	return f
}

// cart は testCustomerID のカートを取得する。まだ無ければ空カートを
// 作って登録してから返す（find-or-create ではなく「テストの前提として
// 用意する」ためのヘルパーである点に注意）。
func (f *fixture) cart() *domaincart.Cart {
	f.t.Helper()
	id, err := domaincart.NewCustomerID(testCustomerID)
	if err != nil {
		f.t.Fatalf("failed to build cart customer id fixture: %v", err)
	}
	c, ok := f.cartRepo.carts[id]
	if !ok {
		c = domaincart.NewCart(id)
		f.cartRepo.carts[id] = c
	}
	return c
}

// addCartItem は testCustomerID のカートに明細を追加する。
func (f *fixture) addCartItem(productID string, quantity int) {
	f.t.Helper()
	pid, err := domaincart.NewProductID(productID)
	if err != nil {
		f.t.Fatalf("failed to build cart product id fixture: %v", err)
	}
	if err := f.cart().AddItem(pid, quantity); err != nil {
		f.t.Fatalf("failed to add cart item fixture: %v", err)
	}
}

// addProduct は catalog に商品を登録し、生成された Product を返す。
func (f *fixture) addProduct(name string, priceAmount int64) *domaincatalog.Product {
	f.t.Helper()
	price := mustMoney(f.t, priceAmount)
	p, err := domaincatalog.NewProduct(name, "テスト用商品説明", price)
	if err != nil {
		f.t.Fatalf("failed to build product fixture: %v", err)
	}
	f.catalogRepo.products[p.ID()] = p
	return p
}

// addStock は指定商品の在庫を quantity で登録する。
func (f *fixture) addStock(productID domaincatalog.ProductID, quantity int) {
	f.t.Helper()
	pid, err := domaininventory.NewProductID(productID.String())
	if err != nil {
		f.t.Fatalf("failed to build inventory product id fixture: %v", err)
	}
	s, err := domaininventory.NewStock(pid, quantity)
	if err != nil {
		f.t.Fatalf("failed to build stock fixture: %v", err)
	}
	f.inventoryRepo.stocks[pid] = s
}

// addCoupon はクーポンをリポジトリに登録する。
func (f *fixture) addCoupon(c *domaincoupon.Coupon) {
	f.couponRepo.coupons[c.Code()] = c
}

func mustMoney(t *testing.T, amount int64) shared.Money {
	t.Helper()
	m, err := shared.NewMoney(amount, shared.JPY)
	if err != nil {
		t.Fatalf("failed to build money fixture: %v", err)
	}
	return m
}

// farFutureExpiry / farPastExpiry は「PlaceOrderUseCase.Execute が内部で
// time.Now() を直接呼んでいる」（ドメイン層の NewOrder/Coupon.Use だけが
// now を引数化されており、アプリケーション層の呼び出し元自体は未注入）
// という現状の実装を踏まえ、実行時刻に依存せず期限切れ・有効を
// 決定的にテストするための固定値である。
var (
	farFutureExpiry = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	farPastExpiry   = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
)

// ---- 正常系 ----

func TestPlaceOrderUseCase_Execute_Success(t *testing.T) {
	t.Run("cart with two items is placed as an order with catalog snapshots and reserved stock", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		productA := f.addProduct("商品A", 1000)
		productB := f.addProduct("商品B", 1500)
		f.addStock(productA.ID(), 5)
		f.addStock(productB.ID(), 5)
		f.addCartItem(productA.ID().String(), 2)
		f.addCartItem(productB.ID().String(), 1)

		out, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.OrderID == "" {
			t.Fatal("expected non-empty OrderID in output")
		}

		if f.txManager.calls != 1 {
			t.Fatalf("txManager.Do called %d times, want 1 (1 usecase execution = 1 transaction)", f.txManager.calls)
		}

		savedOrder, ok := f.orderRepo.orders[domainorder.OrderID(out.OrderID)]
		if !ok {
			t.Fatalf("order %s was not saved", out.OrderID)
		}
		if savedOrder.Status() != domainorder.StatusPending {
			t.Fatalf("Status() = %v, want %v", savedOrder.Status(), domainorder.StatusPending)
		}

		items := savedOrder.Items()
		if len(items) != 2 {
			t.Fatalf("len(Items()) = %d, want 2", len(items))
		}
		byProductID := map[string]domainorder.OrderItem{}
		for _, it := range items {
			byProductID[it.ProductID()] = it
		}

		// 明細のスナップショット(商品名・単価)が catalog の現在値からコピー
		// されていることを確認する。
		gotA, ok := byProductID[productA.ID().String()]
		if !ok {
			t.Fatalf("order item for product A not found")
		}
		if gotA.ProductName() != "商品A" || gotA.UnitPrice().Amount() != 1000 || gotA.Quantity() != 2 {
			t.Fatalf("order item A = %+v, want name=商品A price=1000 qty=2", gotA)
		}
		gotB, ok := byProductID[productB.ID().String()]
		if !ok {
			t.Fatalf("order item for product B not found")
		}
		if gotB.ProductName() != "商品B" || gotB.UnitPrice().Amount() != 1500 || gotB.Quantity() != 1 {
			t.Fatalf("order item B = %+v, want name=商品B price=1500 qty=1", gotB)
		}

		total, err := savedOrder.TotalAmount()
		if err != nil {
			t.Fatalf("TotalAmount() unexpected error: %v", err)
		}
		if total.Amount() != 3500 { // 1000*2 + 1500*1
			t.Fatalf("TotalAmount().Amount() = %d, want 3500", total.Amount())
		}

		// カート数量に応じて在庫が引き当てられていることを確認する。
		stockPidA, _ := domaininventory.NewProductID(productA.ID().String())
		if got := f.inventoryRepo.stocks[stockPidA].Reserved(); got != 2 {
			t.Fatalf("stock A Reserved() = %d, want 2", got)
		}
		stockPidB, _ := domaininventory.NewProductID(productB.ID().String())
		if got := f.inventoryRepo.stocks[stockPidB].Reserved(); got != 1 {
			t.Fatalf("stock B Reserved() = %d, want 1", got)
		}

		// カートが空になった状態で保存されていることを確認する。
		cartCID, _ := domaincart.NewCustomerID(testCustomerID)
		savedCart, ok := f.cartRepo.carts[cartCID]
		if !ok {
			t.Fatalf("cart was not saved")
		}
		if !savedCart.IsEmpty() {
			t.Fatalf("expected cart to be cleared, got %d items", len(savedCart.Items()))
		}
	})

	t.Run("valid coupon discounts the payable amount and is marked as used", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 5)
		f.addCartItem(product.ID().String(), 2) // total 2000

		code, err := domaincoupon.NewCouponCode("SUMMER500")
		if err != nil {
			t.Fatalf("failed to build coupon code fixture: %v", err)
		}
		coupon, err := domaincoupon.NewAmountCoupon(code, mustMoney(t, 500), farFutureExpiry, 10)
		if err != nil {
			t.Fatalf("failed to build coupon fixture: %v", err)
		}
		f.addCoupon(coupon)

		out, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{
			CustomerID: testCustomerID,
			CouponCode: "SUMMER500",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		savedOrder := f.orderRepo.orders[domainorder.OrderID(out.OrderID)]
		if savedOrder.CouponCode() != "SUMMER500" {
			t.Fatalf("CouponCode() = %q, want %q", savedOrder.CouponCode(), "SUMMER500")
		}
		if savedOrder.DiscountAmount().Amount() != 500 {
			t.Fatalf("DiscountAmount().Amount() = %d, want 500", savedOrder.DiscountAmount().Amount())
		}
		payable, err := savedOrder.PayableAmount()
		if err != nil {
			t.Fatalf("PayableAmount() unexpected error: %v", err)
		}
		if payable.Amount() != 1500 { // 2000 - 500
			t.Fatalf("PayableAmount().Amount() = %d, want 1500", payable.Amount())
		}

		// クーポンが 1 回消費され、その状態が保存されていることを確認する。
		savedCoupon, ok := f.couponRepo.coupons[code]
		if !ok {
			t.Fatalf("coupon was not saved")
		}
		if savedCoupon.UsedCount() != 1 {
			t.Fatalf("UsedCount() = %d, want 1", savedCoupon.UsedCount())
		}
	})
}

// ---- 異常系(ドメインルール違反) ----

func TestPlaceOrderUseCase_Execute_DomainRuleViolations(t *testing.T) {
	t.Run("customer does not exist", func(t *testing.T) {
		f := newFixture(t)
		// withCustomer を呼ばないことで「顧客が実在しない」を再現する。

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		if len(f.orderRepo.saved) != 0 {
			t.Fatalf("order must not be saved, but %d were saved", len(f.orderRepo.saved))
		}
	})

	t.Run("cart does not exist", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		// f.cart() を呼ばないため cartRepo は空のまま = FindByCustomerID が
		// shared.ErrNotFound を返す状況を再現する。

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("cart exists but is empty", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		f.cart() // 空カートだけを登録する。

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("cart item references a product that is not in the catalog", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		f.addCartItem("nonexistent-product", 1)
		// catalogRepo には何も登録しない。

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("stock is not registered for the product", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addCartItem(product.ID().String(), 1)
		// inventoryRepo には何も登録しない。

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("stock is insufficient for the requested quantity", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 1) // 在庫は1個だけ
		f.addCartItem(product.ID().String(), 2)

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("coupon code has an invalid format", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 5)
		f.addCartItem(product.ID().String(), 1)

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{
			CustomerID: testCustomerID,
			CouponCode: "bad", // 4文字未満・小文字を含み CouponCode の形式ルールに違反する
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		// 形式チェックは CouponCode コンストラクタが行うため、
		// リポジトリへの問い合わせ自体が発生しないことも確認する。
		if f.couponRepo.findCalls != 0 {
			t.Fatalf("FindByCode should not be called for a malformed code, called %d times", f.couponRepo.findCalls)
		}
	})

	t.Run("coupon code is not registered", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 5)
		f.addCartItem(product.ID().String(), 1)
		// couponRepo には何も登録しない。

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{
			CustomerID: testCustomerID,
			CouponCode: "UNKNOWN1",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("coupon is expired", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 5)
		f.addCartItem(product.ID().String(), 1)

		code, err := domaincoupon.NewCouponCode("EXPIRED1")
		if err != nil {
			t.Fatalf("failed to build coupon code fixture: %v", err)
		}
		coupon, err := domaincoupon.NewAmountCoupon(code, mustMoney(t, 100), farPastExpiry, 10)
		if err != nil {
			t.Fatalf("failed to build coupon fixture: %v", err)
		}
		f.addCoupon(coupon)

		_, err = f.useCase.Execute(context.Background(), order.PlaceOrderInput{
			CustomerID: testCustomerID,
			CouponCode: "EXPIRED1",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})

	t.Run("coupon has reached its usage limit", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 5)
		f.addCartItem(product.ID().String(), 1)

		code, err := domaincoupon.NewCouponCode("LIMIT001")
		if err != nil {
			t.Fatalf("failed to build coupon code fixture: %v", err)
		}
		coupon, err := domaincoupon.NewAmountCoupon(code, mustMoney(t, 100), farFutureExpiry, 1)
		if err != nil {
			t.Fatalf("failed to build coupon fixture: %v", err)
		}
		// usageLimit=1 のクーポンを事前に 1 回消費させ、「利用上限に
		// 達したクーポン」という状態をフィクスチャとして作る。
		if err := coupon.Use(time.Now()); err != nil {
			t.Fatalf("failed to pre-consume coupon fixture: %v", err)
		}
		f.addCoupon(coupon)

		_, err = f.useCase.Execute(context.Background(), order.PlaceOrderInput{
			CustomerID: testCustomerID,
			CouponCode: "LIMIT001",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
	})
}

// ---- 異常系(インフラ層エラーの伝播) ----

func TestPlaceOrderUseCase_Execute_RepositoryErrorPropagation(t *testing.T) {
	t.Run("order repository Save failure propagates as-is", func(t *testing.T) {
		f := newFixture(t)
		f.withCustomer()
		product := f.addProduct("商品A", 1000)
		f.addStock(product.ID(), 5)
		f.addCartItem(product.ID().String(), 1)

		injectedErr := errors.New("db: connection lost")
		f.orderRepo.saveErr = injectedErr

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: testCustomerID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, injectedErr) {
			t.Fatalf("expected error to wrap injected error, got %v", err)
		}
		// インフラ層のエラーはドメインルール違反(422)ではなく、
		// そのまま(500 相当)呼び出し元へ伝播すべきである。
		if shared.IsDomainRuleError(err) {
			t.Fatalf("infra error must not be classified as a domain rule error: %v", err)
		}

		// fakeTxManager は fn(ctx) をそのまま呼ぶだけで実際のロールバックは
		// 行わない(コメント参照)。そのためこのテストで確認できるのは
		// 「Order の Save が失敗した場合、それ以降の手順である
		// cart.Clear()+Save が実行されないままエラーが返る」という
		// ユースケース内の"手順の順序"だけであり、在庫引当が実際に
		// ロールバックされるかどうかはこのテストの対象外である
		// (実ロールバックの検証はインフラ層の統合テストの責務)。
		cartCID, _ := domaincart.NewCustomerID(testCustomerID)
		savedCart, ok := f.cartRepo.carts[cartCID]
		if !ok {
			t.Fatalf("cart fixture must still be present")
		}
		if savedCart.IsEmpty() {
			t.Fatalf("cart must not have been cleared because order Save failed before cart.Save was reached")
		}
	})
}

// ---- 入力検証(トランザクションを開く前に弾かれることの確認) ----

func TestPlaceOrderUseCase_Execute_InputValidation(t *testing.T) {
	t.Run("empty customer id is rejected before opening a transaction", func(t *testing.T) {
		f := newFixture(t)

		_, err := f.useCase.Execute(context.Background(), order.PlaceOrderInput{CustomerID: ""})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !shared.IsDomainRuleError(err) {
			t.Fatalf("expected domain rule error, got %v", err)
		}
		// customerID の妥当性検証は txManager.Do の外(Execute の冒頭)で
		// 行われるため、不正な入力ではトランザクションが一度も開かれない
		// ことを確認する。
		if f.txManager.calls != 0 {
			t.Fatalf("txManager.Do must not be called when input validation fails early, got %d calls", f.txManager.calls)
		}
	})
}
