package customer

import (
	"strings"
	"unicode/utf8"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// maxNameRunes は氏名として許容する最大文字数（rune 単位）である。
const maxNameRunes = 50

// maxAddresses は 1 顧客が登録できる住所の上限である。
// 住所帳が無制限に肥大化するのを防ぐ業務上の割り切りであり、
// cart.maxDistinctItems（カート明細の上限）と同種の設計判断である。
const maxAddresses = 5

// Customer は顧客を表す集約ルート（Aggregate Root）である。
//
// 【本サンプルにおける「子エンティティを持つ集約」のショーケース】
// order.Order も cart.Cart も内部に明細（子要素）を持つが、この
// Customer + Address の組み合わせは、DDD が言う「集約（Aggregate）」の
// 本質がもっとも分かりやすく現れるように設計されている。
//
//   - Address は独自の識別子（AddressID）を持つ「子エンティティ」であり、
//     単なる値の集まり（値オブジェクト）ではない。
//   - Customer 集約は「デフォルト住所は、住所が 1 件でもあれば必ず
//     ちょうど 1 つ存在する」という不変条件を持つ。この条件は Address を
//     1 件ずつ独立に検証しても守れない。兄弟にあたる他の住所すべての
//     状態を横断的に見る必要があるからである。
//
// つまり「デフォルト住所は必ず 1 つ」という不変条件は、子エンティティ
// Address 単体の責務では絶対に守れない。この不変条件を守れるのは、
// 全住所を見渡せる立場にある集約ルート Customer だけである。だからこそ
// Address の生成・変更操作はすべて Customer のメソッド
// （AddAddress / ChangeDefaultAddress / RemoveAddress）に集約し、
// 外部から Address を直接 new して Customer に注入するような経路は
// 一切用意しない。
//
// この「集約＝トランザクション整合性を保証する単位」という考え方が、
// Customer 集約の境界そのものである。Customer と Address は必ず 1 つの
// トランザクションの中でまとめて読み書きされ（customer_repository.go の
// Save を参照）、他のどの集約もこの境界をまたいで Address を直接
// 更新することはできない。
type Customer struct {
	id        CustomerID
	name      string
	email     string
	addresses []Address
}

// NewCustomer は新規顧客を登録するコンストラクタである。
func NewCustomer(name, email string) (*Customer, error) {
	if name == "" {
		return nil, shared.NewDomainRuleError("customer: name must not be empty")
	}
	if utf8.RuneCountInString(name) > maxNameRunes {
		return nil, shared.NewDomainRuleError("customer: name must be at most %d characters, got %d", maxNameRunes, utf8.RuneCountInString(name))
	}
	if email == "" {
		return nil, shared.NewDomainRuleError("customer: email must not be empty")
	}
	// 「@ を含むか」という簡易チェックのみ行い、RFC 5322 準拠の厳密な
	// メールアドレス検証はあえて行わない。本サンプルの主眼は集約
	// （Customer + Address）の設計を示すことにあり、メールアドレス検証の
	// 厳密さそのものは主題ではない。実運用では正規表現による形式検証に
	// 加え、確認メールの送信・到達性チェックまで含めて検証するのが定石だが、
	// それらはこのサンプルの学習目的にとって過剰な複雑さになる。
	if !strings.Contains(email, "@") {
		return nil, shared.NewDomainRuleError("customer: email must contain @, got %q", email)
	}

	return &Customer{
		id:        GenerateCustomerID(),
		name:      name,
		email:     email,
		addresses: nil,
	}, nil
}

// ReconstructCustomer は永続化層から読み込んだデータをもとに Customer を
// 再構築する。
//
// NewCustomer との違いは order.ReconstructOrder / cart.ReconstructCart と
// 同じ理由による。「新規登録」と「DB からの復元」を同じコンストラクタに
// まとめてしまうと、意図の違いが型の上で区別できなくなる。
// リポジトリ実装（infrastructure 層）からのみ呼ばれることを想定している。
func ReconstructCustomer(id CustomerID, name, email string, addresses []Address) *Customer {
	addressesCopy := make([]Address, len(addresses))
	copy(addressesCopy, addresses)

	return &Customer{
		id:        id,
		name:      name,
		email:     email,
		addresses: addressesCopy,
	}
}

// ID は顧客の識別子を返す。
func (c *Customer) ID() CustomerID {
	return c.id
}

// Name は顧客の氏名を返す。
func (c *Customer) Name() string {
	return c.name
}

// Email は顧客のメールアドレスを返す。
func (c *Customer) Email() string {
	return c.email
}

// Addresses は顧客が登録した住所の一覧を返す。
//
// 呼び出し元が内部スライスを直接書き換えて集約の不変条件（デフォルトは
// 必ず1つ等）を壊せないよう、コピーを返す（cart.Cart.Items() /
// order.Order.Items() と同じカプセル化の理由による）。
func (c *Customer) Addresses() []Address {
	addresses := make([]Address, len(c.addresses))
	copy(addresses, c.addresses)
	return addresses
}

// DefaultAddress は現在のデフォルト住所を返す。
// 住所が 1 件も登録されていない顧客にはデフォルト住所が存在しないため、
// ドメインルール違反として扱う。
func (c *Customer) DefaultAddress() (Address, error) {
	for _, a := range c.addresses {
		if a.isDefault {
			return a, nil
		}
	}
	return Address{}, shared.NewDomainRuleError("customer: customer %s has no default address", c.id)
}

// AddAddress は顧客に新しい配送先住所を追加する。
//
// 郵便番号・都道府県・市区町村・番地等はすべて必須項目として扱う。
// また住所帳が無制限に増えるのを防ぐため、上限（maxAddresses）を超える
// 追加はドメインルール違反として拒否する。
//
// 【最初の 1 件は自動的にデフォルトになる】
// 「デフォルト住所は住所があれば必ず 1 つ」という不変条件は、Address と
// いう子エンティティ単体では守れない。もし「デフォルト住所を選ぶ操作」を
// 住所追加とは別の操作として利用者に強制すると、「住所を 1 件追加した
// 直後、まだデフォルトを選んでいない」という不変条件違反の時間帯が
// 生まれてしまう。そこで、集約ルートである Customer が「これが最初の
// 住所である」という事実を検知し、追加操作の中で自動的にデフォルトへ
// 昇格させる。これにより「住所を追加した瞬間、常にデフォルトが存在する」
// という不変条件を、利用者に追加の操作を要求することなく常に成立させ
// 続けられる。2 件目以降は既存のデフォルトをそのまま維持し、変更したい
// 場合は ChangeDefaultAddress を明示的に呼んでもらう。
func (c *Customer) AddAddress(postalCode, prefecture, city, line string) (AddressID, error) {
	if postalCode == "" {
		return "", shared.NewDomainRuleError("customer: postal code must not be empty")
	}
	if prefecture == "" {
		return "", shared.NewDomainRuleError("customer: prefecture must not be empty")
	}
	if city == "" {
		return "", shared.NewDomainRuleError("customer: city must not be empty")
	}
	if line == "" {
		return "", shared.NewDomainRuleError("customer: address line must not be empty")
	}
	if len(c.addresses) >= maxAddresses {
		return "", shared.NewDomainRuleError("customer: customer must not have more than %d addresses", maxAddresses)
	}

	isDefault := len(c.addresses) == 0

	addr := Address{
		id:         GenerateAddressID(),
		postalCode: postalCode,
		prefecture: prefecture,
		city:       city,
		line:       line,
		isDefault:  isDefault,
	}
	c.addresses = append(c.addresses, addr)
	return addr.id, nil
}

// ChangeDefaultAddress は指定した住所を新しいデフォルト住所に変更する。
// 対象の住所が顧客に登録されていない場合はドメインルール違反として扱う。
//
// 「既存のデフォルトを降ろす」と「新しいデフォルトを立てる」という
// 2 つの書き換えを、この 1 つのメソッドの中でアトミックに行う点に注意する。
// もしこの 2 ステップを外部（アプリケーション層）に分割して行わせて
// しまうと、その途中の瞬間には「デフォルトが 0 個」または「デフォルトが
// 2 個」という不変条件違反の状態が一時的にせよ発生しうる。集約ルートの
// メソッド 1 回の呼び出し（＝1 回のドメイン操作）の中で完結させることで、
// 呼び出し前後のどちらの時点で見ても不変条件が破れていないことを保証する。
func (c *Customer) ChangeDefaultAddress(id AddressID) error {
	found := -1
	for i, a := range c.addresses {
		if a.id == id {
			found = i
			break
		}
	}
	if found == -1 {
		return shared.NewDomainRuleError("customer: address %s is not registered for customer %s", id, c.id)
	}

	for i := range c.addresses {
		c.addresses[i].isDefault = false
	}
	c.addresses[found].isDefault = true
	return nil
}

// RemoveAddress は指定した住所を顧客の住所帳から削除する。
// 対象の住所が顧客に登録されていない場合はドメインルール違反として扱う。
//
// 【設計判断: デフォルト住所は「他に住所が残る場合」削除できない】
// 仮にこの制約を設けず、デフォルト住所の削除を無条件に許してしまうと、
// 残った住所に対してどちらか一方を選ばざるを得ない。
//   - デフォルトが存在しない状態を許容する
//     → 「住所が 1 件以上あればデフォルトは必ず 1 つ」という不変条件が破れる。
//   - システムが残った住所のどれかを勝手に新しいデフォルトへ昇格させる
//     → どの住所を新しいデフォルトにするかは顧客本人の意思決定であり、
//     システムが暗黙に肩代わりしてよい判断ではない（例えば配送先を
//     間違えるリスクにつながる）。
//
// そのため、他に住所が残る場合はまず ChangeDefaultAddress で明示的に
// 新しいデフォルトを選んでもらうことを利用者に強制する設計とした。
//
// 一方、削除後に住所が 0 件になる（＝最後の 1 件を削除する）場合は、
// 「住所が 1 件以上あればデフォルトは必ず 1 つ」という不変条件の前提
// （1 件以上ある）自体が成立しなくなるだけであり、不変条件は破れない。
// そのためこのケースは削除を許可する。
func (c *Customer) RemoveAddress(id AddressID) error {
	found := -1
	for i, a := range c.addresses {
		if a.id == id {
			found = i
			break
		}
	}
	if found == -1 {
		return shared.NewDomainRuleError("customer: address %s is not registered for customer %s", id, c.id)
	}

	if c.addresses[found].isDefault && len(c.addresses) > 1 {
		return shared.NewDomainRuleError("customer: cannot remove default address %s while other addresses remain; change the default address first", id)
	}

	c.addresses = append(c.addresses[:found], c.addresses[found+1:]...)
	return nil
}
