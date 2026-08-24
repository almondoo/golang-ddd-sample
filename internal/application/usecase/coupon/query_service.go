package coupon

import "context"

// CouponQueryService はクーポンの参照専用の問い合わせを表すポートである。
//
// domaincoupon.Repository（ドメイン層）が「集約を読み書きする」ための
// インターフェースであるのに対し、こちらは「画面表示・API レスポンスに
// 必要な形へ加工済みのデータを返す」ためのインターフェースである。
// 両者を分けているのは CQRS（コマンドとクエリの責務分離）の考え方に基づく。
type CouponQueryService interface {
	// FindByCode は指定コードのクーポンを DTO として取得する。
	// 該当するクーポンが存在しない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByCode(ctx context.Context, code string) (*CouponDTO, error)
}
