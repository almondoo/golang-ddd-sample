package persistence

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shipping"
)

// ShipmentModel は shipments テーブルに対応する GORM 用の構造体である。
//
// ドメインの Shipment 集約とは意図的に別の型として定義している。
// GORM のタグはあくまで永続化の都合であり、ドメイン層の型にこれを
// 持ち込むと「ドメインが ORM の存在を知っている」ことになり、
// 依存性のルールに反する（ProductModel / OrderModel と同じ理由）。
type ShipmentModel struct {
	ID      string `gorm:"primaryKey"`
	OrderID string `gorm:"index"`
	Address string
	Status  string
}

// TableName は GORM に対して物理テーブル名を明示する。
func (ShipmentModel) TableName() string {
	return "shipments"
}

// toDomain は永続化モデルからドメイン集約を復元する。
// DB から読み出した値は過去にドメイン層のバリデーションを通過済みという
// 前提のもと、ReconstructShipment（検証を行わない再構築コンストラクタ）を
// 使う。ただし Status だけは NewStatus で検証する。これは「取り得る 3
// 状態のいずれかである」という DB 側では保証できない制約を、復元のたびに
// 確認するための最終防衛線である（order_model.go の orderFromModels と
// 同じ方針）。
func (m ShipmentModel) toDomain() (*shipping.Shipment, error) {
	id, err := shipping.NewShipmentID(m.ID)
	if err != nil {
		return nil, err
	}
	orderID, err := shipping.NewOrderID(m.OrderID)
	if err != nil {
		return nil, err
	}
	status, err := shipping.NewStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return shipping.ReconstructShipment(id, orderID, m.Address, status), nil
}

// shipmentModelFromDomain はドメイン集約から永続化モデルを組み立てる。
func shipmentModelFromDomain(s *shipping.Shipment) ShipmentModel {
	return ShipmentModel{
		ID:      s.ID().String(),
		OrderID: s.OrderID().String(),
		Address: s.Address(),
		Status:  s.Status().String(),
	}
}
