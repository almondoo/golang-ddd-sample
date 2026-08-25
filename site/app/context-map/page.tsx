import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { ContextMapDiagram } from "@/components/diagrams/context-map-diagram";

export const metadata: Metadata = {
  title: "コンテキストマップとは",
  description:
    "売り場同士のつきあい方を1枚の地図にする。コンテキストマップの考え方をGo DDDサンプルのPlaceOrderUseCaseで図解します。",
};

export default function ContextMapPage() {
  return (
    <ConceptLayout
      slug="context-map"
      eyebrow="04 / コンテキストマップ"
      title="つきあい方を、1枚の地図にする"
      lead="誰が誰を呼ぶか。どこが結合しているか。それを1枚にまとめておけば、結合は怖くない。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/application/usecase/order/place_order.go\n"}
              {"// PlaceOrderUseCaseは5コンテキストのドメインパッケージをimportする"}
            </span>
            {"\n\n"}
            {"import (\n"}
            {'    domaincart "…/internal/domain/cart"\n'}
            {'    domaincatalog "…/internal/domain/catalog"\n'}
            {'    domaincoupon "…/internal/domain/coupon"\n'}
            {'    domaincustomer "…/internal/domain/customer"\n'}
            {'    domaininventory "…/internal/domain/inventory"\n'}
            {"    …\n)"}
          </>
        ),
        note: (
          <>
            これは<code className="font-mono">Customer-Supplier</code>や
            <code className="font-mono">ACL</code>・
            <code className="font-mono">OHS</code>のような標準パターンには当てはまらない。単一チームが単一デプロイ単位として書く「モジュラーモノリスのapplication層統合」であると、docs/context-map.mdは正直に評価している。
          </>
        ),
      }}
      summary="地図にしておけば、結合は怖くない。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <ContextMapDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        境界づけられたコンテキストが「部屋」だとしたら、コンテキストマップは「誰の部屋が誰とどうつながっているか」の見取り図。orderのapplication層がcustomer/cart/catalog/inventory/couponを直接オーケストレーションし、cartのクエリサービスはcatalogのテーブルへ実行時に結合したSQLを発行する。
      </p>
    </ConceptLayout>
  );
}
