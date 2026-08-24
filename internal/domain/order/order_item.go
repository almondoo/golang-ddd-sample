package order

import "github.com/almondoo/golang-ddd-sample/internal/domain/shared"

// OrderItem は注文明細（ある商品を、いくらで、何個購入したか）を表す
// 値オブジェクトである。
//
// 【最重要の設計判断: なぜ商品への参照ではなく「スナップショット」を持つのか】
// cart.CartItem は productID しか持たなかった（価格は catalog を都度参照する
// 設計だった）。しかし OrderItem は productName と unitPrice を
// 「注文が確定した瞬間の値」としてコピーして保持する。
//
// 理由は、注文が「過去のある時点で成立した取引の記録」だからである。
// もし OrderItem が ProductID だけを持ち、表示のたびに catalog の最新価格を
// 引いてくる設計にすると、次のような業務的に許されない事態が起きる。
//   - 注文確定後に商品の値段が改定されると、過去の注文の合計金額まで
//     変わって見えてしまう（会計・監査上、過去の取引額は不変であるべき）。
//   - 商品名が変更・廃止されると、過去の注文明細の表示が壊れる、あるいは
//     消えてしまう。
//
// そのため OrderItem は catalog.Product への参照を一切持たず、注文確定時点
// での名前・単価を「写し取って」自分自身の中に保持する。これにより、
// catalog 側でその後何が起きても、過去の注文は当時の事実のまま不変で
// あり続けられる。この考え方は会計ドメインでは "point-in-time snapshot"
// と呼ばれ、DDD における値オブジェクトの典型的な使いどころでもある。
type OrderItem struct {
	productID   string
	productName string
	unitPrice   shared.Money
	quantity    int
}

// NewOrderItem は注文確定時に使うコンストラクタである。
// カート明細と、その時点の catalog.Product から得た情報を組み合わせて
// スナップショットを作る（実際の組み立ては application 層の
// PlaceOrderUseCase が担う）。
func NewOrderItem(productID, productName string, unitPrice shared.Money, quantity int) (OrderItem, error) {
	if productID == "" {
		return OrderItem{}, shared.NewDomainRuleError("order: order item product id must not be empty")
	}
	if productName == "" {
		return OrderItem{}, shared.NewDomainRuleError("order: order item product name must not be empty")
	}
	if quantity < 1 {
		return OrderItem{}, shared.NewDomainRuleError("order: order item quantity must be at least 1, got %d", quantity)
	}
	return OrderItem{
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}, nil
}

// ReconstructOrderItem は永続化層から読み込んだデータをもとに OrderItem を
// 再構築する。DB に保存されている値は過去に NewOrderItem の検証を通過
// 済みという前提のもと、再度の検証は行わない（ReconstructCartItem 等と
// 同じ設計上の理由による）。
func ReconstructOrderItem(productID, productName string, unitPrice shared.Money, quantity int) OrderItem {
	return OrderItem{
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}
}

// ProductID はこの明細が指す商品の ID（スナップショット時点の文字列）を返す。
func (i OrderItem) ProductID() string {
	return i.productID
}

// ProductName はこの明細の商品名（注文確定時点でのスナップショット）を返す。
func (i OrderItem) ProductName() string {
	return i.productName
}

// UnitPrice はこの明細の単価（注文確定時点でのスナップショット）を返す。
func (i OrderItem) UnitPrice() shared.Money {
	return i.unitPrice
}

// Quantity はこの明細の数量を返す。
func (i OrderItem) Quantity() int {
	return i.quantity
}

// Subtotal は単価 × 数量の小計を返す。
func (i OrderItem) Subtotal() (shared.Money, error) {
	return i.unitPrice.Multiply(i.quantity)
}
