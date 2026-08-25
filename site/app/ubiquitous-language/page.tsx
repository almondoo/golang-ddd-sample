import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { UbiquitousLanguageDiagram } from "@/components/diagrams/ubiquitous-language-diagram";

export const metadata: Metadata = {
  title: "ユビキタス言語とは",
  description:
    "店員さんも開発者もコードも、同じ言葉で話す。ユビキタス言語をGo DDDサンプルのOrder型で図解します。",
};

export default function UbiquitousLanguagePage() {
  return (
    <ConceptLayout
      slug="ubiquitous-language"
      eyebrow="02 / ユビキタス言語"
      title="みんな、同じ言葉で話す"
      lead="店員さんも、開発者も、コードも、同じ言葉を使う。翻訳はいらない。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/order/order.go"}
            </span>
            {"\n"}
            {"type Order struct {\n"}
            {"    id         OrderID\n"}
            {"    customerID CustomerID\n"}
            {"    items      []OrderItem\n"}
            {"    status     Status\n"}
            {"}"}
          </>
        ),
        note: (
          <>
            業務で「注文」と呼ぶものは、コードでも変わらず{" "}
            <code className="font-mono">Order</code> と呼ばれる。フォルダ名も{" "}
            <code className="font-mono">internal/domain/order</code>{" "}
            のまま。
          </>
        ),
      }}
      summary="現場の言葉とコードの言葉をそろえると、翻訳のズレによるバグが減る。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <UbiquitousLanguageDiagram />
        </figure>
      </Reveal>
    </ConceptLayout>
  );
}
