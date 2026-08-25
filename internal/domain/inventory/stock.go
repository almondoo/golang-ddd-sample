package inventory

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// Stock は商品ごとの在庫を表す集約ルート（Aggregate Root）である。
//
// 集約の識別子には ProductID をそのまま使っている（サロゲートの StockID を
// 別途発行していない）。これは「1 商品につき在庫レコードは常にちょうど 1 つ」
// という業務ルールをそのままモデルに反映させた設計判断である（cart.Cart が
// CustomerID をそのまま識別子に使っているのと同じ考え方）。
//
// quantity（実在庫）と reserved（引当済み数量）を分けて持つ「引当モデル」を
// 採用している。注文が確定した時点では、まだ商品を出荷（実在庫の減少）
// していなくても「他の注文には割り当てられない数量」を先に確保しておく
// 必要がある。これを reserved として quantity とは別に管理することで、
//   - 引当（Reserve）: 注文確定時に「取り置き」を行う。実在庫はまだ減らさない。
//   - 解放（Release）: 注文キャンセル時に取り置きを解除する。実在庫は変えない。
//   - 消込（ConsumeReserved）: 出荷時に取り置き分を実際に払い出す。
//     このとき初めて実在庫（quantity）と引当（reserved）の両方を減らす。
//
// という 3 つの操作を明確に区別できる。
type Stock struct {
	productID ProductID
	quantity  int
	reserved  int
	// version は楽観ロック（optimistic lock）のための版数である。
	//
	// 並行するトランザクション同士の競合（同じ在庫を同時に読んで同時に
	// 更新する lost update）は、Reserve 等の不変条件チェックだけでは
	// 防げない。それらのチェックは「1つの Stock インスタンスに対する
	// 操作」の整合性しか守れないためである。version をドメインモデル自身に
	// 持たせているのは、「この集約が読み込まれた時点の版数」を業務ロジックと
	// 一体で扱えるようにするため、つまり集約の整合性境界（consistency
	// boundary）をアプリケーション層やインフラ層に漏らさず、Stock 自身の
	// 状態として管理するためである。永続化層はこの version を使って
	// 「読んでから保存するまでの間に他のトランザクションが更新していないか」を
	// UPDATE ... WHERE version = ? という条件付き更新で検出し、競合していれば
	// shared.ErrConflict を返す。
	//
	// version == 0 は「まだ一度も永続化されていない」ことを表す規約である
	// （INSERT 時に永続化層が version = 1 を採番する）。この規約により、
	// リポジトリ実装は「DB に行が存在するかどうか」を別途問い合わせて
	// INSERT/UPDATE を判定する必要がなくなり、集約自身が持つ version だけで
	// 判定できる（find-or-create のような「先に存在確認してから分岐する」
	// パスに潜む TOCTOU（Time-Of-Check-Time-Of-Use）競合を避けられる）。
	version int
}

// NewStock は新しい在庫を初期数量とともに生成する。
// 新規作成した集約の版数は 0（未永続化）から始める。version フィールドの
// コメントに記した規約のとおり、実際に version = 1 が採番されるのは
// リポジトリの Save が INSERT を行った時点である。
func NewStock(productID ProductID, quantity int) (*Stock, error) {
	if quantity < 0 {
		return nil, shared.NewDomainRuleError("inventory: quantity must not be negative, got %d", quantity)
	}
	return &Stock{productID: productID, quantity: quantity, reserved: 0, version: 0}, nil
}

// ReconstructStock は永続化層から読み込んだデータをもとに Stock を再構築する。
//
// NewStock は「在庫を新規に立ち上げる」という業務上の意図を表すのに対し、
// ReconstructStock は「すでに存在する在庫を DB から復元する」という別の
// 意図を表す。DB の値は過去に検証済みという前提のもと、検証を行わない
// 「素通し」の関数として分離している（cart.ReconstructCart と同じ判断）。
// リポジトリ実装（infrastructure 層）からのみ呼ばれることを想定している。
//
// version は DB から読み出した時点の版数をそのまま渡す。この Stock が
// 後で Save される際、リポジトリはこの version を使って
// 「保存しようとしている内容が読み込み時点から更新されていないか」を
// 楽観ロックで検証する。
func ReconstructStock(productID ProductID, quantity, reserved, version int) *Stock {
	return &Stock{productID: productID, quantity: quantity, reserved: reserved, version: version}
}

