// Package main は本アプリケーションの composition root（コンポジション・
// ルート）である。
//
// 【composition root とは何か、なぜここに全部集めるのか】
// これまでの各層（domain / application / infrastructure / presentation）は、
// 依存性のルール（内側の層は外側の層を知らない）に従って作られてきた。
// たとえば application 層は「catalog.Repository」というインターフェースには
// 依存するが、それを実装する persistence.ProductRepository（GORM 実装）の
// 存在は知らない。しかし実行時には、どこかで「インターフェース ← 具体的な
// 実装」という配線を実際に行い、依存性を注入（Dependency Injection）しなければ
// プログラムは動かない。その「唯一、全レイヤーの存在を知っていてよい場所」が
// この main.go であり、これを composition root と呼ぶ。
//
// 内側の層が外側の層を知らずに済んでいるのは、この main（正確には main が
// 呼び出す initializeServer）が代わりに「どの実装をどのインターフェースに
// 差し込むか」を一手に引き受けているからである。逆に言えば、main 以外の
// どこかで new persistence.XxxRepository(...) のような「具体的な実装を
// 知ったうえでの組み立て」が現れたら、それは依存性のルール違反の兆候である。
//
// 【wire の導入について】
// 依存の組み立て自体は github.com/almondoo/wire（コンパイル時 DI コード
// 生成ツール）が生成する initializeServer（cmd/api/wire_gen.go）に委ねて
// いる。main はもはや「設定を読み込み、初期化された mux でサーバーを
// 起動するだけ」の薄い層になり、composition root としての配線の詳細は
// cmd/api/wire.go（設計図）と cmd/api/providers.go（手組みの provider）に
// 移った。
package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// 1. 設定値を環境変数から読み込む。
	//
	// 学習用サンプルなので専用の config パッケージや viper 等は導入せず、
	// os.Getenv + フォールバック値というもっとも単純な形にとどめている。
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=ddd_sample port=5432 sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. 依存関係の組み立ては wire が生成した initializeServer に委ねる。
	//
	// DB 接続の確立・スキーマの用意・各レイヤーのコンポーネント生成・
	// HTTP コントローラの mux 登録まで、すべてこの 1 呼び出しの中で
	// 行われる（本サンプルにイベントバスは存在せず、コンテキストを
	// またぐ処理はアプリケーション層からの直接呼び出しで行っている。
	// 詳細は internal/application/usecase/order/place_order.go の
	// コメントを参照）。実体は cmd/api/wire_gen.go にあり、その設計図は
	// cmd/api/wire.go（wireinject タグ付き、通常ビルドからは除外される）
	// に書かれている。
	mux, err := initializeServer(dsn)
	if err != nil {
		slog.Error("failed to initialize server", "error", err)
		os.Exit(1)
	}

	// 3. HTTP サーバーを起動する。
	//
	// 学習用サンプルとして単純さを優先し、http.ListenAndServe をそのまま
	// 使っている。production サーバーでは、シグナル受信によるグレースフル
	// シャットダウン（http.Server.Shutdown）や ReadTimeout/WriteTimeout 等の
	// タイムアウト設定を追加するのが一般的である。
	addr := ":" + port
	slog.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
