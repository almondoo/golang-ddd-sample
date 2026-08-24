// Package coupon は coupon コンテキストのユースケース（アプリケーションサービス）を
// まとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - issue_coupon.go … 書き込み系（コマンド）のユースケース。
//     domaincoupon.Repository と tx.Manager に依存し、集約を生成・検証し・保存する。
//   - get_coupon.go    … 読み取り系（クエリ）のユースケース。
//     CouponQueryService にのみ依存し、ドメイン層を経由せず DTO を直接返す。
//
// コマンドとクエリの区別はパッケージではなく、ファイル名と依存する型
// （Repository+tx.Manager か、QueryService のみか）で表現する方針を
// cart/catalog コンテキストと統一している。
package coupon

import "time"

// CouponDTO はクーポンを画面表示・API レスポンス用に平坦化したデータ転送
// オブジェクト（DTO）である。
//
// ドメインの Coupon 集約とは別に DTO を用意しているのは、クエリ側（読み取り）
// とコマンド側（書き込み）で最適な形が異なるためである（CQRS の考え方）。
// amount / ratePercent はどちらか一方だけが意味を持つという Coupon 集約の
// 内部事情をそのまま反映し、DTO 側でも両方のフィールドを平坦に保持する
// （割引方式の判定は DiscountType フィールドを見て行う）。
type CouponDTO struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	DiscountType string    `json:"discountType"`
	Amount       int64     `json:"amount"`
	RatePercent  int       `json:"ratePercent"`
	ExpiresAt    time.Time `json:"expiresAt"`
	UsageLimit   int       `json:"usageLimit"`
	UsedCount    int       `json:"usedCount"`
}
