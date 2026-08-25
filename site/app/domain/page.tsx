import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { DomainDiagram } from "@/components/diagrams/domain-diagram";

export const metadata: Metadata = {
  title: "ドメインとは",
  description:
    "ソフトの前に、業務の世界がある。DDDの出発点であるドメイン(業務の世界)の考え方をGo DDDサンプルのECの世界で図解します。",
};

export default function DomainPage() {
  return (
    <ConceptLayout
      slug="domain"
      eyebrow="01 / ドメイン"
      title="ソフトの前に、業務の世界がある"
      lead="商品を売り、カートに入れ、注文し、届ける。DDDはこの「世界」を中心に据える。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// README.mdより\n"}
              {"// EC(ネットショップ)を題材に、商品カタログ・カート・注文・顧客・\n"}
              {"// 在庫・配送・クーポンの7コンテキストを実装している"}
            </span>
            {"\n\n"}
            {"internal/domain/\n"}
            {"├── catalog/    "}
            <span className="text-background/50">{"// 商品カタログ"}</span>
            {"\n├── cart/       "}
            <span className="text-background/50">{"// カート"}</span>
            {"\n├── order/      "}
            <span className="text-background/50">{"// 注文"}</span>
            {"\n├── customer/   "}
            <span className="text-background/50">{"// 顧客"}</span>
            {"\n├── inventory/  "}
            <span className="text-background/50">{"// 在庫"}</span>
            {"\n├── shipping/   "}
            <span className="text-background/50">{"// 配送"}</span>
            {"\n└── coupon/     "}
            <span className="text-background/50">{"// クーポン"}</span>
          </>
        ),
        note: (
          <>
            Evansの定義では、ドメインとは「利用者の活動や関心の領域」。この7つのディレクトリ名はコードの都合で決めた分割ではなく、EC業務そのものの区切り方をそのまま写し取ったもの。
          </>
        ),
      }}
      summary="まずソフトの前に、業務の世界がある。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <DomainDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        ソフトウェアが解決したいのは「業務の世界」そのもの。このリポジトリならECの世界(商品を売り、カートに入れ、注文し、届ける)。DDDはコードの都合ではなく、この世界を中心に据えて設計する。
      </p>
    </ConceptLayout>
  );
}
