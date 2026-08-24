package order

import (
	"context"
	"errors"
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincart "github.com/almondoo/golang-ddd-sample/internal/domain/cart"
	domaincatalog "github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	domaincoupon "github.com/almondoo/golang-ddd-sample/internal/domain/coupon"
	domaincustomer "github.com/almondoo/golang-ddd-sample/internal/domain/customer"
	domaininventory "github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	domainorder "github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// ユースケース層とドメイン層でパッケージ名が重なるため、ドメイン側に
// domain 接頭辞の別名を付ける（domaincart / domaincatalog / domainorder /
// domaincustomer / domaininventory / domaincoupon）。
// これは internal/application/usecase 配下の各パッケージ（catalog/cart/order/...）が
// internal/domain 配下の同名パッケージと衝突するために本サンプル全体で
// 統一している命名規約であり、このファイルのように複数コンテキストの
// ドメインパッケージをまとめて import する場合に特に効果を発揮する。

// PlaceOrderInput は注文確定ユースケースの入力である。
type PlaceOrderInput struct {
	CustomerID string
	// CouponCode は適用したいクーポンのコードである（任意）。
	// 空文字列であればクーポンを適用しない。
	CouponCode string
}

// PlaceOrderOutput は注文確定ユースケースの出力である。
type PlaceOrderOutput struct {
	OrderID string
}

// PlaceOrderUseCase は「カートの中身をもとに注文を確定する」ユースケースである。
//
// 【このユースケースは複数コンテキストをまたぐオーケストレーション】
// order ドメインパッケージは cart / catalog / customer / inventory / coupon を
// import しないと決めた（境界づけられたコンテキストの自律性）。しかし現実の
// 業務としては「顧客の実在確認」「カートの中身」「その時点の商品価格」
// 「在庫の引当」「クーポンの適用」のすべてを 1 つの注文確定操作の中で
// 整合させなければならない。この「複数コンテキストにまたがる読み取りと
// 調整」を引き受けるのがアプリケーション層の役割である。各集約の
// 不変条件はそれぞれの集約（Customer / Stock / Coupon / Order）自身が守り、
// ユースケースは「どの順番で・どの集約を呼ぶか」という手順だけを
// 組み立てる。
type PlaceOrderUseCase struct {
	orderRepo     domainorder.Repository
	cartRepo      domaincart.Repository
	catalogRepo   domaincatalog.Repository
	customerRepo  domaincustomer.Repository
	inventoryRepo domaininventory.Repository
	couponRepo    domaincoupon.Repository
	txManager     tx.Manager
}

// NewPlaceOrderUseCase は PlaceOrderUseCase を生成する。
func NewPlaceOrderUseCase(
	orderRepo domainorder.Repository,
	cartRepo domaincart.Repository,
	catalogRepo domaincatalog.Repository,
	customerRepo domaincustomer.Repository,
	inventoryRepo domaininventory.Repository,
	couponRepo domaincoupon.Repository,
	txManager tx.Manager,
) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		orderRepo:     orderRepo,
		cartRepo:      cartRepo,
		catalogRepo:   catalogRepo,
		customerRepo:  customerRepo,
		inventoryRepo: inventoryRepo,
		couponRepo:    couponRepo,
		txManager:     txManager,
	}
}

