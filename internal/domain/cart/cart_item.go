package cart

// CartItem はカートに入っている「ある商品を何個持っているか」という
// 行（line item）を表すエンティティである。
//
// 集約ルート（Cart）とエンティティ（CartItem）の違い:
//   - 集約ルートは外部から直接参照・操作される唯一の入口であり、
//     集約全体の不変条件（例: 明細は最大 20 件まで）を守る責任を持つ。
//   - CartItem のような「集約内部のエンティティ」は、集約ルートを
//     経由してのみ生成・変更されるべきで、外部が直接 CartItem を
//     new して Cart に注入するようなことは許さない。
//
// そのためフィールドを非公開にし、Cart 経由（AddItem/RemoveItem）でしか
// 状態を変更できないようにしている。CartItem 自身は「値を保持するだけの
// 入れ物」であり、ビジネスルールの判断は Cart 側に集約する。
type CartItem struct {
	productID ProductID
	quantity  int
}

// ProductID はこの明細が指す商品の ID を返す。
func (i CartItem) ProductID() ProductID {
	return i.productID
}

// Quantity はこの明細の数量を返す。
func (i CartItem) Quantity() int {
	return i.quantity
}
