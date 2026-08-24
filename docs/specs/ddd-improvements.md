# DDD 改善スペック — 直すべき点と修正方針

このリポジトリを DDD の確立された原則([ddd-research.md](../ddd-research.md) に出典検証済みで整理)と突き合わせて監査し、直すべき点を優先度順にまとめたものです。各項目は「現状(根拠の file:line)→ 反している/不足している原則 → 修正方針 → 工数」の形で書いています。

**監査の方法**: ①一次ソースから DDD 原則をリサーチ → ②別エージェントが原典 PDF・原文 HTML に当たって逐語 crosscheck(誤引用3件を棄却済み)→ ③検証済み原則を基準に、コード全体を file:line 単位で静的監査。実行時挙動(実 DB でのロック競合・ロールバック)は未検証です。

**前提**: 「ドメインイベント不採用・直接呼び出し」「単一 DB スキーマ」「層ファーストのパッケージ構成」は意図して選んだトレードオフであり、本スペックでは覆しません。ただし**そのコストをドキュメントが過少に説明している箇所**は修正対象です。

## 対応状況(2026-08-24 時点)

| 項目 | 状態 |
|---|---|
| 1(shipping/README の事実誤り) | 対応済み |
| 2(NewOrder の time.Now())| 対応済み(`now` を引数化、PlacedAt のテスト追加) |
| 3(エラーメッセージ言語混在) | 対応済み(英語に統一) |
| 4(トランザクションコストの過少申告) | 対応済み(application/README・order/README に実数と Vernon の位置づけを追記) |
| 5(usecase テストゼロ) | **未対応** |
| 6(コンテキストマップ) | 対応済み([../context-map.md](../context-map.md)) |
| 7(フラット persistence の防壁) | 対応済み(ドキュメント明記。パッケージ分割は未実施・将来の選択肢のまま) |
| 8(JOIN の過大主張) | 対応済み |
| 9(ID 型対照表) | 対応済み(context-map 内) |
| 10(ドメインサービスの記述) | 対応済み |
| 11(NewID の非決定性注記) | 対応済み |
| 12(shared の位置づけ明文化) | 対応済み(context-map 内) |

## 優先度: 高

### 1. shipping/README.md が実装済みの ShipOrderUseCase を「将来の追加」と記述している

- **現状**: `internal/domain/shipping/README.md:61-66` が order との連携を「将来アプリケーション層に追加される発送ユースケース(例: ShipOrderUseCase に相当する調整役)」と仮定形で説明。実際には `internal/application/usecase/order/ship_order.go` が実装済みで、wire 配線・`POST /orders/{id}/ship` まで到達可能
- **問題**: 教材として事実に反する記述。shipping の README だけを読んだ学習者は「この連携は未実装」と誤解する
- **修正方針**: 「注文コンテキストとの関係」節を実在の `ShipOrderUseCase` への参照に書き換える
- **工数**: S

### 2. `NewOrder` がドメイン層内で `time.Now()` を直接呼んでいる

- **現状**: `internal/domain/order/order.go:53` が `placedAt: time.Now()`。一方 `internal/domain/coupon/coupon.go:154-159` と coupon/README は「ドメイン層で time.Now() を呼ばず now を引数で受ける(時刻もまた入力)」という原則を明文化し、同じ `PlaceOrderUseCase` 内で `cp.Use(time.Now())` として正しく注入している(`place_order.go:216`)。矛盾の実害はテストに現れており、`order_test.go` は `PlacedAt()` を一度も検証せず、時刻が必要なテストはすべて `ReconstructOrder` で回避している
- **原則**: ドメインの純粋性・決定性(リポジトリ自身が coupon で採用済みの原則。外部では Three Dots Labs 等の Go DDD 実践も同様)
- **修正方針**: `NewOrder(customerID, items, now)` に `now time.Time` を追加し、`PlaceOrderUseCase` から注入。`order_test.go` に `PlacedAt` の検証を追加
- **工数**: S

