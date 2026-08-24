package shipping

import "github.com/almondoo/golang-ddd-sample/internal/domain/shared"

// Status は配送のライフサイクル上の状態を表す値オブジェクトである。
//
// string を基底型にしているのは DB へそのまま文字列として保存でき、
// かつログ・デバッグ出力でも人間が読める形になるためである（order.Status
// と同じ設計意図）。状態遷移そのもののルール（どの状態からどの状態へ
// 遷移してよいか）はこの型ではなく Shipment 集約（shipment.go）が持つ。
// Status 自身は「取り得る値の集合」と「文字列としての妥当性」だけを
// 守る責務に留める。
type Status string

const (
	// StatusPreparing は配送の準備中（発送前）の状態である。
	// Shipment 生成直後はこの状態から始まる。
	StatusPreparing Status = "preparing"
	// StatusShipped は商品が発送された状態である。
	StatusShipped Status = "shipped"
	// StatusDelivered は商品が配送先に到着した状態である。これ以降の
	// 遷移は想定しない（本サンプルの範囲外とする）。
	StatusDelivered Status = "delivered"
)

// String は Status を文字列として取り出す。
func (s Status) String() string {
	return string(s)
}

// NewStatus は文字列から Status を生成するコンストラクタである。
//
// 主に永続化層が DB に保存された文字列を Shipment 集約へ復元する際に使う
// （ReconstructShipment 経由）。DB の内容が何らかの理由で不正な文字列に
// なっていた場合（手動での UPDATE ミス等）を早期に検出するため、
// あらかじめ定義された 3 つの状態のいずれかであることを検証する。
func NewStatus(s string) (Status, error) {
	switch Status(s) {
	case StatusPreparing, StatusShipped, StatusDelivered:
		return Status(s), nil
	default:
		return "", shared.NewDomainRuleError("shipping: unknown status %q", s)
	}
}
