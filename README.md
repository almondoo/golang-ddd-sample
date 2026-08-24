# golang-ddd-sample

Go で DDD(ドメイン駆動設計)の戦術パターンを学ぶためのサンプルリポジトリです。EC(ネットショップ)を題材に、商品カタログ・カート・注文の 3 コンテキストをオニオンアーキテクチャ + 軽量 CQRS で実装しています。読んで学ぶことを主目的にしており、すべてのコードに日本語の解説コメントを付けています。

## 技術スタック

- Go(標準 `net/http`、Go 1.22+ のメソッド付きルーティング)
- GORM + PostgreSQL(永続化)
- almondoo/wire(コンパイル時 DI コード生成)
- 外部フレームワークはドメイン層に持ち込まない方針

## 動かし方

コードリーディングだけであればビルドとテストのみで完結します。DB は不要です。

```sh
go build ./...
go test ./...
```

API サーバーとして起動する場合は PostgreSQL を用意し、DSN を環境変数で渡します。

```sh
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=ddd_sample port=5432 sslmode=disable"
go run ./cmd/api
```

依存の組み立て(`cmd/api/wire.go` のプロバイダ宣言)を変更した場合は、wire で組み立てコードを再生成します。

```sh
go run github.com/almondoo/wire/cmd/wire ./cmd/api
```

## ディレクトリ構成

```
.
├── cmd/api/                     # エントリポイント
│   ├── main.go                  # 設定読み込みと起動のみ
│   ├── wire.go                  # DI の設計図(プロバイダの宣言)
│   └── wire_gen.go              # wire が生成した組み立てコード
├── docs/
│   └── execution-flow.md        # リクエスト実行順序の図解
├── internal/
│   ├── domain/                  # ドメイン層(最内層・外部依存ゼロ)
│   │   ├── shared/              # 共有カーネル: Money、ドメインイベント、集約ベース
│   │   ├── catalog/             # 商品カタログ: Product 集約、値オブジェクト、リポジトリIF
│   │   ├── cart/                # カート: Cart 集約、CartItem、数量などの不変条件
│   │   └── order/               # 注文: Order 集約、状態遷移、OrderPlaced イベント
│   ├── application/             # ユースケース層(1ファイル = 1ユースケース)
│   │   ├── tx/                  # トランザクション境界のポート(インターフェース)
│   │   ├── catalog/
│   │   │   ├── command/         # 書き込み系: 商品登録、価格改定
│   │   │   └── query/           # 読み取り系: 一覧、詳細
│   │   ├── cart/
│   │   │   ├── command/         # カートへの追加・削除
│   │   │   ├── query/           # カート内容の取得
│   │   │   └── eventhandler/    # OrderPlaced 購読 → カートを空にする
│   │   └── order/
│   │       ├── command/         # 注文確定、支払い、発送、キャンセル
│   │       └── query/           # 注文詳細の取得
│   ├── infrastructure/          # 技術詳細(最外層)
│   │   ├── persistence/         # GORM モデル、リポジトリ実装、クエリサービス実装
│   │   └── event/               # インメモリの同期イベントバス
│   └── presentation/
│       └── controller/          # HTTP コントローラ、ルーティング、エラー→ステータス変換
└── README.md
```

各層の役割と設計判断は、それぞれのディレクトリ直下の README で解説しています。

- [internal/domain/README.md](internal/domain/README.md) — 依存方向のルール、集約・値オブジェクト・リポジトリIF
- [internal/application/README.md](internal/application/README.md) — command / query 分離、トランザクション境界、イベント発行
- [internal/infrastructure/README.md](internal/infrastructure/README.md) — ドメインモデルと GORM モデルの分離(データマッパー)
- [internal/presentation/README.md](internal/presentation/README.md) — DTO 変換、エラーマッピング

## API エンドポイント

| メソッドとパス | ユースケース | 種別 |
|---|---|---|
| `POST /products` | 商品を登録する | command |
| `GET /products` | 商品一覧を取得する | query |
| `GET /products/{id}` | 商品詳細を取得する | query |
| `PUT /products/{id}/price` | 価格を改定する | command |
| `GET /carts/{customerID}` | カート内容を取得する | query |
| `POST /carts/{customerID}/items` | カートに商品を追加する | command |
| `DELETE /carts/{customerID}/items/{productID}` | カートから商品を除く | command |
| `POST /orders` | カートの内容から注文を確定する | command |
| `GET /orders/{id}` | 注文詳細を取得する | query |
| `POST /orders/{id}/pay` | 支払いを記録する(状態遷移) | command |
| `POST /orders/{id}/ship` | 発送を記録する(状態遷移) | command |
| `POST /orders/{id}/cancel` | 注文をキャンセルする | command |

## このサンプルの学習ポイント

**ドメイン層の純粋性。** ドメイン層は GORM も HTTP も知りません。永続化のための struct(GORM タグ付き)は infrastructure 側に置き、ドメインモデルと相互変換します。ORM の都合がビジネスルールの表現を歪めないための分離です。

**軽量 CQRS の非対称性。** command はドメインモデルを経由して不変条件を守りながら書き込み、query はドメインを迂回して GORM から直接 DTO を組み立てます。「読み取りに集約の保護は不要」という判断を、あえて非対称なコードで示しています。

**値オブジェクトによる primitive obsession の回避。** `Money`(負値と通貨不一致を拒否)、各種 ID 型(UUID のラッパー)で、`int64` や `string` の生値がドメインを流れることを防ぎます。

**ドメインイベントによるコンテキスト間連携。** 注文確定時に `Order` 集約が `OrderPlaced` を記録し、ユースケースが永続化後に発行、カート側のハンドラがそれを購読してカートを空にします。注文コンテキストはカートの都合を知りません。

**トランザクション境界はユースケースが決める。** リポジトリは自分でトランザクションを張らず、`tx.Manager` が context 経由で伝搬するトランザクションに参加します。

**依存の組み立ては composition root に集約。** controller は command / query を問わず必ず usecase を経由し、その依存関係は `cmd/api/wire.go` に宣言して wire がコンパイル時にコード生成します。配線ミスが実行時エラーではなくコンパイルエラーとして検出できます。

リクエストが各層をどの順序で通るかは [docs/execution-flow.md](docs/execution-flow.md) の図を参照してください。
