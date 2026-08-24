# 境界づけられたコンテキスト(Bounded Context)

境界づけられたコンテキストとは、ある1つのドメインモデルが一貫した意味を持って通用する範囲のことです。Fowlerは、境界はデプロイ単位ではなくモデルの適用範囲で引かれるとし、物理DBの分離自体は要件ではないとしています(詳細は[ddd-research.md](../ddd-research.md))。

## コンテキスト概観

```mermaid
flowchart LR
    subgraph Contexts["7 bounded context"]
        Catalog["catalog"]
        Cart["cart"]
        Order["order"]
        Customer["customer"]
        Inventory["inventory"]
        Shipping["shipping"]
        Coupon["coupon"]
    end
    Shared["shared(共有カーネル)"]

    Order -->|application層オーケストレーション| Customer
    Order -->|application層オーケストレーション| Cart
    Order -->|application層オーケストレーション| Catalog
    Order -->|application層オーケストレーション| Inventory
    Order -->|application層オーケストレーション| Coupon
    Order -->|application層オーケストレーション(発送時)| Shipping
    Cart -.->|読み取り側 SQL JOIN(スキーマ結合)| Catalog

    Catalog --> Shared
    Cart --> Shared
    Order --> Shared
    Customer --> Shared
    Inventory --> Shared
    Shipping --> Shared
    Coupon --> Shared
```

実線はapplication層による直接呼び出し(ドメインパッケージ自体はimportしません)、点線はクエリサービスが読み取り専用で行うSQLレベルの結合です。

## このリポジトリの7コンテキスト + shared

`internal/domain`配下に、それぞれ独立したGoパッケージとして次の7つのコンテキストがあります。catalog(商品カタログ)・cart(カート)・order(注文)・customer(顧客)・inventory(在庫)・shipping(配送)・coupon(クーポン)。加えて`internal/domain/shared`が全コンテキストで共有される[shared-kernel.md](shared-kernel.md)として存在しますが、これはbounded contextそのものではなく共有カーネルという別の概念です。

## 自律性: IDの重複定義

各コンテキストは他コンテキストの内部モデルに依存すべきではないというDDDの原則に従い、同じ実体(例: 商品)を指すIDでも`catalog.ProductID` / `cart.ProductID` / `order.OrderItem.productID`のように意図的に型を重複定義しています。[cart/ids.go](../../internal/domain/cart/ids.go)のコメントがこの判断を明確に述べています。

> 値としては同じ文字列(UUID)を共有するが、境界づけられたコンテキスト(Bounded Context)同士は互いの内部モデルに依存すべきではない、という DDD の原則に従い、あえて「ID という小さな型を重複させる」ことを選んでいる。コンテキスト間の結合を避けるためのコストとしては、この程度の重複はごく小さい。

`order`ドメインパッケージが`cart`や`catalog`を一切importしないことも同じ自律性の表れです([order.go](../../internal/domain/order/order.go)冒頭コメント参照)。コンテキストをまたぐ調整は、ドメイン層ではなくアプリケーション層(`PlaceOrderUseCase`)の責務です。

## 読み取り側に残るスキーマ結合

書き込み側はこの自律性を型で守っていますが、読み取り側(クエリサービス)には結合が残ります。[cart_query_service.go](../../internal/infrastructure/persistence/cart_query_service.go)は`cart_items`と`products`(catalogコンテキスト所有)をSQLレベルで直接JOINし、`products.name`・`products.price_amount`という物理カラム名をハードコードしています。catalog側のテーブル変更はコンパイルエラーを出さずにこのJOINを実行時に破壊する可能性があります。コード中のコメントと[cart/README.md](../../internal/domain/cart/README.md)は現在この点を「書き込み側は型で分離されるが、読み取り側にはスキーマ結合が残る」と正確に記述しています(かつての「疎結合を破らない」という過大な表現は[specs/ddd-improvements.md](../specs/ddd-improvements.md)項目8として修正済み)。

## 注意点・よくある誤解

- 「7コンテキスト」という数え方はsharedを含めていません(sharedはbounded contextではなくShared Kernel)。この前提は[../context-map.md](../context-map.md)で明文化されています。
- コンテキスト間の関係パターンの分類・ID型対照表・ユースケース別のコンテキスト横断一覧は、コンテキストマップ[../context-map.md](../context-map.md)を参照してください([specs/ddd-improvements.md](../specs/ddd-improvements.md)項目6で作成)。
