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
割り切りであり、`cart` / `catalog` / `order` の 3 コンテキストが混在する。
コンテキストが増えたりリポジトリ実装が肥大化したりした場合は、
`internal/infrastructure/cart` のようなサブパッケージへの分割を
将来的に検討してよい。

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

## イベントバス（Bus）

`internal/infrastructure/event.Bus` は
`internal/domain/shared.EventPublisher` のインメモリ実装である。
`Subscribe` でハンドラーを登録し、`Publish` が呼ばれると同期的に該当する
ハンドラーを実行する。

学習用サンプルとしてシンプルさを優先しているが、実運用では
配信の信頼性（配信漏れ・二重配信の防止）や他プロセスへの配信を考慮し、
下記のような **Outbox パターン** への発展を検討するとよい。

1. ユースケースはドメインイベントを、集約の保存と同一トランザクション内で
   Outbox テーブルに書き込む（コミット可否が一致するため配信漏れが起きない）。
2. 別プロセス（ポーラーやCDC）が Outbox テーブルを読み取り、
   Kafka や SQS 等のメッセージブローカーへ非同期に発行する。
3. 発行済みのレコードには印を付け、二重発行を避ける。

`shared.EventPublisher` というポートを切っていることで、
将来この `Bus` を Outbox ベースの実装に差し替えても、
アプリケーション層・ドメイン層のコードは一切変更する必要がない。
