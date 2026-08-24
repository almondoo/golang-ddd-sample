package main

import (
	"net/http"

	"github.com/almondoo/wire"
	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	cartusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/cart"
	catalogusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/catalog"
	couponusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/coupon"
	customerusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/customer"
	inventoryusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/inventory"
	orderusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/order"
	shippingusecase "github.com/almondoo/golang-ddd-sample/internal/application/usecase/shipping"
	"github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	"github.com/almondoo/golang-ddd-sample/internal/domain/customer"
	"github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shipping"
	"github.com/almondoo/golang-ddd-sample/internal/infrastructure/persistence"
	"github.com/almondoo/golang-ddd-sample/internal/presentation/controller"
)

// このファイルは通常ビルドに含まれる（wireinject タグは付けない）。
// wire が生成する wire_gen.go から呼び出されるのは、ここに定義した
// provideXxx 関数と各パッケージの既存コンストラクタ（NewYyy）である。
// つまり「配線の設計図」は wire.go（wireinject 専用）に、
// 「wire だけでは表現しきれない手組みの初期化ロジック」はこのファイルに、
// という役割分担にしている。

// provideDB は DSN から *gorm.DB を確立し、スキーマを用意する。
//
// サンプルなので AutoMigrate を使う。実運用ではマイグレーションツール
// (golang-migrate 等)でスキーマをバージョン管理する。
func provideDB(dsn string) (*gorm.DB, error) {
	db, err := persistence.NewDB(dsn)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&persistence.ProductModel{},
		&persistence.CartItemModel{},
		&persistence.OrderModel{},
		&persistence.OrderItemModel{},
		&persistence.CustomerModel{},
		&persistence.CustomerAddressModel{},
		&persistence.StockModel{},
		&persistence.ShipmentModel{},
		&persistence.CouponModel{},
	); err != nil {
		return nil, err
	}
	return db, nil
}

// provideMux は 7 コンテキストのコントローラを mux に登録する。
func provideMux(
	catalogController *controller.CatalogController,
	cartController *controller.CartController,
	orderController *controller.OrderController,
	customerController *controller.CustomerController,
	inventoryController *controller.InventoryController,
	shippingController *controller.ShippingController,
	couponController *controller.CouponController,
) *http.ServeMux {
	mux := http.NewServeMux()
	catalogController.Register(mux)
	cartController.Register(mux)
	orderController.Register(mux)
	customerController.Register(mux)
	inventoryController.Register(mux)
	shippingController.Register(mux)
	couponController.Register(mux)
	return mux
}

// infrastructureSet は DB 接続・トランザクション管理という、
// どのコンテキストからも参照される基盤コンポーネントをまとめたセットである。
// wire.Bind で「ポート（インターフェース）← 具体的な実装」の対応関係を
// 宣言している点が、手書き DI の txManager := persistence.NewTxManager(db)
// に相当する部分である。
var infrastructureSet = wire.NewSet(
	provideDB,
	persistence.NewTxManager,
	wire.Bind(new(tx.Manager), new(*persistence.TxManager)),
)

// catalogSet は catalog コンテキストのリポジトリ・クエリサービス・
// ユースケースをまとめたセットである。
var catalogSet = wire.NewSet(
	persistence.NewProductRepository,
	wire.Bind(new(catalog.Repository), new(*persistence.ProductRepository)),
	persistence.NewProductQuery,
	wire.Bind(new(catalogusecase.ProductQueryService), new(*persistence.ProductQuery)),
	catalogusecase.NewRegisterProductUseCase,
	catalogusecase.NewChangePriceUseCase,
	catalogusecase.NewListProductsUseCase,
	catalogusecase.NewGetProductUseCase,
)

