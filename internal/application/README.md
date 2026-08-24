# application 層

## この層の役割

`internal/application` はユースケース（アプリケーションサービス）を置く層である。
「カートに商品を追加する」「注文を確定する」といった、ユーザーの操作単位に
対応する処理をここに実装する。

ユースケース自体はビジネスルールを持たない。ビジネスルールの判断はすべて
ドメイン層のエンティティ・値オブジェクト・ドメインサービスに委譲し、
ユースケースは以下のような「調整役（オーケストレーション）」に徹する。

1. リポジトリから集約を読み込む
2. 集約のメソッドを呼んでドメインルールに従った状態変更を行わせる
3. リポジトリで永続化する
4. 必要であれば、永続化成功後に他コンテキストの集約・リポジトリを
   直接操作する（後述の `PlaceOrderUseCase` が代表例だが、
   `ShipOrderUseCase`・`CancelOrderUseCase` も同様に order 以外の
   コンテキスト（customer/shipping/inventory、inventory/coupon）の
   リポジトリを直接操作する）

## コマンドとクエリの分離（軽量 CQRS）

本サンプルでは CQRS（Command Query Responsibility Segregation）の考え方を
軽量な形で採用する。ただし catalog / cart / order のいずれのコンテキストも
コマンドとクエリを別パッケージには分けていない。`internal/application/usecase/catalog`
のようにコンテキストごとに 1 パッケージとし、その中でファイル名と依存する
型の形からコマンドかクエリかを読み分ける方針にしている
（例: `internal/application/usecase/catalog/register_product.go` は
コマンド、`internal/application/usecase/catalog/list_products.go` はクエリ）。

- **コマンド（更新系）**: 必ずドメイン層を経由する。
  「カートに商品を追加する」のようにビジネスルールに従った状態変更を
  伴う操作は、集約を読み込み・操作し・保存するという手順を踏む。
  これは整合性（不変条件）を守るために不可欠である。依存の形としては
  リポジトリ（ドメイン層のポート）と `tx.Manager` の両方に依存する。
- **クエリ（参照系）**: ドメイン層を経由せず、クエリサービスから
  直接 DTO（Data Transfer Object）を組み立てて返す。
  一覧表示や詳細表示のような「読むだけ」の操作にドメインオブジェクトの
  復元コストをかける必要はなく、SQL で必要な列だけを取得して
  そのまま DTO に詰めた方が高速かつシンプルである。依存の形としては
  QueryService（このパッケージ自身が宣言するポート）1 つにしか依存しない。

このように更新系と参照系でモデルを使い分けることで、
「ドメインモデルは複雑な整合性を守るために存在する」
「参照系は表示に必要な形をいかに効率よく組み立てるかに集中する」
という異なる関心事を無理に 1 つのモデルに押し込めずに済む。パッケージを
分割しなかったのは、両者とも「同じコンテキストの入出力窓口」という
点では同じ役割であり、独立したパッケージにするほどの分離の必要性が
なかったためである。

## トランザクション境界

1 回のユースケース実行は、原則として 1 つのトランザクションに対応する。
これは `internal/application/tx.Manager` インターフェースを介して表現する。

```go
// internal/application/usecase/cart/add_item.go を簡略化した例
func (u *AddItemUseCase) Execute(ctx context.Context, in AddItemInput) error {
    // ... 入力検証で customerID / productID を組み立てる ...
    return u.tx.Do(ctx, func(ctx context.Context) error {
        c, err := u.repo.FindByCustomerID(ctx, customerID)
        if err != nil {
            if errors.Is(err, shared.ErrNotFound) {
                c = domaincart.NewCart(customerID)
            } else {
                return err
            }
        }
        if err := c.AddItem(productID, in.Quantity); err != nil {
            return err
        }
        return u.repo.Save(ctx, c)
    })
}
```

ユースケースは `tx.Manager` というインターフェースにのみ依存し、
その実装が GORM のトランザクションであることを知らない。これにより
ユースケースのコードは「何をトランザクションに含めるべきか」という
ビジネス上の判断だけに集中でき、インフラの実装詳細から独立している。

## コンテキストをまたぐ直接呼び出し（PlaceOrderUseCase の例）

「注文を確定したらカートを空にする」というのは order コンテキストと cart
コンテキストをまたぐ処理である。以前は `Order` 集約が `OrderPlaced` という
ドメインイベントを記録し、ユースケースが `Save` 成功後にそれを配信、
cart 側のイベントハンドラが購読してカートを空にする、という疎結合な
設計を採っていた。しかし本サンプルでは仕組みの理解コストを下げるため、
`PlaceOrderUseCase` が `domaincart.Repository` を直接呼び出す方式に変更した。
以下はその実際のコードの抜粋である。

