# inventory コンテキスト

在庫（在庫数の管理と引当）を扱う境界づけられたコンテキスト（Bounded Context）である。

## 集約境界

在庫は「商品ごとに 1 つの在庫集約（`Stock`）」という粒度でモデル化する。
`Stock` の識別子は商品を示す `ProductID` をそのまま使い、サロゲートの
`StockID` は発行しない。「1 商品につき在庫レコードは常にちょうど 1 つ」
という業務ルールをそのままモデルに反映させた設計判断であり、`cart.Cart`
が `CustomerID` をそのまま識別子に使っているのと同じ考え方である。

`inventory.ProductID` は `catalog.ProductID` とは別の型として定義しており、
inventory コンテキストは catalog パッケージを一切 import しない。値としては
同じ文字列（UUID）を共有するが、境界づけられたコンテキスト同士は互いの
内部モデルに依存すべきではない、という DDD の原則に従っている。

## 引当モデル（quantity / reserved / available）

`Stock` は 3 つの数量概念を持つ。

- `quantity`（実在庫）: 倉庫に実際にある数量。
- `reserved`（引当済み）: 注文確定などにより「取り置き」されている数量。
  まだ出荷（実在庫の減少）はされていない。
- `available`（引当可能残数、`quantity - reserved`）: 新たに引当できる残り数量。

常に `0 <= reserved <= quantity` という不変条件が成り立つ。この不変条件は
`Stock` 集約自身が `Reserve` / `Release` / `ConsumeReserved` / `SetQuantity`
の各操作の中で守り続ける。集約の外から `reserved` や `quantity` を
直接書き換える経路は存在しない。

## 注文フローとの関係

`Stock` はそれ自体では注文を知らない（inventory コンテキストは order
コンテキストに依存しない）。実際の在庫操作は、application 層の注文
ユースケースが `inventory.Repository` を通じて `Stock` 集約を操作する形で
行われる。

- 注文確定（place）: [`PlaceOrderUseCase`](../../application/usecase/order/place_order.go)
  が明細ごとに `Stock.Reserve(n)` を呼び、実在庫を減らさずに「取り置き」を行う。
- 注文キャンセル（cancel）: [`CancelOrderUseCase`](../../application/usecase/order/cancel_order.go)
  が明細ごとに `Stock.Release(n)` を呼び、取り置きを解除する。実在庫（quantity）は変えない。
- 出荷（ship）: [`ShipOrderUseCase`](../../application/usecase/order/ship_order.go)
  が明細ごとに `Stock.ConsumeReserved(n)` を呼び、取り置きしていた分を
  実際に払い出す。このとき `quantity` と `reserved` の両方が減る。

このように「引当」と「消込・解放」を分離しているのは、注文確定と出荷の
間にはタイムラグがあり、その間も実在庫を正しく把握し続ける必要がある
という業務要件を反映したものである。