// cartSet は cart コンテキストのリポジトリ・クエリサービス・ユースケースを
// まとめたセットである。
var cartSet = wire.NewSet(
	persistence.NewCartRepository,
	wire.Bind(new(cart.Repository), new(*persistence.CartRepository)),
	persistence.NewCartQuery,
	wire.Bind(new(cartusecase.CartQueryService), new(*persistence.CartQuery)),
	cartusecase.NewAddItemUseCase,
	cartusecase.NewRemoveItemUseCase,
	cartusecase.NewGetCartUseCase,
)

// orderSet は order コンテキストのリポジトリ・クエリサービス・ユースケースを
// まとめたセットである。PlaceOrderUseCase / ShipOrderUseCase / CancelOrderUseCase は
// cart / catalog / customer / inventory / coupon のリポジトリにも依存するが、
// その配線は wire がグラフ全体から自動で解決する。
var orderSet = wire.NewSet(
	persistence.NewOrderRepository,
	wire.Bind(new(order.Repository), new(*persistence.OrderRepository)),
	persistence.NewOrderQuery,
	wire.Bind(new(orderusecase.OrderQueryService), new(*persistence.OrderQuery)),
	orderusecase.NewPlaceOrderUseCase,
	orderusecase.NewPayOrderUseCase,
	orderusecase.NewShipOrderUseCase,
	orderusecase.NewCancelOrderUseCase,
	orderusecase.NewGetOrderUseCase,
)

// customerSet は customer コンテキストのリポジトリ・クエリサービス・
// ユースケースをまとめたセットである。
var customerSet = wire.NewSet(
	persistence.NewCustomerRepository,
	wire.Bind(new(customer.Repository), new(*persistence.CustomerRepository)),
	persistence.NewCustomerQuery,
	wire.Bind(new(customerusecase.CustomerQueryService), new(*persistence.CustomerQuery)),
	customerusecase.NewRegisterCustomerUseCase,
	customerusecase.NewAddAddressUseCase,
	customerusecase.NewChangeDefaultAddressUseCase,
	customerusecase.NewRemoveAddressUseCase,
	customerusecase.NewGetCustomerUseCase,
)

// inventorySet は inventory コンテキストのリポジトリ・クエリサービス・
// ユースケースをまとめたセットである。
var inventorySet = wire.NewSet(
	persistence.NewStockRepository,
	wire.Bind(new(inventory.Repository), new(*persistence.StockRepository)),
	persistence.NewStockQuery,
	wire.Bind(new(inventoryusecase.StockQueryService), new(*persistence.StockQuery)),
	inventoryusecase.NewSetStockUseCase,
	inventoryusecase.NewGetStockUseCase,
)

// shippingSet は shipping コンテキストのリポジトリ・クエリサービス・
// ユースケースをまとめたセットである。
var shippingSet = wire.NewSet(
	persistence.NewShipmentRepository,
	wire.Bind(new(shipping.Repository), new(*persistence.ShipmentRepository)),
	persistence.NewShipmentQuery,
	wire.Bind(new(shippingusecase.ShipmentQueryService), new(*persistence.ShipmentQuery)),
	shippingusecase.NewDeliverShipmentUseCase,
	shippingusecase.NewGetShipmentUseCase,
)

// couponSet は coupon コンテキストのリポジトリ・クエリサービス・ユースケースを
// まとめたセットである。
var couponSet = wire.NewSet(
	persistence.NewCouponRepository,
	wire.Bind(new(coupon.Repository), new(*persistence.CouponRepository)),
	persistence.NewCouponQuery,
	wire.Bind(new(couponusecase.CouponQueryService), new(*persistence.CouponQuery)),
	couponusecase.NewIssueCouponUseCase,
	couponusecase.NewGetCouponUseCase,
)

// controllerSet はプレゼンテーション層のコントローラと、それらを mux に
// 登録する provideMux をまとめたセットである。
var controllerSet = wire.NewSet(
	controller.NewCatalogController,
	controller.NewCartController,
	controller.NewOrderController,
	controller.NewCustomerController,
	controller.NewInventoryController,
	controller.NewShippingController,
	controller.NewCouponController,
	provideMux,
)
