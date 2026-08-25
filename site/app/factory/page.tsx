import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { FactoryDiagram } from "@/components/diagrams/factory-diagram";

export const metadata: Metadata = {
  title: "ファクトリとは",
  description:
    "正しい作り方を知っている工房。ファクトリ(生成と再構築の分離)の考え方をGo DDDサンプルのNew*/Reconstruct*で図解します。",
};

export default function FactoryPage() {
  return (
    <ConceptLayout
      slug="factory"
      eyebrow="11 / ファクトリ"
      title="正しい作り方を知っている、工房"
      lead="「新しく生まれるデータ」と「すでに検証済みのデータを読み戻すだけ」は、作り方が違う。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/coupon/coupon.go"}
            </span>
            {"\n"}
            <span className="text-background/50">{"// 新規生成: 検証あり"}</span>
            {"\n"}
            {"func NewAmountCoupon(code CouponCode, amount shared.Money, ...) (*Coupon, error) {\n"}
            {"    ...\n"}
            {"}\n"}
            {"func NewRateCoupon(code CouponCode, rate int, ...) (*Coupon, error) {\n"}
            {"    ...\n"}
            {"}\n\n"}
            <span className="text-background/50">
              {"// 再構築: DBから読むだけ。検証は省略(過去に検証済みという前提)"}
            </span>
            {"\n"}
            {"func ReconstructCoupon(id CouponID, ...) *Coupon { ... }"}
          </>
        ),
        note: "この New*(検証あり) / Reconstruct*(検証なし・リポジトリ専用)というペアは、Product・Cart・Order・Customer・Coupon・Stock・Shipmentすべてに一貫している。",
      }}
      summary="ファクトリ専用のクラスは無い。関数名を分けるだけで「新規」と「復元」を取り違えないようにしている。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <FactoryDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        Reconstruct*は「検証をサボっている」のではない。DBの中身はすでに一度NewProductの検証を通過している、という前提に立って再検証を省いているだけ。
      </p>
    </ConceptLayout>
  );
}
