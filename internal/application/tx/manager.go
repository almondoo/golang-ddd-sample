package tx

import "context"

// Manager はトランザクション境界を表すポート（インターフェース）である。
//
// ユースケース（アプリケーションサービス）は「1 回のユースケース実行 = 1
// トランザクション」という原則を守るために、この Manager 経由で境界を明示する。
// fn の中で複数のリポジトリ呼び出し（例: カートを読み込み、更新し、保存する）を
// 行っても、それらはすべて同一トランザクションの中で実行される。
//
// 実装はインフラストラクチャ層（persistence.TxManager）が GORM の
// db.Transaction を使って提供する。アプリケーション層は GORM の存在を
// 一切知らずに「トランザクションで囲む」という意図だけを表現できる。
//
// fn が返すエラーが nil でなければロールバック、nil であればコミットする
// という契約は実装側の責務である。
type Manager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
