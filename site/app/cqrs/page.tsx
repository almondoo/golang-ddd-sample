import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { CqrsDiagram } from "@/components/diagrams/cqrs-diagram";

export const metadata: Metadata = {
  title: "CQRS(軽量版)とは",
  description:
    "書く道と読む道を分ける。CQRSの考え方をGo DDDサンプルのProductQueryServiceで図解します。",
};

export default function CqrsPage() {
  return (
    <ConceptLayout
      slug="cqrs"
      eyebrow="14 / CQRS(軽量版)"
      title="書く道と、読む道を分ける"
      lead="書き込みはドメインを通る。読み取りは、ドメインを素通りしてDBから直接組み立てる。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/application/usecase/catalog/list_products.go"}
            </span>
            {"\n"}
            {"func (uc *ListProductsUseCase) Execute(ctx context.Context) (*Output, error) {\n"}
            {"    return uc.queryService.List(ctx) "}
            <span className="text-background/50">{"// これだけ"}</span>
            {"\n}"}
          </>
        ),
        note: (
          <>
            <code className="font-mono">ProductQueryService</code>{" "}
            の実装(GORM)は<code className="font-mono">db.Find(&#123;&#125;)</code>
            で取得した行を<code className="font-mono">toProductDTO</code>{" "}
            でそのままDTOに変換する。Product集約は一切組み立てない。
          </>
        ),
      }}
      summary="commandとqueryで別データストアに分ける本格CQRSではない。「ドメインを通るかどうか」だけを分ける軽量版。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <CqrsDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        一覧表示のような「ただ表示するだけ」の処理でドメイン集約を都度組み立てるのはオーバーヘッド、というのがクエリサービスを分けている理由。
      </p>
    </ConceptLayout>
  );
}
