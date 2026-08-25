package shared

import (
	"errors"
	"fmt"
)

// ErrNotFound はリソースが見つからなかったことを表すセンチネルエラーである。
//
// 各コンテキストのリポジトリ実装は、対象が見つからない場合にこの値を
// そのまま返すのではなく、%w でラップして返すこと（例:
// fmt.Errorf("cart not found: id=%s: %w", id, shared.ErrNotFound)）。
// こうすることで呼び出し側は errors.Is(err, shared.ErrNotFound) で
// 判定できる一方、エラーメッセージには文脈（どの ID か等）を残せる。
//
// プレゼンテーション層はこのエラーを HTTP 404 に変換する。
var ErrNotFound = errors.New("not found")

// ErrConflict は楽観ロック（optimistic lock）による更新競合を表す
// センチネルエラーである。
//
// 「読み込んだ時点のバージョンと、書き込もうとする時点のバージョンが
// 一致しない」状態、つまり自分が読んでから保存するまでの間に別の
// トランザクションが同じ集約を先に更新してしまった（lost update の
// 発生条件）ことを表す。ErrNotFound 同様、リポジトリ実装は %w でラップして
// 返すこと（例: fmt.Errorf("stock update conflict: id=%s: %w", id,
// shared.ErrConflict)）。呼び出し側は errors.Is(err, shared.ErrConflict) で
// 判定し、リトライ（再読み込みしてから再実行）するかどうかを決められる。
//
// プレゼンテーション層はこのエラーを HTTP 409 Conflict に変換する。
var ErrConflict = errors.New("conflict")

// ruleError はドメインルール（不変条件・ビジネスルール）違反を表す
// エラー型である。amount が負である、在庫が不足している、といった
// 「入力値としては妥当だが業務的に許されない」状態を表現するために使う。
//
// 型を非公開にしているのは、生成経路を NewDomainRuleError に限定し、
// 呼び出し側には errors.As 経由でのみ種別を判定させるためである。
// これにより「ドメインルールエラーである」という事実だけを型で保証し、
// 具体的なメッセージ内容には依存させない設計にしている。
type ruleError struct {
	message string
}

// Error は error インターフェースの実装である。
func (e *ruleError) Error() string {
	return e.message
}

// NewDomainRuleError はドメインルール違反を表すエラーを生成する。
// fmt.Errorf と同様に format 文字列で任意のメッセージを組み立てられる。
func NewDomainRuleError(format string, args ...any) error {
	return &ruleError{message: fmt.Sprintf(format, args...)}
}

// IsDomainRuleError は err がドメインルール違反エラー（またはそれをラップした
// エラー）かどうかを判定する。
//
// プレゼンテーション層はこの関数を使い、ErrNotFound → 404、
// DomainRuleError → 422（Unprocessable Entity）、それ以外 → 500 という
// マッピングを行う。この判定をプレゼンテーション層に閉じ込めることで、
// ドメイン層は HTTP という配送方式を一切知らずに済む。
func IsDomainRuleError(err error) bool {
	var re *ruleError
	return errors.As(err, &re)
}
