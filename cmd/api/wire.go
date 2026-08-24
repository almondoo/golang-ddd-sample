//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/almondoo/wire"
)

// このファイルは DI の「設計図」であり、wireinject ビルドタグ付きのため
// 通常の go build / go test には一切含まれない（コンパイル対象は
// wire_gen.go の方）。wire はリフレクションではなくコード生成で DI を
// 行うツールであるため、ここに書いた配線に過不足や型の不一致があれば
// `go run github.com/almondoo/wire/cmd/wire ./cmd/api` の実行時点で
// エラーとして検出できる（実行時に初めて panic する手書き DI や
// リフレクションベースの DI コンテナとの大きな違いである）。
//
// initializeServer は DSN を受け取り、DB 接続からユースケース・
// コントローラ・mux までの依存グラフ全体を組み立てて *http.ServeMux を
// 返す injector（注入器）である。中身は wire.Build の呼び出しだけで、
// 実際の組み立てコードは wire が providers.go の各 provider を解析して
// 自動生成する。
func initializeServer(dsn string) (*http.ServeMux, error) {
	wire.Build(
		infrastructureSet,
		catalogSet,
		cartSet,
		orderSet,
		customerSet,
		inventorySet,
		shippingSet,
		couponSet,
		controllerSet,
	)
	return nil, nil
}
