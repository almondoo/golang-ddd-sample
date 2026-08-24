# DDD 要素解説

このディレクトリはDDDの構成要素を1要素=1ファイルで解説します。それぞれの実例は本リポジトリの実際のコードを参照しており、図で構造・フローをイメージできるようにしています。

| ファイル | 内容 | 図 |
|---|---|---|
| [aggregate.md](aggregate.md) | 集約/集約ルート/不変条件。Customer+Address、Order+OrderItem、Cart | 集約境界のflowchart(Customer/Order) |
| [entity.md](entity.md) | エンティティ。AddressとCartItem/OrderItemの対比 | 識別子維持 vs 置き換えの比較図 |
| [value-object.md](value-object.md) | 値オブジェクト。Money、CouponCode、各種ID型 | NewMoneyという唯一の関門のflowchart |
| [factory.md](factory.md) | 生成と再構築の分離。New*とReconstruct*のペア | 2経路(生成/再構築)のflowchart |
| [repository.md](repository.md) | リポジトリ。ドメイン層IF・infrastructure実装 | 依存性逆転のflowchart |
| [application-service.md](application-service.md) | アプリケーションサービス(usecase)。オーケストレーション | PlaceOrderUseCaseのflowchart |
| [domain-service.md](domain-service.md) | ドメインサービス。本リポジトリでは未使用とその理由 | ロジック置き場所判定のflowchart |
| [domain-event.md](domain-event.md) | ドメインイベント。本リポジトリでは不採用とその経緯 | 採用時 vs 現状の比較シーケンス図 |
| [bounded-context.md](bounded-context.md) | 境界づけられたコンテキスト。7コンテキスト+shared | コンテキスト概観flowchart |
| [shared-kernel.md](shared-kernel.md) | 共有カーネル。小さく保つ原則とinternal/domain/shared | 共有カーネル参照関係のflowchart |
| [ubiquitous-language.md](ubiquitous-language.md) | ユビキタス言語。ID型対照表を整備済み・エラーメッセージは英語に統一済み(修正の経緯つき) | ID型対照表 |
| [cqrs.md](cqrs.md) | CQRS(軽量版)。commandはドメイン経由、queryはDTO直行 | 2経路(command/query)のflowchart |
| [onion-architecture.md](onion-architecture.md) | オニオンアーキテクチャと依存性逆転。4層と依存方向 | 4層依存関係のflowchart |

より詳しい背景は以下を参照してください。

- [../ddd-research.md](../ddd-research.md) — 出典検証済みのDDD設計原則リサーチ(逐語引用付き)
- [../specs/ddd-improvements.md](../specs/ddd-improvements.md) — 本リポジトリが原則からどこで外れているか、その修正方針
