package shared

import (
	"context"
	"time"
)

// DomainEvent は「ドメイン内で起きた重要な出来事」を表すインターフェースである。
//
// DDD ではエンティティの状態変化そのものではなく、「何が起きたか」という
// 事実（過去形の名前を持つ）をイベントとしてモデリングする。
// 例: OrderPlaced（注文が確定した）、CartItemAdded（カートに商品が追加された）。
//
// イベントは発生時刻を持つ不変の事実であり、集約（Aggregate）はこれを
// Record することで「自分の中で何が起きたか」を記録し、アプリケーション層が
// 永続化成功後にこれを取り出して配信する。
type DomainEvent interface {
	// EventName はイベントの種別を一意に表す文字列を返す。
	// 例: "cart.item_added"、"order.placed" のようなドット区切りが読みやすい。
	EventName() string
	// OccurredAt はイベントが発生した時刻を返す。
	OccurredAt() time.Time
}

// EventPublisher はドメインイベントを配信するためのポート（インターフェース）である。
//
// ドメイン層・アプリケーション層はこのインターフェースにのみ依存し、
// 実際の配信方法（インメモリの同期配信か、メッセージキューへの非同期配信か等）は
// インフラストラクチャ層が実装する。これはヘキサゴナルアーキテクチャにおける
// 「ポートとアダプタ」パターンそのものであり、ドメインを配信の実装詳細から
// 切り離すことを目的としている。
type EventPublisher interface {
	Publish(ctx context.Context, events ...DomainEvent) error
}
