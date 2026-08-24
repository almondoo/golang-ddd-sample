package order

import "github.com/almondoo/golang-ddd-sample/internal/domain/shared"

// Status は注文のライフサイクル上の状態を表す値オブジェクトである。
//
// string を基底型にしているのは DB へそのまま文字列として保存でき、
// かつログ・デバッグ出力でも人間が読める形になるためである。
// 状態遷移そのもののルール（どの状態からどの状態へ遷移してよいか）は
// この型ではなく Order 集約（order.go）が持つ。Status 自身は
// 「取り得る値の集合」と「文字列としての妥当性」だけを守る責務に留める。
type Status string

const (
	// StatusPending は注文が確定した直後の状態である。
	// まだ入金が確認されておらず、発送はできない。
	StatusPending Status = "pending"
	// StatusPaid は入金が確認された状態である。発送が可能になる。
	StatusPaid Status = "paid"
	// StatusShipped は商品が発送された状態である。これ以降の遷移は
	// 想定しない（返品・返金等が必要な場合は本サンプルの範囲外とする）。
	StatusShipped Status = "shipped"
	// StatusCanceled は注文が取り消された状態である。
	// pending または paid からのみ到達でき、一度到達すると
	// そこから他の状態へは遷移できない終端状態である。
	StatusCanceled Status = "canceled"
)

// String は Status を文字列として取り出す。
func (s Status) String() string {
	return string(s)
}

// NewStatus は文字列から Status を生成するコンストラクタである。
//
// 主に永続化層が DB に保存された文字列を Order 集約へ復元する際に使う
// （ReconstructOrder 経由）。DB の内容が何らかの理由で不正な文字列に
// なっていた場合（手動での UPDATE ミス等）を早期に検出するため、
// あらかじめ定義された 4 つの状態のいずれかであることを検証する。
func NewStatus(s string) (Status, error) {
	switch Status(s) {
	case StatusPending, StatusPaid, StatusShipped, StatusCanceled:
		return Status(s), nil
	default:
		return "", shared.NewDomainRuleError("order: unknown status %q", s)
	}
}
