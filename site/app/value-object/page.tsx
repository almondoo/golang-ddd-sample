import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { ValueObjectDiagram } from "@/components/diagrams/value-object-diagram";

export const metadata: Metadata = {
  title: "値オブジェクトとは",
  description:
    "1000円はどの1000円でも同じ。値オブジェクトの考え方をGo DDDサンプルのMoney型で図解します。",
};

export default function ValueObjectPage() {
  return (
    <ConceptLayout
      slug="value-object"
      eyebrow="07 / 値オブジェクト"
      title="1000円は、どの1000円でも同じ"
      lead="IDで区別しない。中身(値)が同じなら、同じもの。それが値オブジェクト。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/shared/money.go"}
            </span>
            {"\n"}
            {"type Money struct {\n"}
            {"    amount   int64 "}
            <span className="text-background/50">
              {"// フィールドは非公開 = 外から書き換えられない"}
            </span>
            {"\n"}
            {"    currency string\n}\n\n"}
            <span className="text-background/50">
              {"// コンストラクタが自分の正しさを自分で検証する(自己検証)"}
            </span>
            {"\n"}
            {"func NewMoney(amount int64, currency string) (Money, error) {\n"}
            {"    if amount < 0 {\n"}
            {"        return Money{}, ErrNegativeAmount\n"}
            {"    }\n"}
            {"    return Money{amount: amount, currency: currency}, nil\n"}
            {"}\n\n"}
            <span className="text-background/50">
              {"// 加算は「新しい値」を返す。元の値は変えない(不変)"}
            </span>
            {"\n"}
            {"func (m Money) Add(other Money) (Money, error) { ... }"}
          </>
        ),
        note: (
          <>
            <code className="font-mono">Money</code>{" "}
            は作られた後に書き換えられない(不変)。金額と通貨という「中身」だけで等しいかどうかが決まる。
          </>
        ),
      }}
      summary="値オブジェクトは「中身で区別する」もの。作った後は変えず、正しさは自分で守る。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <ValueObjectDiagram />
        </figure>
      </Reveal>
    </ConceptLayout>
  );
}
