package query

// CartDTO はカートを画面表示用に平坦化したデータ転送オブジェクトである。
//
// ドメインの Cart 集約とは別に DTO を用意しているのは、
// クエリ側（読み取り）とコマンド側（書き込み）で最適な形が異なるためである
// （CQRS の考え方）。ドメインモデルは価格を一切知らないが、画面には
// 商品名や金額を表示したい。この落差を埋めるのがクエリ側の責務であり、
// query_service.go の実装は cart_items テーブルと catalog コンテキストが
// 所有する products テーブルを SQL レベルで直接 JOIN して求める。
// これはドメイン層を経由しないため、コンテキスト間の疎結合を破らない
// （あくまで読み取り専用の「参照」であり、catalog のドメインロジックを
// 呼び出しているわけではない点に注意）。
type CartDTO struct {
	CustomerID  string        `json:"customerId"`
	Items       []CartItemDTO `json:"items"`
	TotalAmount int64         `json:"totalAmount"`
}

// CartItemDTO はカート内の 1 明細を画面表示用に表したものである。
type CartItemDTO struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	PriceAmount int64  `json:"priceAmount"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}
