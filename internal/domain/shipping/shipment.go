package shipping

import "github.com/almondoo/golang-ddd-sample/internal/domain/shared"

// Shipment は配送を表す集約ルート（Aggregate Root）である。
//
// 【集約境界】
// Shipment は自身の配送先住所（address）と状態（status）だけを持つ
// 小さな集約である。どの注文（Order）に対する配送かは orderID
// （shipping コンテキスト自身の型）で参照するのみで、order 集約への
// 参照や order パッケージへの import は一切行わない。これは order
// 集約が cart / catalog パッケージを import しないのと同じ、コンテキスト
// の自律性を保つための最重要ルールである。
type Shipment struct {
	id      ShipmentID
	orderID OrderID
	// address は配送時点の住所を文字列で写し取ったスナップショットである。
	// 顧客が後で住所を変えても過去の配送記録は変わらない
	// （order.OrderItem がスナップショットである理由と同じ）。
	address string
	status  Status
}

// NewShipment は新規の配送を生成するコンストラクタである。
//
// address は配送先を特定するために必須の情報であるため、空文字列は
// ドメインルール違反として拒否する。
//
// 生成時点の状態は Shipment 集約自身が決める（StatusPreparing で
// 開始する）。呼び出し側が状態を指定できないようにすることで、
// 「生成直後の配送はまだ準備中である」という不変条件を型の外から
// 破れないようにしている。
func NewShipment(orderID OrderID, address string) (*Shipment, error) {
	if address == "" {
		return nil, shared.NewDomainRuleError("shipping: address must not be empty")
	}

	return &Shipment{
		id:      GenerateShipmentID(),
		orderID: orderID,
		address: address,
		status:  StatusPreparing,
	}, nil
}

// ReconstructShipment は永続化層から読み込んだデータをもとに Shipment を
// 再構築する。NewShipment との違いは ReconstructOrder 等と同じ理由による
// （検証済みのデータを前提とし、バリデーションを行わない）。
func ReconstructShipment(id ShipmentID, orderID OrderID, address string, status Status) *Shipment {
	return &Shipment{
		id:      id,
		orderID: orderID,
		address: address,
		status:  status,
	}
}

// ID は配送の識別子を返す。
func (s *Shipment) ID() ShipmentID {
	return s.id
}

// OrderID はこの配送が対象とする注文の ID を返す。
func (s *Shipment) OrderID() OrderID {
	return s.orderID
}

// Address は配送先住所のスナップショットを返す。
func (s *Shipment) Address() string {
	return s.address
}

// Status は配送の現在の状態を返す。
func (s *Shipment) Status() Status {
	return s.status
}

// 【状態遷移】
//
//	preparing --MarkShipped--> shipped --MarkDelivered--> delivered
//
// 上記以外の遷移（例: preparing から MarkDelivered、delivered から
// 再度 MarkShipped）はすべてドメインルール違反として拒否する。
// 状態機械（State Machine）の判断を Shipment 集約に閉じ込めることで、
// 「今どの状態からどの状態へ遷移できるか」という業務ルールが Shipment
// の外（アプリケーション層やプレゼンテーション層）に漏れ出さないように
// している（order.Order の Pay/Ship/Cancel と同じ設計方針）。

// MarkShipped は配送を発送済みにし、preparing から shipped へ遷移させる。
func (s *Shipment) MarkShipped() error {
	if s.status != StatusPreparing {
		return shared.NewDomainRuleError("shipping: cannot mark as shipped a shipment in status %q", s.status)
	}
	s.status = StatusShipped
	return nil
}

// MarkDelivered は配送を配達完了にし、shipped から delivered へ遷移させる。
func (s *Shipment) MarkDelivered() error {
	if s.status != StatusShipped {
		return shared.NewDomainRuleError("shipping: cannot mark as delivered a shipment in status %q", s.status)
	}
	s.status = StatusDelivered
	return nil
}
