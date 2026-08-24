package persistence

import (
	"time"

	"github.com/almondoo/golang-ddd-sample/internal/domain/order"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// OrderModel は orders テーブルに対応する GORM 用の構造体である。
//
// ProductModel と同様、ドメインの Order 集約とは意図的に別の型として
// 定義している。GORM のタグはあくまで永続化の都合であり、ドメイン層の
// 型にこれを持ち込むと「ドメインが ORM の存在を知っている」ことになり、
// 依存性のルールに反する。なお、created_at/updated_at のような監査用
// カラムは持たない（理由は catalog_model.go の ProductModel を参照）。
// PlacedAt はドメインの「注文確定時刻」という業務上の概念であり、
// この監査列の話とは別物として扱う。
type OrderModel struct {
	ID         string `gorm:"primaryKey"`
	CustomerID string
	Status     string
	PlacedAt   time.Time
}

// TableName は GORM に対して物理テーブル名を明示する。
func (OrderModel) TableName() string {
	return "orders"
}

// OrderItemModel は order_items テーブルに対応する GORM 用の構造体である。
//
// 複合主キー（OrderID, ProductID）を使っているのは cart_items と同じ
// 理由（同一注文・同一商品の明細は 1 行しか存在しない、という不変条件を
// DB スキーマのレベルでも保証するため）である。ただし cart_items とは
// 異なり、ProductName/UnitPriceAmount/Currency を独自に持つ。これは
// OrderItem がスナップショットであるという設計をそのまま反映しており、
// products テーブルを JOIN しなくても注文明細単独で当時の内容が
// 復元できることを保証する（詳細は order/README.md を参照）。
type OrderItemModel struct {
	OrderID         string `gorm:"primaryKey"`
	ProductID       string `gorm:"primaryKey"`
	ProductName     string
	UnitPriceAmount int64
	Currency        string
	Quantity        int
}

// TableName は GORM に対して物理テーブル名を明示する。
func (OrderItemModel) TableName() string {
	return "order_items"
}

// orderFromModels は orders 行 + order_items 行群から Order 集約を復元する。
// DB に保存されている値は過去にドメイン層の検証を通過済みという前提の
// もと、ReconstructOrder / ReconstructOrderItem（検証を行わない再構築
// コンストラクタ）を使う。ただし Status だけは NewStatus で検証する。
// これは「取り得る 4 状態のいずれかである」という DB 側では保証できない
// 制約を、復元のたびに確認するための最終防衛線である。
func orderFromModels(model OrderModel, itemModels []OrderItemModel) (*order.Order, error) {
	orderID, err := order.NewOrderID(model.ID)
	if err != nil {
		return nil, err
	}
	customerID, err := order.NewCustomerID(model.CustomerID)
	if err != nil {
		return nil, err
	}
	status, err := order.NewStatus(model.Status)
	if err != nil {
		return nil, err
	}

	items := make([]order.OrderItem, 0, len(itemModels))
	for _, im := range itemModels {
		unitPrice, err := shared.NewMoney(im.UnitPriceAmount, shared.Currency(im.Currency))
		if err != nil {
			return nil, err
		}
		items = append(items, order.ReconstructOrderItem(im.ProductID, im.ProductName, unitPrice, im.Quantity))
	}

	return order.ReconstructOrder(orderID, customerID, items, status, model.PlacedAt), nil
}

// orderModelFromDomain は Order 集約から orders テーブル用のモデルを組み立てる。
func orderModelFromDomain(o *order.Order) OrderModel {
	return OrderModel{
		ID:         o.ID().String(),
		CustomerID: o.CustomerID().String(),
		Status:     o.Status().String(),
		PlacedAt:   o.PlacedAt(),
	}
}

// orderItemModelsFromDomain は Order 集約から order_items テーブル用の
// モデル一覧を組み立てる。
func orderItemModelsFromDomain(o *order.Order) []OrderItemModel {
	items := o.Items()
	models := make([]OrderItemModel, 0, len(items))
	for _, item := range items {
		models = append(models, OrderItemModel{
			OrderID:         o.ID().String(),
			ProductID:       item.ProductID(),
			ProductName:     item.ProductName(),
			UnitPriceAmount: item.UnitPrice().Amount(),
			Currency:        string(item.UnitPrice().Currency()),
			Quantity:        item.Quantity(),
		})
	}
	return models
}
