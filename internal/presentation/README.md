# presentation 層

## この層の役割

`internal/presentation` は、外部（Web ブラウザ、他サービス等）とのやり取りを
担う層である。本サンプルでは主に HTTP コントローラをここに置く想定である。

この層の責務は「外部プロトコル（HTTP）とアプリケーション層のユースケースを
橋渡しすること」に限定される。ビジネスルールの判断は一切行わない。

## この層で行うこと

1. **リクエストのパースとバリデーション**: HTTP リクエストボディ（JSON 等）を
   Go の構造体にデコードし、必須項目の有無のような形式的なバリデーションを行う
  （ビジネスルールとしてのバリデーション、例えば「在庫数を超える数量は
   注文できない」はドメイン層の責務であり、ここでは行わない）。
2. **DTO ↔ コマンド/クエリの変換**: HTTP リクエストの DTO を、
   アプリケーション層が受け取るコマンド（またはクエリ）に変換してユースケースを呼ぶ。
3. **ユースケースの実行結果 ↔ レスポンス DTO の変換**: ユースケースが返した
   結果（またはクエリの返す DTO）を、HTTP レスポンス用の JSON に変換する。
4. **エラー → HTTP ステータスコードへのマッピング**: 下記参照。

## エラー → HTTP ステータスコードのマッピング

ドメイン層・アプリケーション層は HTTP を一切知らないため、
「このエラーはどの HTTP ステータスに対応するか」という判断は
プレゼンテーション層に閉じ込める。本サンプルでは以下の方針を採る。

| エラーの種類 | 判定方法 | HTTP ステータス |
| --- | --- | --- |
| リソースが見つからない | `errors.Is(err, shared.ErrNotFound)` | 404 Not Found |
| ドメインルール違反 | `shared.IsDomainRuleError(err)` | 422 Unprocessable Entity |
| その他（予期しないエラー） | 上記いずれにも該当しない | 500 Internal Server Error |

擬似コードで表すと以下のようになる。

```go
func handleError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, shared.ErrNotFound):
        writeJSON(w, http.StatusNotFound, errorBody(err))
    case shared.IsDomainRuleError(err):
        writeJSON(w, http.StatusUnprocessableEntity, errorBody(err))
    default:
        // 予期しないエラーの詳細をそのまま外部に返すと内部実装が漏れる
        // 恐れがあるため、ログには詳細を出しつつレスポンスは汎用的な
        // メッセージに留める。
        log.Printf("unexpected error: %v", err)
        writeJSON(w, http.StatusInternalServerError, genericErrorBody())
    }
}
```

この方針により、ドメイン層でエラーの種類さえ正しく表現しておけば
（センチネルエラーのラップ、`NewDomainRuleError` の利用）、
プレゼンテーション層は個々のユースケースの内容を知らなくても
一貫したエラーハンドリングができる。