// Execute は注文確定ユースケースを実行する。
//
// 処理の流れ:
//  1. 顧客が実在することを確認する（customer コンテキスト）
//  2. 対象顧客のカートを読み込む（存在しない・空である場合はドメインルール違反）
//  3. カート明細ごとに catalog から現在の商品情報（名前・価格）を取得して
//     OrderItem のスナップショットを組み立て、同時に inventory の在庫を引き当てる
//  4. Order 集約を新規生成する
//  5. クーポンコードが指定されていれば、クーポンを 1 回消費し割引額を注文に適用する
//  6. Order を永続化する
//  7. 永続化が成功した後にカートを空にして保存する
func (uc *PlaceOrderUseCase) Execute(ctx context.Context, in PlaceOrderInput) (PlaceOrderOutput, error) {
	cartCustomerID, err := domaincart.NewCustomerID(in.CustomerID)
	if err != nil {
		return PlaceOrderOutput{}, err
	}
	orderCustomerID, err := domainorder.NewCustomerID(in.CustomerID)
	if err != nil {
		return PlaceOrderOutput{}, err
	}
	customerID, err := domaincustomer.NewCustomerID(in.CustomerID)
	if err != nil {
		return PlaceOrderOutput{}, err
	}

	var output PlaceOrderOutput
	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		// 1回だけ取得して注文時刻とクーポン判定で同じ時刻を使う。
		now := time.Now()

		// 1. 顧客検証: 実在しない顧客の注文は成立しない。
		// customer コンテキストへの参照は「注文を組み立てる」ためのオーケストレーション
		// の一部であり、order ドメインパッケージが customer を import しないという
		// 方針とは矛盾しない（詳細は order/README.md を参照）。
		if _, err := uc.customerRepo.FindByID(ctx, customerID); err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				return shared.NewDomainRuleError("order: customer %s not found", in.CustomerID)
			}
			return err
		}

		// 2. カートの読み込み。
		c, err := uc.cartRepo.FindByCustomerID(ctx, cartCustomerID)
		if err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				// カートがまだ 1 度も作られていない状態は「注文できるものが
				// 何もない」という業務ルール違反として扱う。404（リソースが
				// 存在しない）ではなく 422（ルール違反）として表現したいため、
				// ここで shared.NewDomainRuleError に変換する。
				return shared.NewDomainRuleError("order: cart is empty; add items before placing an order")
			}
			return err
		}
		if c.IsEmpty() {
			return shared.NewDomainRuleError("order: cart is empty; add items before placing an order")
		}

		items := make([]domainorder.OrderItem, 0, len(c.Items()))
		for _, cartItem := range c.Items() {
			productID, err := domaincatalog.NewProductID(cartItem.ProductID().String())
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
					return shared.NewDomainRuleError("order: product %s not found in catalog", productID)
				}
				return err
			}

			// 3. 在庫の引当: 注文確定と同時に「取り置き」を行う（実在庫はまだ
			// 減らさない）。出荷（ShipOrderUseCase）のタイミングで初めて実在庫が
			// 減る（inventory.Stock.ConsumeReserved を参照）。
			inventoryProductID, err := domaininventory.NewProductID(cartItem.ProductID().String())
			if err != nil {
				return err
			}
			stock, err := uc.inventoryRepo.FindByProductID(ctx, inventoryProductID)
			if err != nil {
				if errors.Is(err, shared.ErrNotFound) {
					return shared.NewDomainRuleError("order: stock for product %s is not registered", inventoryProductID)
				}
				return err
			}
			if err := stock.Reserve(cartItem.Quantity()); err != nil {
				// 在庫不足は Stock.Reserve がドメインルール違反として返すので
				// そのまま呼び出し元へ伝播させる。
				return err
			}
			if err := uc.inventoryRepo.Save(ctx, stock); err != nil {
				return err
			}

			orderItem, err := domainorder.NewOrderItem(
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

		// 4. Order 集約の新規生成。
		o, err := domainorder.NewOrder(orderCustomerID, items, now)
		if err != nil {
			return err
		}

		// 5. クーポン適用（任意）。
		if in.CouponCode != "" {
			code, err := domaincoupon.NewCouponCode(in.CouponCode)
			if err != nil {
				return err
			}
			cp, err := uc.couponRepo.FindByCode(ctx, code)
			if err != nil {
				if errors.Is(err, shared.ErrNotFound) {
					return shared.NewDomainRuleError("order: invalid coupon code %s", in.CouponCode)
				}
				return err
			}

			// クーポンの有効期限・利用回数上限の判定は Coupon.Use が担う
			// ドメインルールである。time.Now() をここ（アプリケーション層）で
			// 呼び出すのは、ドメイン層を「時刻に依存しない決定的なコード」に
			// 保つための一貫した方針である（coupon.Coupon.IsExpired のコメントを参照）。
			// 上で取得した now を注文確定時刻と共用し、同一トランザクション内で
			// 時刻がぶれないようにする。
			if err := cp.Use(now); err != nil {
				return err
			}

			total, err := o.TotalAmount()
			if err != nil {
				return err
			}
			discount, err := cp.DiscountFor(total)
			if err != nil {
				return err
			}

			if err := uc.couponRepo.Save(ctx, cp); err != nil {
				return err
			}
			if err := o.ApplyDiscount(code.String(), discount); err != nil {
				return err
			}
		}

		// 6. Order の永続化。
		if err := uc.orderRepo.Save(ctx, o); err != nil {
			return err
		}

		// 以前はここでドメインイベント(OrderPlaced)を発行し、カート側の
		// ハンドラが購読して空にする方式だったが、仕組みの理解コストを
		// 下げるため直接呼び出しに変更した。application 層での直接呼び出しは
		// order→cart の依存を生む(結合度は上がる)一方、処理の流れが一目で
		// 追える。コンテキストをまたぐ反応が増えてきたらドメインイベント+
		// 購読へ戻すのが定石。
		//
		// この呼び出しは txManager.Do の fn の中、つまり Order の Save と
		// 同一トランザクションで行われるため、「注文は確定したがカートは
		// 空にならなかった」という中途半端な状態は起こり得ない。
		c.Clear()
		if err := uc.cartRepo.Save(ctx, c); err != nil {
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
