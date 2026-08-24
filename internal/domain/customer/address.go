package customer

// Address は顧客の配送先住所を表す、Customer 集約に属する子エンティティ
// （集約内部のエンティティ）である。
//
// 【集約内部のエンティティという位置づけ】
// Address は CartItem（cart コンテキスト）や OrderItem（order コンテキスト）
// と同じく、集約ルートの外からは直接生成・変更させないという方針を採る。
// フィールドをすべて非公開にし、公開コンストラクタ（NewAddress のような
// 関数）を用意していないのはそのためである。Address を新しく作る唯一の
// 経路は Customer.AddAddress であり、状態を書き換える唯一の経路は
// Customer.ChangeDefaultAddress / Customer.RemoveAddress である。
//
// なぜここまで厳密に「集約ルート経由」を強制するのか: Customer 集約は
// 「デフォルト住所は、住所が 1 件でもあれば必ずちょうど 1 つ存在する」
// という不変条件を持つ。この条件は Address 1 件だけを見て判断できるもの
// ではなく、兄弟にあたる他の Address すべての isDefault の状態を横断的に
// 見なければ判定・維持できない。もし Address 自身が isDefault を自由に
// 書き換えられるメソッドを公開してしまうと、「この住所をデフォルトにしたが、
// 元のデフォルト住所のフラグは誰も倒さなかった」といった不整合が
// 容易に発生してしまう。集約ルートである Customer だけが全住所を
// 見渡せる立場にあるため、整合性の維持を一手に引き受けている。
type Address struct {
	id         AddressID
	postalCode string
	prefecture string
	city       string
	line       string
	isDefault  bool
}

// ID はこの住所の識別子を返す。
func (a Address) ID() AddressID {
	return a.id
}

// PostalCode は郵便番号を返す。
func (a Address) PostalCode() string {
	return a.postalCode
}

// Prefecture は都道府県を返す。
func (a Address) Prefecture() string {
	return a.prefecture
}

// City は市区町村を返す。
func (a Address) City() string {
	return a.city
}

// Line は番地・建物名など住所の残りの部分を返す。
func (a Address) Line() string {
	return a.line
}

// IsDefault はこの住所が顧客のデフォルト住所かどうかを返す。
func (a Address) IsDefault() bool {
	return a.isDefault
}

// ReconstructAddress は永続化層から読み込んだデータをもとに Address を
// 再構築する。
//
// 通常 Address は Customer.AddAddress を経由してのみ生成され、その過程で
// 各フィールドの非空制約や住所件数の上限といった不変条件がチェックされる。
// しかし DB から読み込む際は「すでに過去に検証済みのデータをそのまま
// 復元する」だけなので、再度の検証は行わない（cart.ReconstructCartItem /
// order.ReconstructOrderItem と同じ設計上の理由による）。
// リポジトリ実装（infrastructure 層）からのみ呼ばれることを想定している。
func ReconstructAddress(id AddressID, postalCode, prefecture, city, line string, isDefault bool) Address {
	return Address{
		id:         id,
		postalCode: postalCode,
		prefecture: prefecture,
		city:       city,
		line:       line,
		isDefault:  isDefault,
	}
}
