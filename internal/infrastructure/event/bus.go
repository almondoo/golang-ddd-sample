package event

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// Handler はドメインイベントを処理する関数の型である。
// 例えば「注文が確定したら在庫を引き当てる」「メールを送る」といった
// 副作用をここに実装する。
type Handler func(ctx context.Context, e shared.DomainEvent) error

// Bus はインメモリのイベントバス（shared.EventPublisher の実装）である。
//
// 【同期配信 vs 非同期配信のトレードオフ】
// このサンプルでは学習のしやすさを優先し、Publish が呼ばれたその場で
// 同期的にハンドラーを実行する、最もシンプルな実装を採用している。
// 同期配信のメリットは実装が単純で、ハンドラーのエラーを呼び出し元に
// そのまま伝播できること。デメリットは、ハンドラーが遅い・失敗する場合に
// ユースケース全体の応答が遅延・失敗してしまうこと、また単一プロセス内で
// しか配信できないこと（他サービスへは配信できない）である。
//
// 実際のプロダクションシステムでは、ここを Kafka / SQS / RabbitMQ 等の
// メッセージブローカーに差し替え、
//  1. ドメインイベントをまず Outbox テーブルに書き込む（同一トランザクション内）
//  2. 別プロセスが Outbox を読み取りブローカーに非同期発行する
//
// という Outbox パターンを使うことで、配信の信頼性とユースケースの
// 応答性を両立させることが多い。shared.EventPublisher というポートを
// 切っているのは、まさにこの実装差し替えを容易にするためである。
type Bus struct {
	handlers map[string][]Handler
}

// NewBus は Bus を生成する。
func NewBus() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

// コンパイル時に Bus が shared.EventPublisher を満たすことを保証するアサーション。
var _ shared.EventPublisher = (*Bus)(nil)

// Subscribe は eventName に対応するハンドラーを登録する。
// 同じ eventName に複数のハンドラーを登録した場合、登録順に呼び出される。
func (b *Bus) Subscribe(eventName string, h Handler) {
	b.handlers[eventName] = append(b.handlers[eventName], h)
}

// Publish はイベントを購読中のハンドラーへ同期的に配信する。
//
// 複数イベント・複数ハンドラーを渡された場合、先頭から順に処理し、
// いずれかのハンドラーがエラーを返した時点で処理を中断してそのエラーを
// 返す（それ以降のイベント・ハンドラーは実行されない）。
// これは「一部だけ成功した中途半端な状態」を避け、失敗を呼び出し元に
// 早期かつ明確に伝えるためのシンプルな方針である。
func (b *Bus) Publish(ctx context.Context, events ...shared.DomainEvent) error {
	for _, e := range events {
		for _, h := range b.handlers[e.EventName()] {
			if err := h(ctx, e); err != nil {
				return err
			}
		}
	}
	return nil
}
