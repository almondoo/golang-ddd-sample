# リポジトリ(Repository)

リポジトリとは、集約の永続化を抽象化するポート(インターフェース)です。インターフェースはドメイン層に置き、具体的な実装(SQL発行など)はインフラストラクチャ層に置くことで、依存の向きを逆転させます。ドメイン層は「どう保存されるか」を一切知らずに「保存できる」という契約にだけ依存します。

## なぜ必要か

ドメイン層がDB操作の詳細(GORM、SQL)に直接依存すると、ビジネスルールの表現が技術的関心事に引きずられ、DBやORMを差し替える際にドメイン層まで書き換えが必要になります。インターフェースをドメイン層側で定義し、実装をインフラ層に置く「依存性逆転の原則」によって、この結合を断ち切ります。

## 依存性逆転のイメージ

```mermaid
flowchart TB
    subgraph Domain["domain/order(最内層)"]
        IF["Repository(interface)<br>FindByID / Save"]
    end
    subgraph Application["application/usecase/order"]
        UC["PlaceOrderUseCase"]
    end
    subgraph Infrastructure["infrastructure/persistence(最外層)"]
        Impl["OrderRepository(GORM実装)"]
    end
    UC -->|依存(ポートを呼ぶ)| IF
    Impl -->|依存(インターフェースを実装)| IF
    UC -.->|実行時はDIで注入されたImplを呼ぶ| Impl
```

依存の矢印(コンパイル時のimport)はimplからinterfaceへ向かいますが、実行時の呼び出しはusecaseからimplの実装に対して行われます。この非対称性は[onion-architecture.md](onion-architecture.md)でも扱います。

## このリポジトリでの実例

[internal/domain/order/repository.go](../../internal/domain/order/repository.go)がインターフェースを定義し、[internal/infrastructure/persistence/order_repository.go](../../internal/infrastructure/persistence/order_repository.go)がGORMで実装します。実装ファイルの`var _ order.Repository = (*OrderRepository)(nil)`という1行が、コンパイル時に「この実装がインターフェースを満たしている」ことを保証します。

`Repository`はCart単位・Order単位のように必ず集約単位でメソッドを揃えており、`CartItem`や`OrderItem`だけを個別に取得・保存するAPIは存在しません([cart/repository.go](../../internal/domain/cart/repository.go)のコメント参照)。集約の外から内部エンティティを個別に触れる経路を作らないためです。

Evansの原典は「リポジトリは集約ルートにのみ提供する」「トランザクション制御はリポジトリではなくクライアント(アプリケーション層)が持つ」と述べています(逐語引用は[ddd-research.md](../ddd-research.md)を参照)。本リポジトリのリポジトリ実装は`Save`の中で`db.Transaction`を自ら開始せず、`DBFromContext`で現在のトランザクション用DBハンドルを取得するだけに留めており、この原則に沿っています。トランザクション境界の実際の制御は[application-service.md](application-service.md)を参照してください。

## 注意点・よくある誤解

- 「1集約=1リポジトリ」は本リポジトリでは実質そうなっていますが、Evansの原典自体はこの1:1対応を要求していません([ddd-research.md](../ddd-research.md)の「検証で棄却・訂正された主張」参照)。
- リポジトリの`Save`は新規作成・更新の両方を担う upsert として実装されており、`Create`/`Update`のようなメソッド分割はしていません。
