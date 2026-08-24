package cart

import "context"

// Repository はカート集約の永続化を抽象化するポート（インターフェース）である。
//
// ドメイン層がこのインターフェースを定義し、実装（GORM を使った具体的な
// SQL 操作等）はインフラストラクチャ層が提供する。依存の向きは
// インフラストラクチャ層 → ドメイン層 となり、ドメイン層は「どうやって
// 保存するか」を一切知らないまま「保存できる」という契約だけに依存できる。
// これが依存性逆転の原則（Dependency Inversion Principle）であり、
// ヘキサゴナルアーキテクチャにおける「ポート」に相当する。
//
// メソッドを Cart 単位（集約単位）で用意しているのも重要な点である。
// CartItem 単体を取得・保存する API を用意しないのは、集約の外から
// 内部エンティティを個別に触れるようにしてしまうと、Cart が守るべき
// 不変条件（数量上限や明細数の上限）を経由せずに状態を変えられてしまう
// 経路が生まれてしまうためである。
type Repository interface {
	// FindByCustomerID は指定した顧客のカートを取得する。
	// 該当するカートが存在しない場合は shared.ErrNotFound をラップしたエラーを返す。
	FindByCustomerID(ctx context.Context, id CustomerID) (*Cart, error)
	// Save はカートの現在の状態を永続化する。
	// 集約全体を単位として保存することで、明細の一部だけが更新されて
	// 集約全体としては不整合な状態になる、という事態を避ける。
	Save(ctx context.Context, c *Cart) error
}
