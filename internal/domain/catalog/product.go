package catalog

import (
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// 名前・説明文の最大文字数。DB カラム長やクライアント表示上の制約を
// 想定した業務ルールとしてドメイン層に置く（DB 側の VARCHAR 長と一致させる
// 運用を想定しているが、あくまで「業務的に許される長さ」を表現するもので
// あり、DB の都合をドメインに持ち込んでいるわけではない）。
const (
	maxProductNameLength        = 100
	maxProductDescriptionLength = 1000
)

// Product は商品カタログの集約ルート（Aggregate Root）である。
//
// この集約が守る不変条件（Invariant）は次の 2 つだけ:
//  1. name は必須かつ maxProductNameLength 文字以内である
//  2. description は maxProductDescriptionLength 文字以内である
//
// （price の非負性は Money 自身が保証するため、ここで重複してチェックしない）
//
// 商品は「カートに入れられる」「注文される」など他コンテキストから参照
// されるが、それらのコンテキストは ProductID という識別子だけを保持し、
// Product 集約そのものを取り込まない。これは集約境界を小さく保ち、
// コンテキスト間の結合を弱くするための意図的な設計判断である
// （詳細は README.md を参照）。
//
// なお本サンプルの catalog コンテキストは他コンテキストへ通知すべき
// ような重要な出来事（ドメインイベント）を持たないため、
// shared.AggregateBase は埋め込んでいない。イベントが必要になった時点
// （例: 値下げをカート側に通知したい等）で初めて追加すればよい。
type Product struct {
	id          ProductID
	name        string
	description string
	price       shared.Money
}

// NewProduct は新規商品を登録する際に使うコンストラクタである。
// 「生成」と「再構築」を分けている理由は ReconstructProduct のコメントを参照。
func NewProduct(name, description string, price shared.Money) (*Product, error) {
	if err := validateProductName(name); err != nil {
		return nil, err
	}
	if err := validateProductDescription(description); err != nil {
		return nil, err
	}
	return &Product{
		id:          GenerateProductID(),
		name:        name,
		description: description,
		price:       price,
	}, nil
}

// ReconstructProduct は永続化層（リポジトリ）が DB から読み出した値を
// 使って Product を復元するための関数である。
//
// なぜ NewProduct とは別に用意するのか？
// NewProduct は「これから新しく生まれる商品」に対する業務ルール（入力
// バリデーション）を課す入口である。一方、DB から読み出す値は過去に
// NewProduct のバリデーションを通過して保存されたものであり、再度同じ
// チェックを課す必要はない（むしろ、将来バリデーションルールを厳しく
// 変更した場合に、過去のデータが読み込めなくなるという事故を防げる）。
// この「生成コンストラクタ」と「再構築コンストラクタ」の分離は、DDD の
// 実装パターンとしてよく使われる。
func ReconstructProduct(id ProductID, name, description string, price shared.Money) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
	}
}

// ID は商品の識別子を返す。
func (p *Product) ID() ProductID {
	return p.id
}

// Name は商品名を返す。
func (p *Product) Name() string {
	return p.name
}

// Description は商品説明を返す。
func (p *Product) Description() string {
	return p.description
}

// Price は現在の価格を返す。
func (p *Product) Price() shared.Money {
	return p.price
}

// ChangePrice は商品の価格を変更する。
//
// 「同じ価格への変更を拒否する」という一見些細なルールをここに置いているのは、
// ドメインメソッドが単なる setter ではなく、業務的に意味のある操作である
// ことを示す教材的な例である。実務では「変更履歴を残したい」「無意味な
// 更新でイベントを発火させたくない」といった理由でこの種のガードが
// 実際に使われる。
func (p *Product) ChangePrice(newPrice shared.Money) error {
	if p.price.Equals(newPrice) {
		return shared.NewDomainRuleError("catalog: new price must differ from current price (%s)", p.price.String())
	}
	p.price = newPrice
	return nil
}

// validateProductName は商品名の業務ルールを検証する。
func validateProductName(name string) error {
	if name == "" {
		return shared.NewDomainRuleError("catalog: product name must not be empty")
	}
	if len([]rune(name)) > maxProductNameLength {
		return shared.NewDomainRuleError("catalog: product name must be %d characters or fewer", maxProductNameLength)
	}
	return nil
}

// validateProductDescription は商品説明の業務ルールを検証する。
func validateProductDescription(description string) error {
	if len([]rune(description)) > maxProductDescriptionLength {
		return shared.NewDomainRuleError("catalog: product description must be %d characters or fewer", maxProductDescriptionLength)
	}
	return nil
}
