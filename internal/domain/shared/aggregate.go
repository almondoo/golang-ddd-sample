package shared

// AggregateBase は集約（Aggregate）が共通で持つ「発生したドメインイベントを
// 一時的に保持する」機能を提供する基底構造体である。
//
// 各コンテキストの集約ルート（例: Cart, Order）はこの構造体を埋め込む
// （embedding）ことで、Record / PullEvents の機能を無償で獲得できる。
//
// なぜ集約がイベントを直接パブリッシュしないのか？
// 集約はあくまでドメインルールを守ることだけに責任を持ち、
// 「いつ・どうやって」イベントを外部に伝えるかは知るべきではない
// （単一責任の原則）。また、イベントは「永続化が成功して初めて確定した
// 事実」になるため、集約の中で即座に配信してしまうと、後続の DB 保存が
// 失敗した場合に「起きていないこと」を配信してしまう不整合が生じる。
//
// そこでアプリケーション層のユースケースが
//  1. 集約を操作する（結果として Record が呼ばれ、イベントが集約内に溜まる）
//  2. リポジトリで永続化する
//  3. 永続化が成功したら PullEvents でイベントを取り出し、EventPublisher で配信する
//
// という順序を守ることで、「実際に起きて確定したこと」だけを配信できる。
type AggregateBase struct {
	events []DomainEvent
}

// Record は集約内で発生したドメインイベントを記録する。
// 集約のドメインロジック（例: カートに商品を追加するメソッド）の内部から呼ぶ。
func (a *AggregateBase) Record(e DomainEvent) {
	a.events = append(a.events, e)
}

// PullEvents は記録済みのイベントをすべて取り出し、内部の記録をクリアする。
// アプリケーション層が永続化成功後に一度だけ呼び出すことを想定している。
// 呼び出すたびに空になるため、二重配信を防ぐことができる。
func (a *AggregateBase) PullEvents() []DomainEvent {
	events := a.events
	a.events = nil
	return events
}
