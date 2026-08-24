package main

import (
	"net/http"

	"github.com/almondoo/wire"
	"gorm.io/gorm"

	cartcommand "github.com/almondoo/golang-ddd-sample/internal/application/cart/command"
	"github.com/almondoo/golang-ddd-sample/internal/application/cart/eventhandler"
	cartquery "github.com/almondoo/golang-ddd-sample/internal/application/cart/query"
	catalogcommand "github.com/almondoo/golang-ddd-sample/internal/application/catalog/command"
	catalogquery "github.com/almondoo/golang-ddd-sample/internal/application/catalog/query"
	ordercommand "github.com/almondoo/golang-ddd-sample/internal/application/order/command"
	orderquery "github.com/almondoo/golang-ddd-sample/internal/application/order/query"
	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
	"github.com/almondoo/golang-ddd-sample/internal/infrastructure/event"
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
	); err != nil {
		return nil, err
	}
	return db, nil
}

// provideEventBus はインメモリのイベントバスを生成し、購読を登録する。
//
// 購読の配線も composition root の責務 — order コンテキストは cart の
// 存在を知らないまま連携する。ClearCartOnOrderPlaced（cart 側の
// イベントハンドラ）を OrderPlacedEventName に紐づけるのはここだけであり、
// order パッケージ自身はどこにも cart を import していない。
// wire.Build の依存グラフだけでは「購読の登録」という副作用付きの手順は
// 表現できないため、手組みの provider として用意している。
func provideEventBus(clearCart *eventhandler.ClearCartOnOrderPlaced) *event.Bus {
	bus := event.NewBus()
	bus.Subscribe(order.OrderPlacedEventName, clearCart.Handle)
	return bus
}

// provideMux は 3 コンテキストのコントローラを mux に登録する。
func provideMux(
	catalogController *controller.CatalogController,
	cartController *controller.CartController,
	orderController *controller.OrderController,
) *http.ServeMux {
	mux := http.NewServeMux()
	catalogController.Register(mux)
	cartController.Register(mux)
	orderController.Register(mux)
	return mux
}

// infrastructureSet は DB 接続・トランザクション管理・イベントバスという、
// どのコンテキストからも参照される基盤コンポーネントをまとめたセットである。
// wire.Bind で「ポート（インターフェース）← 具体的な実装」の対応関係を
// 宣言している点が、手書き DI の txManager := persistence.NewTxManager(db)
// に相当する部分である。
var infrastructureSet = wire.NewSet(
	provideDB,
	persistence.NewTxManager,
	wire.Bind(new(tx.Manager), new(*persistence.TxManager)),
	provideEventBus,
	wire.Bind(new(shared.EventPublisher), new(*event.Bus)),
)

// catalogSet は catalog コンテキストのリポジトリ・クエリサービス・
// ユースケースをまとめたセットである。
var catalogSet = wire.NewSet(
	persistence.NewProductRepository,
	wire.Bind(new(catalog.Repository), new(*persistence.ProductRepository)),
	persistence.NewProductQuery,
	wire.Bind(new(catalogquery.ProductQueryService), new(*persistence.ProductQuery)),
	catalogcommand.NewRegisterProductUseCase,
	catalogcommand.NewChangePriceUseCase,
	catalogquery.NewListProductsUseCase,
	catalogquery.NewGetProductUseCase,
)

// cartSet は cart コンテキストのリポジトリ・クエリサービス・ユースケース・
// イベントハンドラをまとめたセットである。
var cartSet = wire.NewSet(
	persistence.NewCartRepository,
	wire.Bind(new(cart.Repository), new(*persistence.CartRepository)),
	persistence.NewCartQuery,
	wire.Bind(new(cartquery.CartQueryService), new(*persistence.CartQuery)),
	cartcommand.NewAddItemUseCase,
	cartcommand.NewRemoveItemUseCase,
	cartquery.NewGetCartUseCase,
	eventhandler.NewClearCartOnOrderPlaced,
)

// orderSet は order コンテキストのリポジトリ・クエリサービス・ユースケースを
// まとめたセットである。PlaceOrderUseCase は cart / catalog のリポジトリにも
// 依存するが、その配線は wire がグラフ全体から自動で解決する。
var orderSet = wire.NewSet(
	persistence.NewOrderRepository,
	wire.Bind(new(order.Repository), new(*persistence.OrderRepository)),
	persistence.NewOrderQuery,
	wire.Bind(new(orderquery.OrderQueryService), new(*persistence.OrderQuery)),
	ordercommand.NewPlaceOrderUseCase,
	ordercommand.NewPayOrderUseCase,
	ordercommand.NewShipOrderUseCase,
	ordercommand.NewCancelOrderUseCase,
	orderquery.NewGetOrderUseCase,
)

// controllerSet はプレゼンテーション層のコントローラと、それらを mux に
// 登録する provideMux をまとめたセットである。
var controllerSet = wire.NewSet(
	controller.NewCatalogController,
	controller.NewCartController,
	controller.NewOrderController,
	provideMux,
)
