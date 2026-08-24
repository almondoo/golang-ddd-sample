package order

import "time"

// OrderPlacedEventName は OrderPlaced イベントの種別を表す一意な文字列である。
// "order.placed" のようなドット区切りは、将来 order コンテキストが
// 他のイベント（例: order.canceled）を持った場合にも一貫した命名で
// 見分けられるようにするための慣習である。
const OrderPlacedEventName = "order.placed"

// OrderPlaced は「注文が確定した」という事実を表すドメインイベントである。
//
// なぜ Pay/Ship/Cancel はイベントを記録しないのか（設計判断のメモ）:
// 本サンプルでは学習用に、他コンテキスト（cart）が反応すべき最初の
// マイルストーンである「注文確定」だけをイベント化している。
// 実際のシステムでは OrderPaid や OrderShipped も同様にイベント化し、
// 例えば「発送完了メールを送る」「在庫を確定的に引き当てる」といった
// 副作用のトリガーに使うのが自然な拡張である。
type OrderPlaced struct {
	orderID    OrderID
	customerID CustomerID
	occurredAt time.Time
}

// NewOrderPlaced は OrderPlaced イベントを生成する。
// 発生時刻はこの関数が呼ばれた瞬間（= 注文が確定した瞬間）を採用する。
func NewOrderPlaced(orderID OrderID, customerID CustomerID) OrderPlaced {
	return OrderPlaced{
		orderID:    orderID,
		customerID: customerID,
		occurredAt: time.Now(),
	}
}

// EventName は shared.DomainEvent インターフェースの実装である。
func (e OrderPlaced) EventName() string {
	return OrderPlacedEventName
}

// OccurredAt は shared.DomainEvent インターフェースの実装である。
func (e OrderPlaced) OccurredAt() time.Time {
	return e.occurredAt
}

// OrderID はこのイベントが指す注文の ID を返す。
func (e OrderPlaced) OrderID() OrderID {
	return e.orderID
}

// CustomerID はこの注文を行った顧客の ID を返す。
//
// カート側のイベントハンドラ（clear_cart_on_order_placed.go）は
// この値をもとに「どの顧客のカートを空にすべきか」を判断する。
// order パッケージが cart パッケージを知らなくても、イベントに
// 必要な情報を載せておくことで、受け手側だけが変換を行えばよくなる。
func (e OrderPlaced) CustomerID() CustomerID {
	return e.customerID
}
