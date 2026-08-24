// Package eventhandler は cart コンテキストが「他コンテキストで起きた
// 出来事に反応する」ためのイベントハンドラを集めたパッケージである。
//
// なぜ command パッケージではなくこの専用パッケージに置くのか？
// command パッケージのユースケースは「HTTP など外部からの直接の指示」に
// よって起動されるのに対し、ここに置くハンドラは「他コンテキストの
// ドメインイベント」によって起動される。トリガーの種類が異なるため
// パッケージを分け、「このカート操作は誰の意思で実行されるのか」を
// ディレクトリ構成だけで把握できるようにしている。
package eventhandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ClearCartOnOrderPlaced は order コンテキストの OrderPlaced イベントに
// 反応してカートを空にするイベントハンドラである。
//
// 【これはコンテキスト間の疎結合な連携を示す教材的な例である】
// order パッケージは cart パッケージを一切知らない（import しない）。
// それでも「注文が確定したらカートを空にする」という cart 側の業務ルールを
// 実現できるのは、cart 側がイベントを購読して自律的に反応しているからである。
// order 側は「注文が確定した」という事実を発信するだけで、その事実に
// 誰が・どう反応するかには関与しない。これにより、将来「注文確定時に
// ポイントを付与する」「注文確定時に在庫を引き当てる」といった新しい
// 反応が増えても、order 側のコードは一切変更する必要がない
// （Open-Closed Principle をコンテキスト間の連携に適用した形）。
type ClearCartOnOrderPlaced struct {
	cartRepo cart.Repository
}

// NewClearCartOnOrderPlaced は ClearCartOnOrderPlaced を生成する。
//
// tx.Manager を受け取らない点に注意する。このハンドラは
// PlaceOrderUseCase が txManager.Do の中（= 注文の Save と同じ
// トランザクション）で publisher.Publish を呼んだ結果として起動される。
// 同期インメモリバス（event.Bus）は Publish をその場で呼び出すだけなので、
// Handle に渡ってくる ctx にはすでに注文確定と同じトランザクション用の
// *gorm.DB が埋め込まれている。ここで新たにトランザクションを開始して
// しまうと、注文の保存とカートのクリアが別トランザクションになり、
// 「注文は確定したのにカートは空にならない」という不整合が起こり得る
// 窓が生まれてしまう。そのため、このハンドラは自前でトランザクションを
// 開始せず、渡された ctx をそのまま cart.Repository に渡すことで、
// 呼び出し元が既に開いているトランザクションに乗り続ける設計にしている。
func NewClearCartOnOrderPlaced(cartRepo cart.Repository) *ClearCartOnOrderPlaced {
	return &ClearCartOnOrderPlaced{cartRepo: cartRepo}
}

// Handle は shared.DomainEvent を受け取り処理する。
//
// この関数のシグネチャは internal/infrastructure/event.Handler
// （func(ctx, shared.DomainEvent) error）にそのまま代入できる形にして
// ある。cart（アプリケーション層）が event パッケージ（インフラ層）を
// import してしまうと依存の向きが逆転するため、ここでは素の
// メソッドとして定義するに留め、bus.Subscribe への配線は main（DI の
// 組み立て役）が h.Handle という形で行う想定である。
func (h *ClearCartOnOrderPlaced) Handle(ctx context.Context, e shared.DomainEvent) error {
	placed, ok := e.(order.OrderPlaced)
	if !ok {
		// Bus は EventName で購読先を振り分けるため、"order.placed" として
		// 登録された Handle には本来 order.OrderPlaced 以外は渡ってこない。
		// それでも型アサーションに失敗した場合は、配線ミス（別イベントを
		// 誤って同じ名前で購読してしまった等）という「本来あってはならない
		// プログラミングエラー」であるため、黙って無視するのではなく
		// エラーとして表面化させる。
		return fmt.Errorf("cart: eventhandler: unexpected event type %T for %q", e, e.EventName())
	}

	customerID, err := cart.NewCustomerID(placed.CustomerID().String())
	if err != nil {
		return err
	}

	c, err := h.cartRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			// すでにカートが存在しない（＝すでに空、あるいは元から
			// 何も入れていなかった）状態は、このハンドラにとって
			// エラーではない。「空にする」という操作の結果として見れば、
			// カートが存在しないことと明細が 0 件であることは同じ状態を
			// 指すため、ここでは何もせず正常終了する（冪等性の確保）。
			return nil
		}
		return err
	}

	c.Clear()
	return h.cartRepo.Save(ctx, c)
}
