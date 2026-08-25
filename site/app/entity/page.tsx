import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { EntityDiagram } from "@/components/diagrams/entity-diagram";

export const metadata: Metadata = {
  title: "エンティティとは",
  description:
    "名前が変わってもあなたはあなた。エンティティの考え方をGo DDDサンプルのCustomer型で図解します。",
};

export default function EntityPage() {
  return (
    <ConceptLayout
      slug="entity"
      eyebrow="06 / エンティティ"
      title="名前が変わっても、あなたはあなた"
      lead="見た目や名前が変わっても、IDが同じなら同じ人。それがエンティティ。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/customer/customer.go"}
            </span>
            {"\n"}
            {"type Customer struct {\n"}
            {"    id   CustomerID "}
            <span className="text-background/50">{"// 同一性の唯一の根拠"}</span>
            {"\n"}
            {"    name string      "}
            <span className="text-background/50">{"// 変わってもよい"}</span>
            {"\n}\n\n"}
            <span className="text-background/50">{"// 等価性はIDだけで判定する"}</span>
            {"\n"}
            {"func (c *Customer) Equals(other *Customer) bool {\n"}
            {"    return c.id == other.id\n"}
            {"}"}
          </>
        ),
        note: (
          <>
            <code className="font-mono">Customer</code> は名前(
            <code className="font-mono">name</code>)が変わっても、
            <code className="font-mono">CustomerID</code>{" "}
            が同じなら同一人物として扱う。
          </>
        ),
      }}
      summary="エンティティは「IDで区別する」もの。中身が変わっても、IDが同じなら同一。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <EntityDiagram />
        </figure>
      </Reveal>
    </ConceptLayout>
  );
}
