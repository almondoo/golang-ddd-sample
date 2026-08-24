// Package controller は presentation 層のコントローラを提供する。
// MVC の controller に相当し、HTTP リクエストを usecase の入出力へ
// 変換することだけを責務とする（ビジネスルールは持たない）。
// パッケージ名を http にすると標準ライブラリ net/http と衝突するため controller としている。
package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// errorResponse はエラー時に返す共通の JSON ボディ。
type errorResponse struct {
	Error string `json:"error"`
}

// WriteJSON は任意のペイロードを JSON でレスポンスに書き出す。
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// ヘッダ送信後はステータスを変えられないため、ログに残すことしかできない
		slog.Error("failed to encode response", "error", err)
	}
}

// WriteError はドメイン層・アプリケーション層のエラーを HTTP ステータスへ変換して返す。
//
// 変換規則:
//   - shared.ErrNotFound(をラップしたもの) → 404 Not Found
//   - shared.NewDomainRuleError で作られたビジネスルール違反 → 422 Unprocessable Entity
//   - それ以外(想定外の失敗) → 500 Internal Server Error
//
// ドメイン層は HTTP を知らないため、この変換は presentation 層の責務になる。
func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case shared.IsDomainRuleError(err):
		WriteJSON(w, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
	default:
		// 内部エラーの詳細はクライアントに漏らさず、ログにのみ残す
		slog.Error("internal server error", "error", err)
		WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

// WriteBadRequest はリクエスト形式の不備(JSON パース失敗など)を 400 で返す。
// ビジネスルール違反(422)とは区別する: 400 は「読めない」、422 は「読めたがルールに反する」。
func WriteBadRequest(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusBadRequest, errorResponse{Error: msg})
}

// DecodeJSON はリクエストボディを dst にデコードする。
// 未知のフィールドは拒否し、クライアントのタイプミスに早く気づけるようにする。
func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
