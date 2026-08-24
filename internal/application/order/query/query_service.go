package query

import "context"

// OrderQueryService は注文の参照専用の問い合わせを表すポートである。
//
// order.Repository（ドメイン層）が「集約を読み書きする」ためのインター
// フェースであるのに対し、こちらは「画面表示に必要な形へ加工済みの
// データを返す」ためのインターフェースである。cart / catalog と同じ
// 軽量 CQRS の適用例であり、実装（infrastructure 層）は Order 集約を
// 経由せず、GORM で order・order_items テーブルを直接読んで DTO を組み立てる。
type OrderQueryService interface {
	// FindByID は id に対応する注文を DTO として取得する。
	// 見つからない場合は shared.ErrNotFound をラップしたエラーを返す。
	//
	// カートのクエリとは異なり「存在しない場合は空を返す」のではなく
	// エラーとして扱う。注文はカートと違って「まだ作られていない」状態が
	// 存在せず（GET /orders/{id} は必ず特定の注文を指す）、指定した ID の
	// 注文が存在しないことは呼び出し側の誤り、または削除済みの可能性を
	// 示す正当なエラーだからである。
	FindByID(ctx context.Context, id string) (*OrderDTO, error)
}
