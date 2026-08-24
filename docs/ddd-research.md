# DDD 設計原則リサーチ(出典検証済み)

このリポジトリを DDD の観点で評価・改善するために収集した設計原則の一覧です。収集した主張はすべて一次ソース(原典 PDF・原文 HTML)に当たって再検証しており、各項目に確度を付記しています。逐語引用は原文(英語)のまま載せ、日本語で補足します。

- **確度ラベル**: 「逐語確認」= 原文から該当文言を直接確認 / 「実質確認」= 文言は異なるが内容を原文で確認 / 「推定適用」= 原典は別の文脈での主張であり、本リポジトリへの適用は解釈を含む
- このリポジトリのどこが原則から外れているか・どう直すかは [specs/](specs/) の改善スペックに分離します(本ファイルは原則集)

## 集約の設計 — Vernon の4ルール

出典: Vaughn Vernon, *Effective Aggregate Design* Part I–II ([Part I](https://www.dddcommunity.org/wp-content/uploads/files/pdf_articles/Vernon_2011_1.pdf) / [Part II](https://www.dddcommunity.org/wp-content/uploads/files/pdf_articles/Vernon_2011_2.pdf))。すべて逐語確認。

**ルール1: 真の不変条件を整合性境界の中でモデリングする。** トランザクションで守るべき不変条件を持つ範囲だけが集約になります。

> "And a properly designed bounded context modifies only one aggregate instance per transaction in all cases." (Part I, p.3)

いわゆる「1トランザクション=1集約」の原文です。ただし Vernon 自身が同じページで、これは絶対則ではなく目標だと明言しています。

> "Limiting the modification of one aggregate instance per transaction may sound overly strict. However, it is a rule of thumb and should be the goal in most cases." (Part I, p.3)

**ルール2: 集約は小さく設計する。** 大きなクラスタ集約は性能・スケールともに破綻します。実務статистикとして、

> "Niclas [Hedhman] reported that his team was able to design approximately 70% of all aggregates with just a root entity containing some value-typed properties. The remaining 30% had just two to three total entities." (Part I, p.5)

**ルール3: 他の集約は ID で参照する。**

> "Prefer references to external aggregates only by their globally unique identity, not by holding a direct object reference (or 'pointer')." (Part II, p.7)

**ルール4: 境界の外は結果整合性を使う。**

> "Thus, if executing a command on one aggregate instance requires that additional business rules execute on one or more other aggregates, use eventual consistency." (Part II, p.8)

Part II はさらに Evans の原典(DDD 本 p.128)を引いています: "Any rule that spans AGGREGATES will not be expected to be up-to-date at all times."

### ルールを破ってよい場合 — "Reasons To Break the Rules"

Part II には逸脱を明示的に扱う節があります(節名逐語: **"Reasons To Break the Rules"**、p.9–10)。

> "An experienced DDD practitioner may at times decide to persist changes to multiple aggregate instances in a single transaction, but only with good reason. What might some reasons be? I discuss four reasons here."

列挙されている4つの理由(節見出し逐語):

1. **Reason One: User Interface Convenience** — UI の都合(一括操作など)
2. **Reason Two: Lack of Technical Mechanisms** — メッセージング・タイマー等の技術的手段が使えない場合
3. **Reason Three: Global Transactions** — レガシー統合等でグローバルトランザクションが要求される場合
4. **Reason Four: Query Performance** — クエリ性能上の理由

なお「user-aggregate affinity(同時に同じ集約群を触るユーザーが1人に限られるか)」は、独立した第5の理由ではなく Reason Two の中の補助的な判断材料です(逐語確認)。

> "Consider an additional factor that could further support diverting from the rule: user-aggregate affinity."

**要点**: 複数集約を1トランザクションで更新することは「禁止」ではなく「正当な理由を明示できる場合に限り経験者の判断で許される逸脱」です。

## 境界づけられたコンテキストと統合パターン

- Fowler の [BoundedContext](https://martinfowler.com/bliki/BoundedContext.html) は、境界はデプロイ単位ではなく**モデルの適用範囲**で引かれるとし、1アプリ内の表現の分離(インメモリモデルと RDB モデルの分離など)も複数コンテキストの一形態に数えています。物理 DB の分離自体は要件ではありません(実質確認)
- Microsoft の [domain-analysis](https://learn.microsoft.com/en-us/azure/architecture/microservices/model/domain-analysis) は、コンテキスト間の関係パターンとして Customer-Supplier / Open Host Service + Published Language / Anti-Corruption Layer / Separate Ways を列挙し、コンテキストマップで統合点と責務を文書化することを求めています(実質確認)
- **Shared Kernel は小さく保ち、変更は両者協議** — Evans の DDD Reference から逐語確認([ddd-crew/context-mapping](https://github.com/ddd-crew/context-mapping) 経由、原典 [domainlanguage.com/ddd/reference](https://www.domainlanguage.com/ddd/reference/)):

> "Designate with an explicit boundary some subset of the domain model that the teams agree to share. **Keep this kernel small.** ... This explicitly shared stuff has special status, and shouldn't be changed without consultation with the other team."

## ドメインイベントとコンテキスト間連携

- Fowler の [DomainEvent](https://martinfowler.com/eaaDev/DomainEvent.html) は、イベントを監査証跡・イベントソーシングなどに有効としつつ、独特のアーキテクチャスタイルを持ち込むため**労力対効果を検討すべきトレードオフ**として提示しています(既定ではない)(実質確認)
- モジュラーモノリスの代表的参照実装 [kgrzybek/modular-monolith-with-ddd](https://github.com/kgrzybek/modular-monolith-with-ddd) は、イベント優先どころか**イベント限定**です(README 逐語確認):

> "**Modules communicate each other only asynchronously using Events Bus** - direct method calls are not allowed"

- 一方 Fowler の [MonolithFirst](https://martinfowler.com/bliki/MonolithFirst.html) は、良い境界を最初から引くのは難しく、まず注意深くモジュール化したモノリスで境界を発見せよという立場です。ただしこの記事はモノリス対マイクロサービスの文脈であり、「モノリス内の同期的なコンテキスト横断トランザクション」への直接の判定ではありません(**推定適用**)
- 「モノリス内の同期・単一トランザクションのコンテキスト横断オーケストレーション」に対する明示的な可否を述べた一次ソースは見つかっていません。純粋主義(Vernon ルール4・kgrzybek)は否、実利主義は許容という**係争中の論点**です

## CQRS

出典: Fowler, [CQRS](https://martinfowler.com/bliki/CQRS.html)(逐語確認)。

> "...beware that for most systems CQRS adds risky complexity."

> "Despite these benefits, you should be very cautious about using CQRS."

CQRS はシステム全体ではなく特定の bounded context に限って適用すべきであり、複雑なクエリが動機なら ReportingDatabase のような単純な代替も検討せよ、としています。読み取りモデルがドメインモデルを迂回すること自体は CQRS の定義そのものです。コンテキスト横断の JOIN を読み取り側で行うことの明示的な可否を述べた一次ソースはありません(迂回の許容からの**推定適用**)。

## アプリケーション層とリポジトリ

- **トランザクション制御はリポジトリではなくクライアント(アプリケーション層)が持つ** — Evans 原典(Ch.6)逐語確認:

> "Leave transaction control to the client. While the REPOSITORY will insert and delete from the database, it will ordinarily not commit anything."

- **リポジトリは集約ルートにのみ提供する** — Evans 原典逐語確認:

> "Only provide repositories for AGGREGATE roots that actually need direct access."

注意: よく引用される「1集約=1リポジトリ」は原文の表現ではありません。Evans はむしろ1つのリポジトリが型階層(例: TradeOrder ← BuyOrder/SellOrder)をまとめて扱うことを明示的に許しています。

- **不変条件はエンティティ(集約)のメソッドが守る** — Microsoft(逐語確認):

> "The entity's methods take care of the invariants and rules of the entity instead of having those rules spread across the application layer."

- Go の実例として Three Dots Labs は、アプリケーションサービス層を "We also have no logic here: just some orchestration." と説明しています([DDD Lite in Go](https://threedots.tech/post/ddd-lite-in-go-introduction/)、逐語確認)

## 貧血ドメインモデル(Anemic Domain Model)

- Fowler([AnemicDomainModel](https://martinfowler.com/bliki/AnemicDomainModel.html)、逐語確認): 見た目はドメインモデルでも振る舞いがなく "little more than bags of getters and setters" になっているものはアンチパターンで、ロジックをサービスに出し切ると "you essentially end up with Transaction Scripts" になる
- ただし Microsoft は明示的な例外を認めています(逐語確認):

> "if your microservice or Bounded Context is very simple (a CRUD service), the anemic domain model in the form of entity objects with just data properties might be good enough"

つまり「単純な CRUD コンテキストの貧血モデルは許容、ビジネスルールが濃いコンテキストでの貧血がアンチパターン」という使い分けです。

## Go でのパッケージ構成

- 最も引用される Go DDD 参照実装 [Wild Workouts](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example) は**コンテキストファースト**構成です(GitHub API でツリーを実測確認): `internal/trainer/` と `internal/trainings/` は各自 `app/`(command/query)・`domain/`・`adapters/`・`ports/` を持ちます。ただし `internal/users/` はフラットなパッケージで、全コンテキストが同じ内部構成という一律ルールではありません
- 「層ファースト vs コンテキストファースト」のどちらかを正典とする一次ソースは存在しません。Go コミュニティの慣行はコンテキストファースト寄り、というのが実態です(コミュニティ合意レベル)

## 戦略的設計の最低限

- [ddd-crew/ddd-starter-modelling-process](https://github.com/ddd-crew/ddd-starter-modelling-process) は Context Map(Connect ステップ)と Bounded Context Canvas(Define ステップ)を推奨ツールとして提示しています(実質確認)。なお「どんな規模でも有用」という文言は README には存在しません(検証で棄却)
- Microsoft の domain-analysis はユビキタス言語(コンテキストごとの共有語彙)を「中心」とし、コンテキストマップによる関係・統合点の文書化を基本成果物としています(実質確認)

## 検証で棄却・訂正された主張

crosscheck(原典への反証優先の照合)で、当初のリサーチから以下を訂正しました。**これらを引用しないでください。**

| 当初の主張 | 検証結果 |
|---|---|
| 「モジュール間の同期呼び出しが3回以上連鎖したらイベントに再設計せよ」(kgrzybek README とされていた) | **出典に存在しない**。README は同期呼び出し自体を禁止しており、回数のしきい値はどこにもない |
| Evans「1集約ルートにつき1リポジトリ」 | **不正確**。原文は「集約ルートにのみ提供せよ」であり、1:1 は要求していない(型階層を1リポジトリで扱う例を明示) |
| ddd-crew「コンテキストマップはどんな規模のプロジェクトでも有用」 | **文言が存在しない**。適用場面はシナリオ(グリーンフィールド/移行/大型プログラム等)で説明されている |
| user-aggregate affinity を「第5の理由」とする整理 | **不正確**。Reason Two 内の補助的判断材料 |

## 出典一覧

- Vaughn Vernon, Effective Aggregate Design [Part I](https://www.dddcommunity.org/wp-content/uploads/files/pdf_articles/Vernon_2011_1.pdf) / [Part II](https://www.dddcommunity.org/wp-content/uploads/files/pdf_articles/Vernon_2011_2.pdf) / [Part III](https://www.dddcommunity.org/wp-content/uploads/files/pdf_articles/Vernon_2011_3.pdf)
- Eric Evans, Domain-Driven Design(2003 最終原稿 PDF、Ch.6 Repository / Shared Kernel)/ [DDD Reference](https://www.domainlanguage.com/ddd/reference/)
- Martin Fowler bliki: [BoundedContext](https://martinfowler.com/bliki/BoundedContext.html) / [CQRS](https://martinfowler.com/bliki/CQRS.html) / [AnemicDomainModel](https://martinfowler.com/bliki/AnemicDomainModel.html) / [DomainEvent](https://martinfowler.com/eaaDev/DomainEvent.html) / [MonolithFirst](https://martinfowler.com/bliki/MonolithFirst.html)
- Microsoft Learn: [domain-analysis](https://learn.microsoft.com/en-us/azure/architecture/microservices/model/domain-analysis) / [microservice-domain-model](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/microservice-domain-model)
- [ddd-crew/ddd-starter-modelling-process](https://github.com/ddd-crew/ddd-starter-modelling-process) / [ddd-crew/context-mapping](https://github.com/ddd-crew/context-mapping)
- [ThreeDotsLabs/wild-workouts-go-ddd-example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example) / [DDD Lite in Go](https://threedots.tech/post/ddd-lite-in-go-introduction/)
- [kgrzybek/modular-monolith-with-ddd](https://github.com/kgrzybek/modular-monolith-with-ddd)
