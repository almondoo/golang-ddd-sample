// Package order は order コンテキストのユースケース（アプリケーションサービス）を
// まとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - place_order.go / pay_order.go / ship_order.go / cancel_order.go … 書き込み系
//     （コマンド）のユースケース。order.Repository と tx.Manager に依存し、集約を
//     読み込み・操作し・保存する（place_order.go のみ cart / catalog の
//     Repository にも依存する。理由は同ファイルのコメントを参照）。
//   - get_order.go                                                    … 読み取り系
//     （クエリ）のユースケース。OrderQueryService にのみ依存し、ドメイン層を
//     経由せず DTO を直接返す。
//
// かつては command / query という別パッケージに分割していたが、パッケージを
// 分けるほどの独立性はなく、コマンドとクエリの区別はパッケージではなく
// ファイル名と依存する型（Repository+tx.Manager か、QueryService のみか）で
// 表現する方針に統合した。
package order

import "time"

// OrderDTO は注文を画面表示用に平坦化したデータ転送オブジェクトである。
//
// catalog / cart の DTO と同じ理由で、ドメインの Order 集約とは別に
// 用意している。ドメイン集約は不変条件を守るための構造を持つが、
// クエリ側はそれを一切気にせず「表示に必要な形」に最適化してよい
// （CQRS の考え方）。TotalAmount はドメイン側では毎回明細から計算する
// メソッドだが、DTO では計算済みの値をそのままフィールドとして持たせる。
type OrderDTO struct {
	ID          string         `json:"id"`
	CustomerID  string         `json:"customerId"`
	Status      string         `json:"status"`
	TotalAmount int64          `json:"totalAmount"`
	PlacedAt    time.Time      `json:"placedAt"`
	Items       []OrderItemDTO `json:"items"`
	// CouponCode は適用済みクーポンのコードである。空文字列は未適用を表す。
	CouponCode string `json:"couponCode"`
	// DiscountAmount は適用済みの割引額である。未適用の場合は 0。
	DiscountAmount int64 `json:"discountAmount"`
	// PayableAmount は実際に支払うべき金額（TotalAmount - DiscountAmount）である。
	PayableAmount int64 `json:"payableAmount"`
}

// OrderItemDTO は注文内の 1 明細を画面表示用に表したものである。
// ドメインの OrderItem と同様、ここでの値は注文確定時点のスナップショットで
// あり、確定後に catalog 側の商品名・価格が変わっても影響を受けない。
type OrderItemDTO struct {
	ProductID       string `json:"productId"`
	ProductName     string `json:"productName"`
	UnitPriceAmount int64  `json:"unitPriceAmount"`
	Quantity        int    `json:"quantity"`
	Subtotal        int64  `json:"subtotal"`
}