### 3. ドメインルール違反エラーメッセージの言語が混在している

- **現状**: `internal/domain/**` の `NewDomainRuleError` は全コンテキストで英語(例: `cart/cart.go:117`)。ところが order の usecase 3ファイルだけ日本語(`place_order.go:115,128,152,167,207`「顧客が存在しません」等、`ship_order.go:133`、`cancel_order.go:75`)。同じ usecase 層でも `coupon/issue_coupon.go:74,84` は英語。同じ 422 エラー面に2言語が混在し、学習者がどちらの規約に従うべきか判断できない
- **原則**: ユビキタス言語の一貫性(用語・表現の一貫はコード全体で維持する)
- **修正方針**: 英語に統一(コードベースの多数派に合わせる。「コメントは日本語・識別子とメッセージは英語」という既存方針とも整合)。日本語化したい場合は presentation 層での変換として設計し直す
- **工数**: S

### 4. 注文確定の「1トランザクションのコスト」をドキュメントが過少申告している

- **現状**: `PlaceOrderUseCase`(`place_order.go:108-259`)は1トランザクションで最大 **4集約種・約23インスタンス**(Cart + Order + Stock×最大20 + Coupon)を更新し、他 **5コンテキスト**のドメインパッケージを import している(`place_order.go:9-15`)。しかし `internal/application/README.md` と `internal/domain/order/README.md` はコストを「order→cart の結合」という単一の関係としてしか説明していない。`ShipOrderUseCase` も同様に3集約種を1トランザクションで更新
- **原則**: Vernon のルール1「1トランザクションで変更する集約インスタンスは1つ」— ただし同 Part II の "Reasons To Break the Rules" が「正当な理由(UI の都合/技術的手段の欠如/グローバルトランザクション/クエリ性能)があれば経験者の判断で逸脱可、ただし理由を明示せよ」としている([ddd-research.md](../ddd-research.md) 参照、逐語確認済み)
- **修正方針**: コードは変えない(学習コスト優先の意図的な選択)。application/README.md と order/README.md に、(a) 実際の集約数・依存コンテキスト数、(b) Vernon ルールに対する明示的な逸脱でありどの "Reason" に相当するか(本サンプルは Reason Two: イベント基盤という技術的手段を意図的に持たない)、(c) 本番規模での帰結(ロック競合・ロールバック範囲・段階的分離の難しさ)を追記
- **工数**: S(ドキュメントのみ)

### 5. usecase 層のテストがゼロ

- **現状**: `*_test.go` は `internal/domain/*` の8ファイルのみ。リポジトリ内で最も複雑な `PlaceOrderUseCase`(264行・5リポジトリ・在庫引当ループ・クーポン分岐)に実行可能な検証が一切ない。モック/フェイクのヘルパも皆無
- **原則**: ポート(1〜2メソッドの interface)を切っている設計の見返りはテスト容易性であり、それを示さないと設計の説得力が半減する
- **修正方針**: モックライブラリは使わず手書きフェイク(map ベースの fakeRepository、`fn(ctx)` を素通しする fakeTxManager)で、まず `PlaceOrderUseCase`(正常系/在庫不足/クーポン無効/顧客不在)、次に find-or-create 分岐を持つ `AddItemUseCase`・`SetStockUseCase` をテストする
- **工数**: M

### 6. コンテキストマップがない

- **現状**: `docs/` には実行順序図のみ。8コンテキスト間の関係(order の app 層が5コンテキストに依存、cart→products の読み取り JOIN、shared kernel)は各所の README に散在し、全体像を1枚で見られる文書がない
- **原則**: コンテキストマップとユビキタス言語の文書化は戦略的設計の基本成果物(Microsoft domain-analysis、ddd-crew。[ddd-research.md](../ddd-research.md) 参照)
- **修正方針**: `docs/context-map.md` を新規作成。コンテキストを箱、連携を「関係パターン名付きの有向エッジ」で mermaid 図化(order→他4コンテキスト: アプリケーション層オーケストレーション / cart→catalog: 読み取り側のスキーマ参照 / 全体→shared: Shared Kernel)。既存 README の記述を集約するだけで書ける
- **工数**: S〜M

