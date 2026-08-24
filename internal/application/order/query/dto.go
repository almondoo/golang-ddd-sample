package query

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
