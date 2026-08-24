# ユビキタス言語(Ubiquitous Language)

ユビキタス言語とは、ドメインエキスパートと開発者が同じ言葉を同じ意味で使う、コードとコミュニケーションの両方に共通する語彙です。コード中の型名・メソッド名・エラーメッセージが、業務で実際に使われる用語と一致していることを目指します。Microsoftのdomain-analysisは、ユビキタス言語(コンテキストごとの共有語彙)を戦略的設計の「中心」と位置づけています(詳細は[ddd-research.md](../ddd-research.md))。

## なぜ必要か

コードの用語と業務会話の用語がずれていると、その都度「コード上のこの言葉は業務的には何を指すか」という翻訳が必要になり、認識のずれや誤りが生まれやすくなります。型名・メソッド名を業務用語にそのまま合わせることで、コードを読むこと自体がドメイン知識の学習になります。

## このリポジトリでの現状

型名・メソッド名のレベルでは、`Customer.ChangeDefaultAddress`や`Order.Ship`のように業務用語がそのまま使われており、一貫性があります。

**ID型の対照表は[context-map.md](../context-map.md)に整備されています**([specs/ddd-improvements.md](../specs/ddd-improvements.md)項目9対応)。`ProductID`が3コンテキスト、`CustomerID`が3コンテキスト、`OrderID`が2コンテキストに重複定義されており([bounded-context.md](bounded-context.md)参照)、これらは同一のUUID空間を共有します。実際に重複定義されているID型は次の通りです(コードを走査して確認)。

| ID型 | 定義箇所(コンテキスト) |
|---|---|
| `ProductID` | [catalog](../../internal/domain/catalog/product_id.go) / [cart](../../internal/domain/cart/ids.go) / [inventory](../../internal/domain/inventory/product_id.go) |
| `CustomerID` | [order](../../internal/domain/order/ids.go) / [cart](../../internal/domain/cart/ids.go) / [customer](../../internal/domain/customer/ids.go) |
| `OrderID` | [order](../../internal/domain/order/ids.go) / [shipping](../../internal/domain/shipping/ids.go) |

いずれも値としては同じUUID文字列空間を共有しますが、型としては別物であり、コンテキストをまたいで直接代入するとコンパイルエラーになります(意図的な設計、詳細は[bounded-context.md](bounded-context.md))。

**エラーメッセージは英語に統一済みです。** かつて`order`のusecase層の3ファイル(`place_order.go`・`ship_order.go`・`cancel_order.go`)だけが日本語のメッセージを使っており、同じ422エラー面に2言語が混在していましたが、[specs/ddd-improvements.md](../specs/ddd-improvements.md)項目3として英語に統一しました。現在は`internal`配下の`NewDomainRuleError`メッセージはすべて英語です(例: [cart.go](../../internal/domain/cart/cart.go)の`"cart: quantity must be between..."`)。

## 注意点・よくある誤解

- ユビキタス言語は「日本語で統一する」「英語で統一する」という言語選択そのものではなく、「一貫した用語を使う」ことが本質です。本リポジトリはコメントを日本語、識別子・エラーメッセージを英語に統一する方針です(過去に`order`のusecase層3ファイルがこの方針から外れていましたが、修正済みです)。
- 型名の一貫性(例: `Money`、`DomainRuleError`)があるからといって、エラーメッセージという「実行時に人が読むテキスト」まで一貫しているとは限りません。両方を別々に確認する必要があります。
