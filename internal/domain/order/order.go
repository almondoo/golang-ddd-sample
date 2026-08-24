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
type Order struct {
	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     Status
	placedAt   time.Time
	// couponCode は適用済みクーポンのコードである。空文字列は「クーポン未適用」を表す。
	// discountAmount と合わせて ApplyDiscount 経由でのみ設定される。
	couponCode string
	// discountAmount は適用済みの割引額である。クーポン未適用時はゼロ値の
	// Money（amount=0, currency未設定）のままとなる。
	discountAmount shared.Money
}

// NewOrder は新規注文を確定するコンストラクタである。
//
// 呼び出し時点で明細が 1 件もない注文は業務的に意味を持たない
// （「何も買わない注文」は存在しない）ため、items の非空をここで検証する。
// 個々の OrderItem 自体の妥当性（数量 1 以上等）は NewOrderItem
// コンストラクタがすでに保証済みという前提に立つ。
//
// 生成と同時に StatusPending をセットする。
//
// 【時刻もまた入力】
// ドメイン層では time.Now() を呼ばず、時刻を引数で受けることで決定的な
// テストを可能にする（coupon.Coupon.Use / IsExpired と同じ方針）。もし
// 内部で time.Now() を呼んでしまうと、PlacedAt の値がテスト実行時刻に
// 依存してしまい、再現性のあるテストが書けなくなる。呼び出し側
// （PlaceOrderUseCase）が now を注入する。
func NewOrder(customerID CustomerID, items []OrderItem, now time.Time) (*Order, error) {
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
		placedAt:   now,
	}
	return o, nil
}

// ReconstructOrder は永続化層から読み込んだデータをもとに Order を再構築する。
//
// NewOrder との違いは ReconstructProduct / ReconstructCart と同じ理由による。
// couponCode が空文字列の場合はクーポン未適用として扱い、discountAmount は
// そのまま素通しする（永続化層側で「未適用ならゼロ値の Money」を組み立てる
// ことを前提とする）。
func ReconstructOrder(id OrderID, customerID CustomerID, items []OrderItem, status Status, placedAt time.Time, couponCode string, discountAmount shared.Money) *Order {
	itemsCopy := make([]OrderItem, len(items))
	copy(itemsCopy, items)

	return &Order{
		id:             id,
		customerID:     customerID,
		items:          itemsCopy,
		status:         status,
		placedAt:       placedAt,
		couponCode:     couponCode,
		discountAmount: discountAmount,
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

// CouponCode は適用済みクーポンのコードを返す。空文字列は未適用を表す。
func (o *Order) CouponCode() string {
	return o.couponCode
}

// DiscountAmount は適用済みの割引額を返す。クーポン未適用の場合はゼロ値の
// Money（amount=0）が返る。
func (o *Order) DiscountAmount() shared.Money {
	return o.discountAmount
}

// PayableAmount は実際に支払うべき金額（TotalAmount - DiscountAmount）を返す。
// クーポン未適用の場合は discountAmount がゼロ円のままなので、
// 実質的に TotalAmount と同額になる。
func (o *Order) PayableAmount() (shared.Money, error) {
	total, err := o.TotalAmount()
	if err != nil {
		return shared.Money{}, err
	}
	if o.couponCode == "" {
		// 未適用時は discountAmount の通貨単位が未設定（ゼロ値）のままであり、
		// そのまま total.Subtract(discountAmount) を呼ぶと通貨不一致で
		// エラーになってしまう。そのため未適用時は total をそのまま返す。
		return total, nil
	}
	return total.Subtract(o.discountAmount)
}

// ApplyDiscount はクーポンによる割引を注文に適用する。
//
// 【割引の妥当性は Order 集約自身が守る】
// 「どのクーポンが有効か」「割引額をいくら計算するか」はクーポンコンテキスト
// （coupon.Coupon）の責務だが、「その割引を注文に適用してよいか」という
// 判断（注文の状態・二重適用・金額の整合性）は注文コンテキストの不変条件
// であり、Order 集約自身が守るべきものである。そのため、割引額の計算結果を
// 受け取った後の適用可否チェックはすべてこのメソッドに閉じ込める。
func (o *Order) ApplyDiscount(couponCode string, discount shared.Money) error {
	if o.status != StatusPending {
		return shared.NewDomainRuleError("order: cannot apply discount to an order in status %q", o.status)
	}
	if o.couponCode != "" {
		return shared.NewDomainRuleError("order: coupon %q is already applied to this order", o.couponCode)
	}
	if couponCode == "" {
		return shared.NewDomainRuleError("order: coupon code must not be empty")
	}

	total, err := o.TotalAmount()
	if err != nil {
		return err
	}
	// 割引額が合計金額を超える場合（通貨不一致を含む）はドメインルール違反
	// として拒否する。Money.Subtract は「通貨不一致」と「結果が負になる」の
	// 両方をすでに検証してくれるため、ここではその結果を借りて「割引後に
	// 支払うべき金額」を検証する。戻り値そのものは使わず、検証のためだけに
	// 呼び出している点に注意する（実際の PayableAmount 計算は別メソッドで
	// 都度行う）。
	if _, err := total.Subtract(discount); err != nil {
		return shared.NewDomainRuleError("order: discount %v must not exceed total amount %v: %v", discount, total, err)
	}

	o.couponCode = couponCode
	o.discountAmount = discount
	return nil
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
