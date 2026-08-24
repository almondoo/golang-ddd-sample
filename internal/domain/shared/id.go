package shared

import "github.com/google/uuid"

// NewID は新しい UUID 文字列を生成する。
//
// 各コンテキストは、この共有関数をラップして CartID や OrderID のような
// 型付き ID（例: type CartID string）を定義することを想定している。
// ここで生成ロジックを一箇所に集約しておくことで、将来 UUID v4 から
// ULID 等に切り替える場合も変更箇所を最小化できる。
func NewID() string {
	return uuid.NewString()
}
