package cart

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// maxQuantityPerItem は 1 明細（同一商品）が持てる数量の上限である。
// カートの数量には「注文として現実的な範囲」という業務ルールがあるため、
// 無制限にはせず定数として明示する。
const maxQuantityPerItem = 99

// maxDistinctItems はカートが持てる明細（商品種別）の上限である。
// カートが無限に肥大化するのを防ぎ、注文確定時の処理コストも抑える。
const maxDistinctItems = 20

// Cart はショッピングカートを表す集約ルート（Aggregate Root）である。
//
// 集約の識別子には CustomerID をそのまま使っている（サロゲートの CartID を
// 別途発行していない）。これは「1 顧客につきカートは常にちょうど 1 つ」
// という業務ルールをそのままモデルに反映させた設計判断である。
// もし将来「同一顧客が複数のカートを持てる」（例: 保存済みカートの複数管理）
// という要件が生まれたら、そのときは CustomerID とは別に CartID を発行する
// 設計へ切り替える必要がある。今の要件に対して過剰な柔軟性を先取りしない、
// という YAGNI の実践例でもある。
//
// 価格を一切持たない設計について:
// カートは「何をいくつ欲しいか」という意思表示だけを保持し、金額計算は
// 行わない。価格は注文確定（Order コンテキスト）のタイミングで catalog の
// 最新価格を参照して初めて確定させる。こうすることで、
//   - カートに入れてから購入するまでの間に商品価格が変わっても、
//     カート側は何も知らなくてよい（追随して更新する必要がない）
//   - 「カートに入れた時点の価格」と「注文時の価格」のどちらを正とするか、
//     という曖昧な業務ルールの決定をカートコンテキストが背負わずに済む
//
// というメリットがある。参考画面（クエリ側）で金額を表示したい場合は、
// query_service.go のようにカート明細と商品テーブルを SQL レベルで
// 結合して求める。
type Cart struct {
	customerID CustomerID
	items      []CartItem
}

// NewCart は空のカートを新規生成する。
func NewCart(customerID CustomerID) *Cart {
	return &Cart{
		customerID: customerID,
		items:      nil,
	}
}

// ReconstructCart は永続化層から読み込んだデータをもとに Cart を再構築する。
//
// NewCart との違い: NewCart は「新規カートを作る」というドメイン上の意図を
// 表すのに対し、ReconstructCart は「すでに存在するカートを DB から復元する」
// という別の意図を表す。両者を同じコンストラクタにまとめてしまうと、
// 「新規作成なのか復元なのか」を呼び出し側が区別できなくなり、
// 誤って初期状態を上書きしてしまうバグの温床になる。
// リポジトリ実装（infrastructure 層）からのみ呼ばれることを想定している。
func ReconstructCart(customerID CustomerID, items []CartItem) *Cart {
	return &Cart{
		customerID: customerID,
		items:      items,
	}
}

// ReconstructCartItem は永続化層のレコードから CartItem を再構築する。
//
// 通常 CartItem は Cart.AddItem を経由してのみ生成され、その過程で
// 数量の不変条件（1〜99）がチェックされる。しかし DB から読み込む際は
// 「すでに過去に検証済みのデータをそのまま復元する」だけなので、
// 再度検証を強制する必要はない（むしろ検証ロジックを二重に持つと、
// 将来ルールを変えたときに過去データの復元が失敗するリスクがある）。
// そのため ReconstructCartItem は検証を行わない「素通し」の関数として
// 分離し、AddItem 経由の正規ルートとは明確に区別している。
func ReconstructCartItem(productID ProductID, quantity int) CartItem {
	return CartItem{productID: productID, quantity: quantity}
}

// CustomerID はこのカートの持ち主を返す。
func (c *Cart) CustomerID() CustomerID {
	return c.customerID
}

// Items はカート内の明細一覧を返す。
//
// 呼び出し元がスライスを直接書き換えて集約の不変条件を壊せないよう、
// コピーを返す。内部スライスへの参照をそのまま渡してしまうと、
// 「集約の外から集約の内部状態を直接改変できてしまう」という
// カプセル化違反になる。
func (c *Cart) Items() []CartItem {
	items := make([]CartItem, len(c.items))
	copy(items, c.items)
	return items
}

// IsEmpty はカートに明細が 1 つもないかどうかを返す。
func (c *Cart) IsEmpty() bool {
	return len(c.items) == 0
}

// AddItem はカートに商品を追加する。
//
// 同一商品がすでにカートにある場合は、新しい明細を増やすのではなく
// 既存の明細に数量をマージする（「同じ商品を 2 回に分けて追加する」操作は
// 業務的には「合計数量が増える」という 1 つの事実にすぎない）。
//
// 注意（意図的な設計判断）: この時点では商品が catalog コンテキストに
// 実在するかどうかを検証しない。カートコンテキストは catalog に依存しない
// ため、検証しようとすると catalog への同期呼び出しが必要になり、
// コンテキスト間が強く結合してしまう。存在検証は「実際に注文する」という
// より重い操作（注文確定ユースケース）のタイミングまで遅延させる。
// これにより、カートへの追加は catalog の可用性に左右されない、
// 疎結合で軽量な操作であり続けられる（トレードオフとして、
// 実在しない商品 ID が一時的にカートに混入する可能性はある）。
func (c *Cart) AddItem(productID ProductID, quantity int) error {
	if quantity < 1 || quantity > maxQuantityPerItem {
		return shared.NewDomainRuleError("cart: quantity must be between 1 and %d, got %d", maxQuantityPerItem, quantity)
	}

	for i, item := range c.items {
		if item.productID == productID {
			merged := item.quantity + quantity
			if merged > maxQuantityPerItem {
				return shared.NewDomainRuleError("cart: merged quantity must not exceed %d, got %d", maxQuantityPerItem, merged)
			}
			c.items[i].quantity = merged
			return nil
		}
	}

	if len(c.items) >= maxDistinctItems {
		return shared.NewDomainRuleError("cart: cart must not have more than %d distinct items", maxDistinctItems)
	}

	c.items = append(c.items, CartItem{productID: productID, quantity: quantity})
	return nil
}

// RemoveItem はカートから指定商品の明細を取り除く。
// カートに存在しない商品を指定した場合はドメインルール違反として扱う。
// これは「存在しないものを消そうとする」呼び出し側の誤りをサイレントに
// 無視せず、明示的に気づかせるための設計判断である。
func (c *Cart) RemoveItem(productID ProductID) error {
	for i, item := range c.items {
		if item.productID == productID {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return nil
		}
	}
	return shared.NewDomainRuleError("cart: product %s is not in the cart", productID)
}

// Clear はカート内の明細をすべて取り除く。
// 注文確定後にカートを空にする、といった用途を想定している。
func (c *Cart) Clear() {
	c.items = nil
}
