package command

import (
	"context"
	"errors"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	"github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	"github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// PlaceOrderInput は注文確定ユースケースの入力である。
type PlaceOrderInput struct {
	CustomerID string
}

// PlaceOrderOutput は注文確定ユースケースの出力である。
type PlaceOrderOutput struct {
	OrderID string
}

// PlaceOrderUseCase は「カートの中身をもとに注文を確定する」ユースケースである。
//
// 【なぜこのユースケースが 3 つのコンテキストの Repository に依存するのか】
// order ドメインパッケージは cart / catalog を import しないと決めた
// （境界づけられたコンテキストの自律性）。しかし現実の業務としては
// 「カートの中身」と「その時点の商品価格」の両方を参照しなければ注文は
// 組み立てられない。この「複数コンテキストにまたがる読み取りと調整」を
// 引き受けるのがアプリケーション層の役割である。ドメイン層はあくまで
// 「1 つのコンテキストの中で完結する不変条件」だけを守ればよく、
// コンテキストをまたぐオーケストレーションはその一段上の層（アプリケーション層）
// に置く、というのが本サンプル全体を通じたレイヤリングの方針である。
type PlaceOrderUseCase struct {
	orderRepo   order.Repository
	cartRepo    cart.Repository
	catalogRepo catalog.Repository
	txManager   tx.Manager
	publisher   shared.EventPublisher
}

// NewPlaceOrderUseCase は PlaceOrderUseCase を生成する。
func NewPlaceOrderUseCase(
	orderRepo order.Repository,
	cartRepo cart.Repository,
	catalogRepo catalog.Repository,
	txManager tx.Manager,
	publisher shared.EventPublisher,
) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		catalogRepo: catalogRepo,
		txManager:   txManager,
		publisher:   publisher,
	}
}

// Execute は注文確定ユースケースを実行する。
//
// 処理の流れ:
//  1. 対象顧客のカートを読み込む（存在しない・空である場合はドメインルール違反）
//  2. カート明細ごとに catalog から現在の商品情報（名前・価格）を取得し、
//     OrderItem のスナップショットを組み立てる
//  3. Order 集約を新規生成する（この時点で内部に OrderPlaced が記録される）
//  4. Order を永続化する
//  5. 永続化が成功した後にだけ、記録済みのイベントを配信する
//
// 【カートのクリアは誰の責任か】
// このユースケースはカートを明示的に空にする呼び出しを行わない。
// 代わりに、Save 成功後に配信される OrderPlaced イベントを
// cart コンテキスト側のイベントハンドラ（ClearCartOnOrderPlaced）が購読し、
// 自分自身の責任でカートを空にする。order 側から cart.Repository.Save を
// 直接呼んで空にすることも技術的には可能だが、そうすると「注文確定」という
// 1 つのユースケースが cart 集約の状態変更まで背負うことになり、
// 責任の境界が曖昧になる。イベント駆動にすることで「カートが注文確定に
// 反応してどう振る舞うか」という判断を cart コンテキスト自身に委ねられる。
func (uc *PlaceOrderUseCase) Execute(ctx context.Context, in PlaceOrderInput) (PlaceOrderOutput, error) {
	cartCustomerID, err := cart.NewCustomerID(in.CustomerID)
	if err != nil {
		return PlaceOrderOutput{}, err
	}
	orderCustomerID, err := order.NewCustomerID(in.CustomerID)
	if err != nil {
		return PlaceOrderOutput{}, err
	}

	var output PlaceOrderOutput
	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		c, err := uc.cartRepo.FindByCustomerID(ctx, cartCustomerID)
		if err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				// カートがまだ 1 度も作られていない状態は「注文できるものが
				// 何もない」という業務ルール違反として扱う。404（リソースが
				// 存在しない）ではなく 422（ルール違反）として表現したいため、
				// ここで shared.NewDomainRuleError に変換する。
				return shared.NewDomainRuleError("order: カートが空です。商品を追加してから注文してください")
			}
			return err
		}
		if c.IsEmpty() {
			return shared.NewDomainRuleError("order: カートが空です。商品を追加してから注文してください")
		}

		items := make([]order.OrderItem, 0, len(c.Items()))
		for _, cartItem := range c.Items() {
			productID, err := catalog.NewProductID(cartItem.ProductID().String())
			if err != nil {
				return err
			}

			product, err := uc.catalogRepo.FindByID(ctx, productID)
			if err != nil {
				if errors.Is(err, shared.ErrNotFound) {
					// 注文確定という「重い」操作のタイミングで初めて商品の
					// 実在を検証する（cart.AddItem の時点ではあえて検証しない
					// 設計だった。理由は cart/README.md を参照）。ここで見つから
					// ないのは「カートに入れた後に商品が削除された」等の業務上
					// 起こりうる状況であり、システム的な 404 ではなく注文という
					// 操作自体のドメインルール違反として扱う。
					return shared.NewDomainRuleError("order: 商品 %s が見つかりません", productID)
				}
				return err
			}

			orderItem, err := order.NewOrderItem(
				product.ID().String(),
				product.Name(),
				product.Price(),
				cartItem.Quantity(),
			)
			if err != nil {
				return err
			}
			items = append(items, orderItem)
		}

		o, err := order.NewOrder(orderCustomerID, items)
		if err != nil {
			return err
		}

		if err := uc.orderRepo.Save(ctx, o); err != nil {
			return err
		}

		// イベントの配信は Save が成功した「後」に行う。Save より前や
		// Save と同時に配信してしまうと、Save が失敗した場合に「実際には
		// 起きなかった注文確定」を配信してしまう不整合が生じる
		// （詳細は shared.AggregateBase・order/README.md を参照）。
		//
		// なお、このユースケース全体が 1 つのトランザクション（txManager.Do
		// の fn）の中にあるため、ここで呼ばれる Publish（同期のインメモリ
		// バス経由で ClearCartOnOrderPlaced.Handle を呼ぶ）もまだコミット前の
		// 同一トランザクションに参加する。これにより「注文は確定したが
		// カートは空にならなかった」というような中途半端な状態を防げる。
		// 実運用で Outbox パターンに切り替える場合は、この Publish は
		// 「Outbox テーブルへの INSERT」に置き換わる想定である。
		if err := uc.publisher.Publish(ctx, o.PullEvents()...); err != nil {
			return err
		}

		output = PlaceOrderOutput{OrderID: o.ID().String()}
		return nil
	})
	if err != nil {
		return PlaceOrderOutput{}, err
	}
	return output, nil
}
