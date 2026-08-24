package persistence

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/customer"
)

// CustomerModel は customers テーブルに対応する GORM 用の構造体である。
//
// ドメインの Customer 集約とは意図的に別の型として定義している。GORM の
// タグはあくまで永続化の都合であり、ドメイン層の型にこれを持ち込むと
// 「ドメインが ORM の存在を知っている」ことになり、依存性のルールに反する
// （詳細は internal/domain/README.md 「なぜここに GORM を置かないのか」を参照）。
type CustomerModel struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string
}

// TableName は GORM に対して物理テーブル名を明示する。
func (CustomerModel) TableName() string {
	return "customers"
}

// CustomerAddressModel は customer_addresses テーブルに対応する GORM 用の
// 構造体である。
//
// order_items（OrderItemModel）と異なり、Address は集約内で独自の識別子
// （AddressID）を持つ子エンティティであるため、主キーは複合キーではなく
// Address 自身の ID（AddressID）を用いる。CustomerID は「どの顧客に
// 属する住所か」を表す外部キー相当の列として持つ。
type CustomerAddressModel struct {
	ID         string `gorm:"primaryKey"`
	CustomerID string `gorm:"index"`
	PostalCode string
	Prefecture string
	City       string
	Line       string
	IsDefault  bool
}

// TableName は GORM に対して物理テーブル名を明示する。
func (CustomerAddressModel) TableName() string {
	return "customer_addresses"
}

// customerFromModels は customers 行 + customer_addresses 行群から Customer
// 集約を復元する。DB に保存されている値は過去にドメイン層の検証を通過済み
// という前提のもと、ReconstructCustomer / ReconstructAddress（検証を
// 行わない再構築コンストラクタ）を使う。
func customerFromModels(model CustomerModel, addressModels []CustomerAddressModel) (*customer.Customer, error) {
	customerID, err := customer.NewCustomerID(model.ID)
	if err != nil {
		return nil, err
	}

	addresses := make([]customer.Address, 0, len(addressModels))
	for _, am := range addressModels {
		addressID, err := customer.NewAddressID(am.ID)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, customer.ReconstructAddress(
			addressID, am.PostalCode, am.Prefecture, am.City, am.Line, am.IsDefault,
		))
	}

	return customer.ReconstructCustomer(customerID, model.Name, model.Email, addresses), nil
}

// customerModelFromDomain は Customer 集約から customers テーブル用の
// モデルを組み立てる。
func customerModelFromDomain(c *customer.Customer) CustomerModel {
	return CustomerModel{
		ID:    c.ID().String(),
		Name:  c.Name(),
		Email: c.Email(),
	}
}

// customerAddressModelsFromDomain は Customer 集約から customer_addresses
// テーブル用のモデル一覧を組み立てる。
func customerAddressModelsFromDomain(c *customer.Customer) []CustomerAddressModel {
	addresses := c.Addresses()
	models := make([]CustomerAddressModel, 0, len(addresses))
	for _, a := range addresses {
		models = append(models, CustomerAddressModel{
			ID:         a.ID().String(),
			CustomerID: c.ID().String(),
			PostalCode: a.PostalCode(),
			Prefecture: a.Prefecture(),
			City:       a.City(),
			Line:       a.Line(),
			IsDefault:  a.IsDefault(),
		})
	}
	return models
}
