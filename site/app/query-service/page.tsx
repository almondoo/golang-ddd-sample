import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { QueryServiceDiagram } from "@/components/diagrams/query-service-diagram";

export const metadata: Metadata = {
  title: "クエリサービスとは",
  description:
    "見るだけなら、正式ルートはいらない。CQRSの読み取り側の port/adapter をGo DDDサンプルのCartQueryServiceで図解します。",
};

export default function QueryServicePage() {
  return (
    <ConceptLayout
      slug="query-service"
      eyebrow="15 / クエリサービス"
      title="見るだけなら、正式ルートはいらない"
      lead="書くときは集約とリポジトリの正式ルートを通る。見るだけの窓口は、ドメインを素通りしてDTOを直接返す。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/infrastructure/persistence/cart_query_service.go\n"}
              {"// cart_items(cartコンテキスト)とproducts(catalogコンテキスト)を\n"}
              {"// SQLレベルで直接JOINする"}
            </span>
            {"\n\n"}
            {"func (q *CartQuery) FindByCustomerID(ctx context.Context, customerID string) (*cartusecase.CartDTO, error) {\n"}
            {"    …\n"}
            {'    db.Table("cart_items").\n'}
            {'        Joins("JOIN products ON products.id = cart_items.product_id").\n'}
            {"        …\n"}
            {"}"}
          </>
        ),
        note: (
          <>
            <code className="font-mono">CartQueryService</code>{" "}
            はapplication層が宣言するport、実装はinfrastructure層の
            <code className="font-mono">CartQuery</code>。catalogの物理スキーマ(テーブル名・カラム名)への結合が読み取り側に残るため、catalog側でスキーマを変更すると実行時にこのクエリが壊れうる、という正直な代償も伴う。
          </>
        ),
      }}
      summary="見るだけなら、正式ルートを通らなくていい。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <QueryServiceDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        CQRSが「書く道と読む道を分ける」という考え方だとしたら、クエリサービスはその読む道を実際に実装するport/adapterの仕組み。
        <code className="font-mono">ProductQueryService</code> /{" "}
        <code className="font-mono">CartQueryService</code>{" "}
        がapplication層のport、GORM実装のQuery型がinfrastructure層のadapterにあたる。
      </p>
    </ConceptLayout>
  );
}
