package query

// ProductDTO は問い合わせ（Query）結果としてクライアントに返す
// データ転送オブジェクト（Data Transfer Object）である。
//
// ドメインの Product 集約をそのまま JSON にシリアライズせず、専用の DTO を
// 経由させているのは、「表示に必要な形」と「業務ルールを守るための内部
// 構造」を分離するためである。DTO はドメインの不変条件を一切持たない
// ただのデータの入れ物であり、Product の内部フィールド構成を変更しても
// API の応答形式を独立して保てる。
type ProductDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceAmount int64  `json:"priceAmount"`
	Currency    string `json:"currency"`
}
