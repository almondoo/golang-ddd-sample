# 実行順序の図解

リクエストが各層をどの順序で通るかを図で示します。オニオンアーキテクチャでは「依存の向き」と「実行時の呼び出し順」が別物であることが分かりにくいため、まず依存関係を確認してから、command / query それぞれの実行順序を追います。

## 層の依存関係

依存はつねに外側から内側(ドメイン)へ向かいます。ドメイン層はどの層にも依存しません。infrastructure がドメインに依存する(リポジトリIFを実装する)方向であって、その逆ではない点がオニオンアーキテクチャの核心です。

```mermaid
flowchart TD
    P["presentation/controller<br>(コントローラ・ルーティング)"]
    A["application<br>(ユースケース・command/query)"]
    D["domain<br>(集約・値オブジェクト・リポジトリIF)"]
    I["infrastructure<br>(GORM実装)"]

    P --> A
    P -. "エラー分類のみ<br>(domain/shared)" .-> D
    A --> D
    I -- "インターフェースを実装" --> D
    I -- "tx.Manager を実装" --> A
```

なお、presentation 層は基本的に application 層だけを呼びますが、エラーを HTTP ステータスへ分類するために `domain/shared` のエラー種別だけは直接参照します(点線)。外側の層が内側の層に依存すること自体はオニオンアーキテクチャのルール上問題ありません — 禁止されているのは内側から外側への依存です。

## 起動時の組み立て順序(cmd/api)

実行時に内側の層が外側を呼べるのは、起動時に依存が注入(DI)されているためです。この組み立ては `github.com/almondoo/wire`(コンパイル時 DI コード生成ツール)が生成しています — `wire.go` が「どのプロバイダをどう組み合わせるか」を書いた設計図、`wire_gen.go` がそこから生成された実際の組み立てコード(`initializeServer` 関数)です。生成コマンドは以下の通りです。

```sh
go run github.com/almondoo/wire/cmd/wire ./cmd/api
```

`main.go` はもはや依存を手で組み立てず、`initializeServer(dsn)` を呼んで `*http.ServeMux` を受け取るだけになりました。生成されたコードがたどる組み立て順序自体は、手書き DI だった頃と変わらず内側の実装から順に行われます。

1. DB 接続を開く(`persistence.NewDB`)
2. リポジトリ実装・クエリサービス実装を生成(GORM 依存)
3. ユースケースを生成(リポジトリIF・`tx.Manager` を注入)
4. HTTP コントローラを生成し、ルーティングに登録
5. サーバー起動

## command の実行順序 — 注文確定(POST /orders)

もっとも登場人物が多い「注文確定」を例にします。顧客の実在を確認し、カートを読み、商品の現在価格を引きながら在庫を引き当て、注文集約を生成し、クーポンが指定されていれば割引を適用してから保存し、最後にカートを空にする、という一連の流れです。

```mermaid
sequenceDiagram
    autonumber
    actor C as クライアント
    participant H as OrderController<br>(presentation)
    participant UC as PlaceOrderUseCase<br>(application/usecase/order)
    participant TX as TxManager<br>(infrastructure)
    participant CUR as CustomerRepository<br>(IFはdomain/実装はinfra)
    participant CR as CartRepository<br>(IFはdomain/実装はinfra)
    participant PR as ProductRepository<br>(IFはdomain/実装はinfra)
    participant SR as StockRepository<br>(IFはdomain/実装はinfra)
    participant O as Order集約<br>(domain/order)
    participant COR as CouponRepository<br>(IFはdomain/実装はinfra)
    participant OR as OrderRepository<br>(IFはdomain/実装はinfra)

    C->>H: POST /orders {customerId, couponCode?}
    H->>UC: Execute(ctx, input)
    UC->>TX: Do(ctx, fn) — トランザクション開始
    activate TX
    TX->>UC: fn(トランザクション入りctx)
    UC->>CUR: FindByID(ctx, customerID) — 顧客検証
    CUR-->>UC: Customer(存在しなければルール違反)
    UC->>CR: FindByCustomerID(ctx, customerID)
    CR-->>UC: Cart(明細つき)
    UC->>PR: FindByID(ctx, productID) × 明細ぶん
    PR-->>UC: Product(現在価格)
    UC->>SR: FindByProductID → Reserve → Save × 明細ぶん — 在庫引当
    SR-->>UC: 引当済み Stock
    UC->>O: NewOrder(customerID, items)
    Note over O: 不変条件を検証
    O-->>UC: Order
    opt couponCode が指定されている場合
        UC->>COR: FindByCode(ctx, code)
        COR-->>UC: Coupon
        UC->>COR: Use(now) → DiscountFor(total) → Save — クーポン適用
        UC->>O: ApplyDiscount(code, discount)
        Note over O: pending・未適用・合計以内であることを検証
    end
    UC->>OR: Save(ctx, order)
    UC->>CR: カートを Clear して Save
    TX-->>UC: コミット
    deactivate TX
    UC-->>H: 出力DTO(注文ID)
    H-->>C: 201 Created (JSON)
```

