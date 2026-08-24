# infrastructure 層

## この層の役割

`internal/infrastructure` は、ドメイン層・アプリケーション層が定義した
インターフェース（ポート）に対する具体的な実装（アダプタ）を置く層である。
DB、外部 API、メッセージング基盤といった「技術的な詳細」はすべてここに閉じ込める。

## GORM モデルとドメインモデルの分離（データマッパー）

ドメイン層のエンティティ・値オブジェクトには GORM のタグを付けない。
代わりに、`internal/infrastructure/persistence` という 1 つのフラットな
パッケージの中に、コンテキストごとにファイルを分けて GORM 用の永続化モデルと
リポジトリ実装を置く（例: `cart_model.go` / `cart_repository.go`、
`catalog_model.go` / `catalog_repository.go`、`order_model.go` /
`order_repository.go`）。永続化モデルは `CartItemModel` / `ProductModel` /
`OrderModel` / `OrderItemModel` のようにエクスポートされた構造体として
`*_model.go` に定義し、対応する `*_repository.go` がそれを使って以下の変換を行う。

- 読み込み時: `CartItemModel` 等 → ドメインの `Cart` 集約
- 保存時: ドメインの `Cart` 集約 → `CartItemModel` 等

この往復変換を担うのがリポジトリ実装の責務であり、これによって
「テーブル構造の都合」と「ドメインモデルの都合」を独立して進化させられる。
例えば、正規化されたテーブル構成であっても、ドメイン側では 1 つの値オブジェクトに
まとめて表現する、といったことが可能になる。

同じ `persistence` パッケージには、各コンテキストのクエリサービス実装
（`cart_query_service.go` / `catalog_query_service.go` /
`order_query_service.go`）も置かれている。こちらは「コマンドとクエリの分離
（軽量 CQRS）」の項（`internal/application/README.md` 参照）で説明した通り
ドメインモデルを経由しないため、上記のデータマッパーとは別物である。
DB の行から直接 DTO を組み立てて返す。

なお、パッケージをコンテキストごとに分割せずフラットな `persistence`
パッケージ 1 つにまとめているのは学習用サンプルとしての単純さを優先した
割り切りであり、実際には `cart` / `catalog` / `order` にとどまらず
8 コンテキスト分・23 ファイル・約 1600 行が同じ `persistence` パッケージに
混在している。

**フラット化が消している防壁。** ドメイン層では「order パッケージは
cart パッケージを import できない」がコンパイルエラーとして機械的に
強制される（`internal/domain/README.md` 「依存関係のルール」参照）。
しかし `persistence` パッケージ内では話が違う。同一パッケージ内の
`OrderModel` や `OrderRepository` の実装コードから `CartModel` や
`cart_repository.go` が定義する型・関数を、import 文を追加することなく
そのまま参照できてしまう。つまりドメイン層が持つ「コンテキストをまたぐ
参照はコンパイルエラーになる」という防壁が、persistence 層には存在しない。

現時点でコンテキストをまたぐ参照は発生していないことを確認済みである
（各 `*_repository.go` / `*_query_service.go` は自コンテキストの
モデルしか参照していない）。ただしこれは「実装者の規律」によって
保たれているだけであり、構造的に禁止されているわけではない点が
コストである。将来この越境が実際に発生した場合の対処としては、
`internal/infrastructure/persistence/cart` のようなコンテキストごとの
サブパッケージへ分割し、ドメイン層と同様にコンパイルエラーとして
越境を検出できるようにするのが定石である（実施するかどうかは規模に
応じた将来の判断とし、現時点では記述に留める）。

## TxManager（トランザクション境界の実装）

`internal/infrastructure/persistence.TxManager` は
`internal/application/tx.Manager` の GORM 実装である。

`Do` メソッドは GORM の `db.Transaction` を使ってトランザクションを開始し、
トランザクション用の `*gorm.DB` を `context.Context` に埋め込んでから
コールバック関数を実行する。

```go
func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
    return m.db.Transaction(func(txDB *gorm.DB) error {
        txCtx := WithDB(ctx, txDB)
        return fn(txCtx)
    })
}
```

各コンテキストのリポジトリ実装は、DB 操作の際に必ず
`persistence.DBFromContext(ctx, r.db)` を呼び出して使用する `*gorm.DB` を
取得する。トランザクション中であれば context に埋め込まれたトランザクション用の
DB ハンドルが、そうでなければリポジトリが保持している通常のコネクションプール用
DB ハンドルが返る。これにより、リポジトリ自身は「今トランザクション中かどうか」を
意識する必要がなく、常に同じ書き方で実装できる。

## コンテキストをまたぐ連携について

本サンプルでは、注文確定後にカートを空にするといったコンテキストを
またぐ連携を、以前はドメインイベント + インメモリのイベントバスで
実現していたが、仕組みの理解コストを下げるため application 層からの
直接呼び出しに変更した（詳細は `internal/application/README.md`
「コンテキストをまたぐ直接呼び出し」を参照）。そのため infrastructure
層にイベントバスの実装は存在しない。

反応するコンテキストが増えるなどして疎結合化のメリットが上回るように
なった場合は、`internal/domain/shared` に `DomainEvent` /
`EventPublisher` のようなポートを再度切り、インメモリバスや
Outbox パターン（集約の保存と同一トランザクションで Outbox テーブルへ
書き込み、別プロセスが Kafka / SQS 等へ非同期に発行する方式）による
実装をこの層に置くのが定石である。
