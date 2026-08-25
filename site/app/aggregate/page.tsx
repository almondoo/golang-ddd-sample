import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { AggregateDiagram } from "@/components/diagrams/aggregate-diagram";

export const metadata: Metadata = {
  title: "集約とは",
  description:
    "おもちゃ箱ごと出し入れする。集約(Aggregate)の考え方をGo DDDサンプルのOrder/OrderItemで図解します。",
};

export default function AggregatePage() {
  return (
    <ConceptLayout
      slug="aggregate"
      eyebrow="08 / 集約"
      title="おもちゃ箱ごと、出し入れする"
      lead="中身に直接手を伸ばさない。箱を通して、まとめて出し入れする。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/order/order.go\n"}
              {"// Order は自身の明細(OrderItem)を抱える集約ルート"}
            </span>
            {"\n"}
            {"type Order struct {\n"}
            {"    id         OrderID\n"}
            {"    customerID CustomerID\n"}
            {"    items      []OrderItem "}
            <span className="text-background/50">{"// 箱の中身"}</span>
            {"\n    status     Status\n}\n\n"}
            <span className="text-background/50">
              {"// 空の注文は作れない(不変条件は箱=集約ルートが守る)"}
            </span>
            {"\n"}
            {"func NewOrder(customerID CustomerID, items []OrderItem, now time.Time) (*Order, error) {\n"}
            {"    if len(items) == 0 {\n"}
            {'        return nil, shared.NewDomainRuleError("order: order must contain at least one item")\n'}
            {"    }\n"}
            {"    ...\n"}
            {"}\n\n"}
            <span className="text-background/50">
              {"// 呼び出し元が中身を直接書き換えられないよう、コピーを返す"}
            </span>
            {"\n"}
            {"func (o *Order) Items() []OrderItem { ... }"}
          </>
        ),
        note: (
          <>
            隣の集約 <code className="font-mono">Cart</code> も同じパターン:
            中身を直接いじらせず、<code className="font-mono">AddItem</code>
            /<code className="font-mono">RemoveItem</code>{" "}
            という入り口(メソッド)を必ず経由させる。
          </>
        ),
      }}
      summary="集約は「まとめて扱う単位」。入り口(集約ルート)だけがルールを守りながら中身を出し入れする。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <AggregateDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        Order箱の中にOrderItemが入っている。箱の外から中身のOrderItemを勝手に書き換えることはできない。
      </p>
    </ConceptLayout>
  );
}