```go
// internal/application/usecase/order/place_order.go を一部抜粋
err = uc.txManager.Do(ctx, func(ctx context.Context) error {
    c, err := uc.cartRepo.FindByCustomerID(ctx, cartCustomerID)
    // ... カートを読み込み、明細ごとに catalog から現在価格を取得して
    //     OrderItem のスナップショットを組み立てる ...

    o, err := domainorder.NewOrder(orderCustomerID, items, now)

    if err := uc.orderRepo.Save(ctx, o); err != nil {
        return err
    }

    // 以前はここでドメインイベント(OrderPlaced)を発行し、カート側の
    // ハンドラが購読して空にする方式だったが、仕組みの理解コストを
    // 下げるため直接呼び出しに変更した。application 層での直接呼び出しは
    // order→cart の依存を生む(結合度は上がる)一方、処理の流れが一目で
    // 追える。コンテキストをまたぐ反応が増えてきたらドメインイベント+
    // 購読へ戻すのが定石。
    c.Clear()
    if err := uc.cartRepo.Save(ctx, c); err != nil {
        return err
    }
    // ...
    return nil
})
```

このユースケース全体は 1 つのトランザクション（`txManager.Do` の fn）の
中にあるため、`Order` の `Save` と `Cart` の `Clear` + `Save` は同一
トランザクションに参加する。途中のどこかで失敗すればトランザクション
全体がロールバックされるため、「注文は確定したがカートは空にならなかった」
のような中途半端な状態は起こり得ない。

**トレードオフ。** 直接呼び出しは order パッケージが cart パッケージを
知ってしまう（結合度が上がる）という代償を伴う。一方でドメインイベント +
購読の方式は order 側が cart の存在を一切知らずに済む疎結合な設計だが、
「どこで何が起きるか」を追うのに複数ファイルを行き来する必要があり、
学習コストが上がる。反応するコンテキストが 1 つだけの今は直接呼び出しの
単純さを優先しているが、「注文確定時にポイントを付与する」「在庫を
引き当てる」のように反応するコンテキストが増えてきたら、ドメインイベント
+ 購読による疎結合化へ戻すのが定石である。

**実際のコストは order→cart の1対1の結合にとどまらない。**
`PlaceOrderUseCase`（[`place_order.go:9-15`](usecase/order/place_order.go)）は
`domaincart` / `domaincatalog` / `domaincustomer` / `domaininventory` /
`domaincoupon` の 5 コンテキストのドメインパッケージを import しており、
1 回のトランザクション（`txManager.Do` の fn）の中で最大 4 集約種・
約 23 集約インスタンス（Cart × 1、Order × 1、Stock × カート明細数
（最大 20）、Coupon × 最大 1）を更新する。`ShipOrderUseCase`
（[`ship_order.go`](usecase/order/ship_order.go)）も同様に Order・
Shipment・Stock の 3 集約種を 1 トランザクションで更新する。

これは Vernon の集約設計ルール1「1 トランザクションで変更する集約
インスタンスは 1 つ」（[ddd-research.md](../../docs/ddd-research.md)
「集約の設計 — Vernon の4ルール」参照）から明確に逸脱している。ただし
これは無自覚な逸脱ではなく、Vernon 自身が Part II
"Reasons To Break the Rules" で挙げる4つの理由のうち **Reason Two:
Lack of Technical Mechanisms**（メッセージング等の技術的手段が使えない
場合の逸脱許容）に相当する意図的な選択である。本サンプルはドメイン
イベント + 非同期配信という技術的手段をあえて持たないと決めており
（前節「コンテキストをまたぐ直接呼び出し」参照）、その結果として
複数集約を1トランザクションにまとめる直接呼び出しを採用している。

**本番規模での帰結。** 学習用サンプルの DB 負荷では顕在化しないが、
実運用規模でこの設計をそのまま採用すると次のようなコストが発生する。

- **並行時のロストアップデート**: 行ロックを複数集約にまたがって同時に
  保持し続けるわけではなく、各保存（`Save`）は書き込み時にのみロックを
  取って解放する。真の危険は「読んでから書く」間に他トランザクションが
  割り込む lost update である。たとえば `Stock.Reserve` は
  `FindByProductID` で読み込んだ在庫スナップショットをもとに引当可能数を
  判定するため、2 つの注文が同時に同じ商品を読み込むと、どちらも
  「引当可能」と判定してしまい実在庫を超えて引き当ててしまう可能性がある
  （クーポンの利用回数上限も同様）。防ぐには `SELECT ... FOR UPDATE`
  （gorm の `clause.Locking`）や楽観ロック（version カラム）が必要だが、
  本サンプルは実装していない（`inventory.Stock.Reserve` のコメントを参照）。
- **ロールバック範囲**: `PlaceOrderUseCase` は 6 つのリポジトリ
  （order/cart/catalog/customer/inventory/coupon）を注入されており、
  そのうち実際に保存（`Save`）を行うのは order/cart/inventory/coupon の
  4 つである。これらのいずれか 1 つの保存が失敗しても、トランザクション
  全体（Cart のクリアも含む）がロールバックされる。失敗の原因（在庫不足・
  クーポン無効・DB 接続断など）によらず影響範囲が一律に広い。
- **段階的分離の難しさ**: 将来 cart・inventory・coupon を別サービス
  （別 DB）に分離しようとすると、この 1 トランザクション前提の実装は
  そのままでは成立しなくなる。分散トランザクションや Saga パターンへの
  書き換えが必要になり、単一 DB スキーマを前提にした現在の設計からの
  移行コストが高い。