補足を3点。

- 顧客検証・在庫引当・クーポン適用はいずれも「注文コンテキスト自身の不変条件ではないが、注文確定という業務操作の一部として必要な調整」であり、`order` ドメインパッケージが `customer` / `inventory` / `coupon` を import することなく、アプリケーション層(`PlaceOrderUseCase`)が各コンテキストのリポジトリを横断して呼び出す形で実現しています。
- カートのクリアはユースケース(application 層)から `cart.Repository` への直接呼び出しであり、`Order` の `Save` と同一トランザクション内で行われます。「注文は確定したがカートが空にならない」という不整合が起きない一方、order→cart という依存(結合度の上昇)を受け入れる判断です。
- コンテキストをまたぐ反応が増えてきたら(例: 注文確定時にポイント付与も行いたい等)、ドメインイベント + 購読による疎結合な連携へ発展させるのが定石です。

なお、発送(`POST /orders/{id}/ship`)も同様に複数コンテキストにまたがる操作です。`ShipOrderUseCase` は `Order.Ship()` で状態を進めた後、顧客のデフォルト住所から `shipping.Shipment` を生成して発送指示と同時に出荷済みにし、注文明細ごとに `inventory.Stock.ConsumeReserved` で引当済みの在庫を実際に払い出す、という「Shipment 生成 + 在庫消込」を 1 トランザクションでまとめて行います。

## query の実行順序 — 商品一覧(GET /products)

query はドメイン層を通りません。クエリサービス(application 層でIF定義・infrastructure 層で GORM 実装)が DB から直接 DTO を組み立てます。集約の不変条件は書き込み時に守られているため、読み取りで再構築するコストを払わない、という軽量 CQRS の判断です。

```mermaid
sequenceDiagram
    autonumber
    actor C as クライアント
    participant H as ProductController<br>(presentation)
    participant Q as ListProductsUseCase<br>(application/usecase/catalog)
    participant QS as ProductQueryService<br>(IFはapplication/実装はinfra)
    participant DB as PostgreSQL

    C->>H: GET /products
    H->>Q: Execute(ctx)
    Q->>QS: List(ctx)
    QS->>DB: SELECT(GORM)
    DB-->>QS: 行データ
    QS-->>Q: []ProductDTO(ドメインを経由しない)
    Q-->>H: []ProductDTO
    H-->>C: 200 OK (JSON)
```

## エラーの変換順序

ドメイン層が返すエラーは、presentation 層で HTTP ステータスに変換されます。ドメイン層は HTTP を知らないため、エラーの「種類」だけを表現します。

| ドメイン層のエラー | 意味 | HTTP ステータス |
|---|---|---|
| `shared.ErrNotFound`(ラップされる) | 対象が存在しない | 404 |
| `shared.NewDomainRuleError(...)` | ビジネスルール違反(例: 在庫上限超過、不正な状態遷移) | 422 |
| 上記以外 | 想定外の失敗 | 500 |
