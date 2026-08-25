import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { NotAdoptedBadge } from "@/components/not-adopted-badge";
import { DomainServiceDiagram } from "@/components/diagrams/domain-service-diagram";

export const metadata: Metadata = {
  title: "ドメインサービスとは",
  description:
    "集約をまたぐ計算の置き場所。ドメインサービス(Domain Service)の考え方と、本リポジトリで採用していない理由を図解します。",
};

export default function DomainServicePage() {
  return (
    <ConceptLayout
      slug="domain-service"
      eyebrow="09 / ドメインサービス"
      title="集約をまたぐ計算の、置き場所"
      lead="1つの集約だけでは判断できないロジックを置く、第三の場所。ただし本リポジトリでは出番がなかった。"
      badge={<NotAdoptedBadge />}
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// domain-service.md より\n"}
              {"// 各コンテキストは集約を1つずつしか持たない"}
            </span>
            {"\n"}
            {"catalog  → Product\n"}
            {"cart     → Cart\n"}
            {"order    → Order\n"}
            {"customer → Customer\n"}
            {"...\n\n"}
            <span className="text-background/50">
              {"// 同一コンテキスト内で複数集約をまたぐ場面が構造的に無い"}
            </span>
          </>
        ),
        note: (
          <>
            コンテキストをまたぐ調整(顧客確認・在庫引当・クーポン適用など)は実在するが、それは
            <code className="font-mono">PlaceOrderUseCase</code>
            (アプリケーションサービス)の仕事。「同一コンテキスト内の複数集約をまたぐ計算」とは性質が違う。
          </>
        ),
      }}
      summary="ドメインサービスの不在は実装漏れではなく、必要になる場面が現状のコードに無いという監査結果。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <DomainServiceDiagram />
        </figure>
      </Reveal>

      <div className="mx-auto mt-8 max-w-xl rounded-3xl border-2 border-border bg-card px-6 py-5">
        <h3 className="mb-2 text-lg font-extrabold">なぜ不採用?</h3>
        <p className="text-sm text-muted-foreground">
          本リポジトリの各コンテキスト(catalog / cart / order / customer / inventory / shipping /
          coupon)はそれぞれ集約を1つずつしか持たないため、「同一コンテキスト内で複数集約をまたぐ」ロジックが構造的に発生しない。無理に人工的な例を作るのは学習コスト優先の方針に反するため、「なぜ無いか」を正直に書く方針にした。将来2つ目の集約が増えて両者をまたぐ計算が必要になったら、そのときがドメインサービス導入の合図になる。
        </p>
      </div>
    </ConceptLayout>
  );
}
