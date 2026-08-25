import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { NotAdoptedBadge } from "@/components/not-adopted-badge";
import { SpecificationDiagram } from "@/components/diagrams/specification-diagram";

export const metadata: Metadata = {
  title: "仕様(Specification)とは",
  description:
    "条件そのものを部品にする。Specificationパターンの考え方と、本リポジトリで採用していない理由を図解します。",
};

export default function SpecificationPage() {
  return (
    <ConceptLayout
      slug="specification"
      eyebrow="18 / 仕様(Specification)"
      title="条件そのものを、部品にする"
      lead="「満たすかどうか」の判定ロジックを、1つのオブジェクトとして持ち運べるようにする。"
      badge={<NotAdoptedBadge />}
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// 教科書的なイメージ(本リポジトリに実装はない)"}
            </span>
            {"\n"}
            {"type InStockSpec struct{}\n"}
            {"func (s InStockSpec) IsSatisfiedBy(p Product) bool {\n"}
            {"    return p.Stock() > 0\n"}
            {"}\n\n"}
            <span className="text-background/50">{"// AND/OR/NOTで組み合わせられるのが特徴"}</span>
            {"\n"}
            {"spec := InStockSpec{}.And(InCategorySpec{Category: \"book\"})"}
          </>
        ),
        note: (
          <>
            実際のクエリは<code className="font-mono">ProductQuery.List</code>{" "}
            /<code className="font-mono">FindByID</code>{" "}
            のような単一キーの問い合わせのみで、複数条件を組み合わせる検索は存在しない。
          </>
        ),
      }}
      summary="Specification不採用は実装漏れではなく、条件を合成・再利用する必要が現状のコードに無いという監査結果。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <SpecificationDiagram />
        </figure>
      </Reveal>

      <div className="mx-auto mt-8 max-w-xl rounded-3xl border-2 border-border bg-card px-6 py-5">
        <h3 className="mb-2 text-lg font-extrabold">なぜ不採用?</h3>
        <p className="text-sm text-muted-foreground">
          理由は2つ。(1)フィルタ・検索条件は軽量CQRSのクエリサービスに閉じており、複数条件を組み合わせる検索がそもそも実装されていない。(2)各コンテキストが単一の集約しか持たず、不変条件も集約のメソッド内で完結する単純なものばかりで、判定ロジックを使い回す動機がない。無理に人工的な例を作るより、「なぜ無いか」を正直に書く方針。
        </p>
      </div>
    </ConceptLayout>
  );
}
