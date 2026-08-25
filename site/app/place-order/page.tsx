import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { PlaceOrderFlow } from "@/components/place-order-flow";

export const metadata: Metadata = {
  title: "実例で見る: 注文確定の旅",
  description:
    "PlaceOrderの一部始終。HTTPリクエストからOrder保存までの流れと、エラーからHTTPステータスへの対応をGo DDDサンプルで図解します。",
};

export default function PlaceOrderPage() {
  return (
    <ConceptLayout
      slug="place-order"
      eyebrow="19 / 実例で見る: 注文確定の旅"
      title="注文確定の旅"
      lead="「注文して」という1回のお願いが、裏でいくつものステップを経てOrderになる。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/application/usecase/order/place_order.go\n"}
              {"// PlaceOrderUseCase.Execute はすべての手順を1トランザクションにまとめる"}
            </span>
            {"\n"}
            {"err = uc.txManager.Do(ctx, func(ctx context.Context) error {\n"}
            {"    "}
            <span className="text-background/50">{"// 1. 顧客確認 → 2. カート読込 → 3. 在庫確保"}</span>
            {"\n"}
            {"    "}
            <span className="text-background/50">{"// 4. Order生成 → 5. クーポン適用 → 6. 保存 → 7. カートを空にする"}</span>
            {"\n"}
            {"    ...\n"}
            {"    return nil\n"}
            {"})"}
          </>
        ),
        note: "途中でどれか1つでも失敗すれば、トランザクションごと巻き戻る。「注文は確定したがカートは空にならなかった」のような中途半端な状態は起こらない。",
      }}
      summary="1つのユースケースが、1つのトランザクションの中で筋道立てて処理を進める。"
    >
      <PlaceOrderFlow />
    </ConceptLayout>
  );
}