## 優先度: 中

### 7. フラットな persistence パッケージがコンパイル時のコンテキスト境界を消している

- **現状**: `internal/infrastructure/persistence` は8コンテキスト分・23ファイル・約1600行が1パッケージ。ドメイン層では「order は cart を import できない」がコンパイルエラーとして強制されるのに、persistence 内では全コンテキストのモデルが相互に見え、将来の越境を止める仕組みがない(現時点で越境は未発生と確認済み)。`infrastructure/README.md` はこの構造的コストに触れていない
- **修正方針**: まず README にこの「失われている防壁」を明記(S)。実際の分割(`persistence/catalog/` 等のサブパッケージ化)は将来の選択肢として記述に留める(実施するなら M)
- **工数**: S(ドキュメント)

### 8. cart の読み取り JOIN が「疎結合を破らない」と過大に主張している

- **現状**: `cart_query_service.go:55` は catalog の物理カラム名(`products.name`, `products.price_amount`)をハードコードしており、catalog 側のテーブル変更でコンパイルエラーなしに実行時破壊される。`cart/dto.go:25-26` と `cart/README.md:28-30` は「コンテキスト間の疎結合を破らない」と説明
- **修正方針**: 「書き込み側は型安全に分離されているが、読み取り側にはスキーマ結合が残る」と正確に書き換える
- **工数**: S

### 9. ID 型の対照表(ミニ用語集)がない

- **現状**: `ProductID`×3、`CustomerID`×3、`OrderID`×2 が意図的に重複定義されているが、同一 UUID 空間を共有することを知るにはファイルを見比べるしかない。また空文字チェックが11ファイルにコピーされている保守コストはどこにも書かれていない
- **修正方針**: コンテキストマップ(項目6)に ID 対照表を含める。domain/README.md に重複の保守コストを1文追記
- **工数**: S

### 10. domain/README.md が存在しないドメインサービスを「この層に置くもの」として挙げている

- **現状**: `internal/domain/README.md:64-65` がドメインサービスを列挙するが、リポジトリ内に実例はゼロ。監査の結論として、現状のコードに「ドメインサービスに移すべき置き違いロジック」は存在しない(コンテキスト横断調整はすべて application 層の責務で正しい)ため、これは記述の過大提示
- **修正方針**: 「本サンプルでは各コンテキストの集約が1つずつであり、ドメインサービスが必要になる場面がないため未使用」と明記する
- **工数**: S

## 優先度: 低

### 11. `shared.NewID()` の非決定性が原則説明と未照合

- **現状**: 全コンテキストの生成コンストラクタがドメイン層内で UUID を生成する。「時刻も入力」という coupon の原則を ID に当てはめれば同種の非決定性だが、テスト実害はなく(ID の値に依存する検証はない)、一般的な DDD 実践でも広く許容されている
- **修正方針**: コードは変えず、domain/README.md に「ID 生成はドメイン内に置く判断をした(時刻と違いテストが値に依存しないため)」と1文追記
- **工数**: S

### 12. README の「7コンテキスト」の数え方が暗黙

- **現状**: shared(共有カーネル)を bounded context に数えない前提が明文化されていない
- **修正方針**: コンテキストマップに「shared は bounded context ではなく Shared Kernel」と明記
- **工数**: S

## 実施順の提案

1. 事実誤りの修正(項目1)と小さなコード修正(項目2・3)
2. ドキュメントの正確化(項目4・8・10・7)
3. 戦略的設計の補完(項目6・9・12)
4. テスト整備(項目5)

すべて実施しても既存の設計判断(直接呼び出し・単一スキーマ・層ファースト)は変わりません。「原則から外れている箇所を、外れていると正確に説明できる状態」にすることが本スペックの目的です。
