package order

import (
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// Order は注文を表す集約ルート（Aggregate Root）である。
//
// 【集約境界】
// Order は自身の明細（OrderItem）を内部に抱え、他のコンテキスト
// （cart, catalog）の集約を直接参照しない。cart・catalog への依存は
// 「注文を組み立てる」というオーケストレーションの中で必要になるが、
// それはドメイン層ではなくアプリケーション層（PlaceOrderUseCase）の
// 責務である。order ドメインパッケージが cart / catalog パッケージを
// import しないことは、コンテキストの自律性を保つための最重要ルールである。
//
// shared.AggregateBase を埋め込んでいるのは、Order が他コンテキストへ
// 通知すべき重要な出来事（OrderPlaced）を持つためである。catalog.Product
// がイベントを持たないのとは対照的に、Order は「注文確定」という
// 業務上の一大イベントの発生源になる。
type Order struct {
	shared.AggregateBase

	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     Status
	placedAt   time.Time
}

// NewOrder は新規注文を確定するコンストラクタである。
//
// 呼び出し時点で明細が 1 件もない注文は業務的に意味を持たない
// （「何も買わない注文」は存在しない）ため、items の非空をここで検証する。
// 個々の OrderItem 自体の妥当性（数量 1 以上等）は NewOrderItem
// コンストラクタがすでに保証済みという前提に立つ。
//
// 生成と同時に StatusPending をセットし、OrderPlaced イベントを記録する。
// このイベントはまだ「配信」されていない点に注意する。配信（Publish）は
// アプリケーション層が Save の成功後に行う責務であり、集約はあくまで
// 「注文確定という事実が発生した」ことを記録するだけに留める
// （理由の詳細は shared.AggregateBase のコメントを参照）。
func NewOrder(customerID CustomerID, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, shared.NewDomainRuleError("order: order must contain at least one item")
	}

	itemsCopy := make([]OrderItem, len(items))
	copy(itemsCopy, items)

	o := &Order{
		id:         GenerateOrderID(),
		customerID: customerID,
		items:      itemsCopy,
		status:     StatusPending,
		placedAt:   time.Now(),
	}
	o.Record(NewOrderPlaced(o.id, o.customerID))
	return o, nil
}

// ReconstructOrder は永続化層から読み込んだデータをもとに Order を再構築する。
//
// NewOrder との違いは ReconstructProduct / ReconstructCart と同じ理由による。
// 加えて重要な点として、ReconstructOrder は OrderPlaced イベントを
// 記録しない。イベントは「今まさに起きた出来事」を表すものであり、
// すでに過去に確定して DB に保存済みの注文を読み込むたびに
// OrderPlaced を再送してしまうと、カートを何度も空にしようとする等の
// 意図しない副作用の再発火につながってしまう。
func ReconstructOrder(id OrderID, customerID CustomerID, items []OrderItem, status Status, placedAt time.Time) *Order {
	itemsCopy := make([]OrderItem, len(items))
	copy(itemsCopy, items)

	return &Order{
		id:         id,
		customerID: customerID,
		items:      itemsCopy,
		status:     status,
		placedAt:   placedAt,
	}
}

// ID は注文の識別子を返す。
func (o *Order) ID() OrderID {
	return o.id
}

// CustomerID はこの注文を行った顧客の ID を返す。
func (o *Order) CustomerID() CustomerID {
	return o.customerID
}

// Items は注文明細の一覧を返す。
// 呼び出し元が内部スライスを直接書き換えられないようコピーを返す
// （cart.Cart.Items() と同じカプセル化の理由による）。
func (o *Order) Items() []OrderItem {
	items := make([]OrderItem, len(o.items))
	copy(items, o.items)
	return items
}

// Status は注文の現在の状態を返す。
func (o *Order) Status() Status {
	return o.status
}

// PlacedAt は注文が確定した時刻を返す。
func (o *Order) PlacedAt() time.Time {
	return o.placedAt
}

// TotalAmount は全明細の小計を合算した注文全体の合計金額を返す。
// 明細が 1 件もない状態は NewOrder が禁じているため通常は発生しないが、
// ReconstructOrder 経由で空の明細を渡された場合に備え、0 円（JPY）を
// フォールバックとして返す。
func (o *Order) TotalAmount() (shared.Money, error) {
	if len(o.items) == 0 {
		return shared.NewMoney(0, shared.JPY)
	}

	total, err := shared.NewMoney(0, o.items[0].UnitPrice().Currency())
	if err != nil {
		return shared.Money{}, err
	}
	for _, item := range o.items {
		subtotal, err := item.Subtotal()
		if err != nil {
			return shared.Money{}, err
		}
		total, err = total.Add(subtotal)
		if err != nil {
			return shared.Money{}, err
		}
	}
	return total, nil
}

// 【状態遷移】
//
//	pending --Pay--> paid --Ship--> shipped
//	pending --Cancel--> canceled
//	paid    --Cancel--> canceled
//
// 上記以外の遷移（例: shipped から Cancel、pending から Ship）は
// すべてドメインルール違反として拒否する。状態機械（State Machine）の
// 判断を Order 集約に閉じ込めることで、「今どの状態からどの状態へ
// 遷移できるか」という業務ルールが Order の外（アプリケーション層や
// プレゼンテーション層）に漏れ出さないようにしている。

// Pay は注文の入金を確認し、pending から paid へ遷移させる。
func (o *Order) Pay() error {
	if o.status != StatusPending {
		return shared.NewDomainRuleError("order: cannot pay an order in status %q", o.status)
	}
	o.status = StatusPaid
	return nil
}

// Ship は注文を発送済みにし、paid から shipped へ遷移させる。
func (o *Order) Ship() error {
	if o.status != StatusPaid {
		return shared.NewDomainRuleError("order: cannot ship an order in status %q", o.status)
	}
	o.status = StatusShipped
	return nil
}

// Cancel は注文を取り消す。pending・paid のいずれからも取り消せるが、
// 発送済み（shipped）・すでに取り消し済み（canceled）の注文は取り消せない。
func (o *Order) Cancel() error {
	if o.status != StatusPending && o.status != StatusPaid {
		return shared.NewDomainRuleError("order: cannot cancel an order in status %q", o.status)
	}
	o.status = StatusCanceled
	return nil
}
