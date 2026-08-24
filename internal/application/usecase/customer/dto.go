// Package customer は customer コンテキストのユースケース（アプリケーション
// サービス）をまとめたパッケージである。
//
// 1 ファイル = 1 ユースケースの原則に従い、ファイルを以下のように分けている。
//   - register_customer.go / add_address.go / change_default_address.go /
//     remove_address.go … 書き込み系（コマンド）のユースケース。
//     domaincustomer.Repository と tx.Manager に依存し、Customer 集約を
//     読み込み・操作し・保存する。
//   - get_customer.go             … 読み取り系（クエリ）のユースケース。
//     CustomerQueryService にのみ依存し、ドメイン層を経由せず DTO を
//     直接返す。
//
// cart / order コンテキストと同じ方針で、コマンドとクエリの区別は
// パッケージではなくファイル名と依存する型（Repository+tx.Manager か、
// QueryService のみか）で表現する。
package customer

// CustomerDTO は顧客を画面表示用に平坦化したデータ転送オブジェクトである。
//
// ドメインの Customer 集約とは別に DTO を用意しているのは、クエリ側
// （読み取り）とコマンド側（書き込み）で最適な形が異なるためである
// （CQRS の考え方）。Addresses は Customer.Addresses() の返す
// []customer.Address（ドメインオブジェクト）ではなく、JSON へそのまま
// シリアライズできる平坦なフィールドの集まりとして表現する。
type CustomerDTO struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Addresses []AddressDTO `json:"addresses"`
}

// AddressDTO は顧客の住所（子エンティティ）1 件を画面表示用に表したものである。
type AddressDTO struct {
	ID         string `json:"id"`
	PostalCode string `json:"postalCode"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Line       string `json:"line"`
	IsDefault  bool   `json:"isDefault"`
}