// ProductID はこの在庫が対象とする商品の ID を返す。
func (s *Stock) ProductID() ProductID {
	return s.productID
}

// Quantity は実在庫数を返す。
func (s *Stock) Quantity() int {
	return s.quantity
}

// Reserved は引当済み数量を返す。
func (s *Stock) Reserved() int {
	return s.reserved
}

// Available は引当可能な残数量（実在庫 - 引当済み）を返す。
func (s *Stock) Available() int {
	return s.quantity - s.reserved
}

// Version は楽観ロックのための版数を返す。
// リポジトリ実装が Save 時の条件付き更新（WHERE version = ?）に使う。
func (s *Stock) Version() int {
	return s.version
}

// SetQuantity は実在庫数を n に設定する（入荷・棚卸し等による在庫数の更新）。
//
// n は 0 以上かつ現在の引当済み数量（reserved）以上でなければならない。
// 引当済みを下回る在庫設定を許してしまうと、reserved <= quantity という
// 不変条件が破れ、「引当した数量より少ない実在庫しかない」という不整合な
// 状態になってしまうため拒否する。
func (s *Stock) SetQuantity(n int) error {
	if n < 0 {
		return shared.NewDomainRuleError("inventory: quantity must not be negative, got %d", n)
	}
	if n < s.reserved {
		return shared.NewDomainRuleError("inventory: quantity (%d) must not be less than reserved (%d)", n, s.reserved)
	}
	s.quantity = n
	return nil
}

// Reserve は n 個の在庫を引き当てる（注文確定時に呼ばれることを想定）。
//
// 在庫引当は「reserved <= quantity」という不変条件を集約自身が守る、
// という設計の核心部分である。Available()（= quantity - reserved）を
// 超える引当は許可しない。
//
// ただしこの不変条件が守るのは「1つの Stock インスタンスに対する操作」の
// 整合性まで。並行するトランザクション同士の競合（同じ在庫を同時に読んで
// 同時に引き当てる lost update）は集約だけでは防げない。この Stock は
// version フィールドを持ち、リポジトリ実装（infrastructure 層）が
// Save 時に WHERE version = ? の条件付き更新で楽観ロックを行うことで、
// 「読み込んでから保存するまでの間に他のトランザクションが先に更新した」
// 競合を検出し shared.ErrConflict を返す。呼び出し側（ユースケース）は
// この場合、在庫を読み直してから Reserve をやり直す（リトライする）ことを
// 想定している。
func (s *Stock) Reserve(n int) error {
	if n < 1 {
		return shared.NewDomainRuleError("inventory: reserve quantity must be at least 1, got %d", n)
	}
	if n > s.Available() {
		return shared.NewDomainRuleError("inventory: insufficient stock: requested %d, available %d", n, s.Available())
	}
	s.reserved += n
	return nil
}

// Release は n 個の引当を解除する（注文キャンセル時に呼ばれることを想定）。
// 実在庫（quantity）は変えず、引当済み数量（reserved）だけを減らす。
func (s *Stock) Release(n int) error {
	if n < 1 {
		return shared.NewDomainRuleError("inventory: release quantity must be at least 1, got %d", n)
	}
	if n > s.reserved {
		return shared.NewDomainRuleError("inventory: cannot release %d, only %d reserved", n, s.reserved)
	}
	s.reserved -= n
	return nil
}

// ConsumeReserved は引当済みの n 個を消込む（出荷時に呼ばれることを想定）。
//
// 出荷は「取り置きしていた在庫を実際に払い出す」操作であるため、
// 実在庫（quantity）と引当済み数量（reserved）の両方を同時に減らす。
// Release との違いはここにある: Release は「取り置きをやめる」だけで
// 実在庫は減らないが、ConsumeReserved は「取り置きしていたものを
// 実際に出す」ため実在庫も減る。
func (s *Stock) ConsumeReserved(n int) error {
	if n < 1 {
		return shared.NewDomainRuleError("inventory: consume quantity must be at least 1, got %d", n)
	}
	if n > s.reserved {
		return shared.NewDomainRuleError("inventory: cannot consume %d, only %d reserved", n, s.reserved)
	}
	s.quantity -= n
	s.reserved -= n
	return nil
}
