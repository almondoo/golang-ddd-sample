import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { BoundedContextDiagram } from "@/components/diagrams/bounded-context-diagram";

export const metadata: Metadata = {
  title: "境界づけられたコンテキストとは",
  description:
    "大きなお店を意味のまとまりで小部屋に分ける、境界づけられたコンテキストをGo DDDサンプルの7つのコンテキストで図解します。",
};

export default function BoundedContextPage() {
  return (
    <ConceptLayout
      slug="bounded-context"
      eyebrow="03 / 境界づけられたコンテキスト"
      title="大きなお店を、小部屋に分ける"
      lead="意味のまとまりごとに部屋を分けると、迷子にならない。このリポジトリには7つの部屋(コンテキスト)がある。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// 同じ「商品ID」でも、部屋(パッケージ)ごとに型を定義しなおす\n"}
              {"// internal/domain/catalog/product_id.go"}
            </span>
            {"\ntype ProductID string\n\n"}
            <span className="text-background/50">
              {"// internal/domain/cart/ids.go"}
            </span>
            {"\ntype ProductID string\n\n"}
            <span className="text-background/50">
              {"// internal/domain/inventory/product_id.go"}
            </span>
            {"\ntype ProductID string"}
          </>
        ),
        note: "値としては同じUUID文字列でも、型としては別物。部屋をまたいでうっかり代入するとコンパイルエラーになる。",
      }}
      summary="部屋(コンテキスト)ごとに言葉と型を分けると、大きなお店でも迷子にならない。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <BoundedContextDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        「商品ID」は3つの部屋にあるけど、部屋ごとに意味も型もちがう。だから混同しない(コンパイルの時点で防げる)。
      </p>
    </ConceptLayout>
  );
}
