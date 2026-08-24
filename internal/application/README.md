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
4. 永続化が成功したら、集約から `PullEvents()` でイベントを取り出し配信する
   （実際にイベントを扱うのは後述の `PlaceOrderUseCase` のみである）

## コマンドとクエリの分離（軽量 CQRS）

本サンプルでは CQRS（Command Query Responsibility Segregation）の考え方を
軽量な形で採用する。

- **コマンド（更新系）**: 必ずドメイン層を経由する。
  「カートに商品を追加する」のようにビジネスルールに従った状態変更を
  伴う操作は、集約を読み込み・操作し・保存するという手順を踏む。
  これは整合性（不変条件）を守るために不可欠である。
- **クエリ（参照系）**: ドメイン層を経由せず、リポジトリやクエリサービスから
  直接 DTO（Data Transfer Object）を組み立てて返す。
  一覧表示や詳細表示のような「読むだけ」の操作にドメインオブジェクトの
  復元コストをかける必要はなく、SQL で必要な列だけを取得して
  そのまま DTO に詰めた方が高速かつシンプルである。

このように更新系と参照系でモデルを使い分けることで、
「ドメインモデルは複雑な整合性を守るために存在する」
「参照系は表示に必要な形をいかに効率よく組み立てるかに集中する」
という異なる関心事を無理に 1 つのモデルに押し込めずに済む。

## トランザクション境界

1 回のユースケース実行は、原則として 1 つのトランザクションに対応する。
これは `internal/application/tx.Manager` インターフェースを介して表現する。

```go
// internal/application/cart/command/add_item.go を簡略化した例
func (u *AddItemUseCase) Execute(ctx context.Context, in AddItemInput) error {
    // ... 入力検証で customerID / productID を組み立てる ...
    return u.tx.Do(ctx, func(ctx context.Context) error {
        c, err := u.repo.FindByCustomerID(ctx, customerID)
        if err != nil {
            if errors.Is(err, shared.ErrNotFound) {
                c = cart.NewCart(customerID)
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

## イベント配信は Save 成功後・同一トランザクション内で行う

ドメインイベントの配信は `txManager.Do` の **内側**、しかも
リポジトリの `Save` が成功した **後** に行う。これは一見奇妙に見えるかもしれない
（コミットされる前にイベントを配信してしまってよいのか、という疑問が浮かぶ）が、
本サンプルでは意図的にこの順序を採っている。

理由は、イベント配信を担う `internal/infrastructure/event.Bus` が
「`Publish` が呼ばれた時点で該当ハンドラーを同期的に（同じ goroutine・同じ ctx で）
呼び出す」インメモリ実装であることにある。ハンドラー内で行われる DB 操作も
`ctx` 経由で **同じトランザクション** に参加するため、`Save` → `Publish`
→ ハンドラーの処理 → （すべて成功すれば）トランザクションのコミット、という
一連の流れが 1 つの原子的な単位になる。途中のどこかで失敗すればトランザクション
全体がロールバックされるため、「注文は確定したがカートは空にならなかった」
のような中途半端な状態は起こり得ない。同期・同一トランザクションという
方式は、整合性を非常にシンプルな形で得られるのが利点である。

このユースケース群のうち、実際にイベントを配信しているのは現時点では
`PlaceOrderUseCase`（`internal/application/order/command/place_order.go`）
だけである。`Cart` / `Product` はイベントを記録しない（`shared.AggregateBase`
を埋め込んでいるのは `Order` のみ）。以下はその実際のコードの抜粋である。

```go
// internal/application/order/command/place_order.go を一部抜粋
err = uc.txManager.Do(ctx, func(ctx context.Context) error {
    c, err := uc.cartRepo.FindByCustomerID(ctx, cartCustomerID)
    // ... カートを読み込み、明細ごとに catalog から現在価格を取得して
    //     OrderItem のスナップショットを組み立てる ...

    o, err := order.NewOrder(orderCustomerID, items)
    // この時点で Order 集約の内部に OrderPlaced が記録される

    if err := uc.orderRepo.Save(ctx, o); err != nil {
        return err
    }

    // イベントの配信は Save が成功した「後」に行う。Save より前や
    // Save と同時に配信してしまうと、Save が失敗した場合に「実際には
    // 起きなかった注文確定」を配信してしまう不整合が生じる。
    //
    // このユースケース全体が 1 つのトランザクション（txManager.Do の fn）の
    // 中にあるため、ここで呼ばれる Publish（同期のインメモリバス経由で
    // ClearCartOnOrderPlaced.Handle を呼ぶ）もまだコミット前の同一
    // トランザクションに参加する。
    if err := uc.publisher.Publish(ctx, o.PullEvents()...); err != nil {
        return err
    }
    // ...
    return nil
})
```

ただし、この方式には限界もある。ハンドラーが外部システム（メール送信・
外部 API 呼び出し等）を呼ぶ場合、その呼び出し自体は DB トランザクションの
外にある「取り消せない副作用」になってしまい、上記の「1 つの原子的な単位」
という前提が崩れる。また、購読者が増えるほど 1 トランザクションが長くなり、
非同期化やスケールアウトもしづらくなる。外部システム連携が必要になったり
配信を非同期化したくなったりした場合は、Outbox パターン（コミット後に
別プロセスが配信を担う方式）への切り替えを検討する。詳細は
`internal/infrastructure/README.md` を参照。
