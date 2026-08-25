import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { OnionDiagram } from "@/components/diagrams/onion-diagram";

export const metadata: Metadata = {
  title: "オニオンアーキテクチャとは",
  description:
    "まん中を、まわりが守る。オニオンアーキテクチャの考え方をGo DDDサンプルの層構成(internal/domain等)で図解します。",
};

export default function OnionArchitecturePage() {
  return (
    <ConceptLayout
      slug="onion-architecture"
      eyebrow="16 / オニオンアーキテクチャ"
      title="まん中を、まわりが守る"
      lead="玉ねぎみたいに層を重ねる。矢印はいつも内向き。まん中(ドメイン)はまわりを知らない。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// 依存の向きはいつも内向き(外側 → 内側)のみ許される"}
            </span>
            {"\n"}
            {"internal/presentation   "}
            <span className="text-background/50">{"// HTTPコントローラ"}</span>
            {"\n"}
            {"internal/infrastructure "}
            <span className="text-background/50">{"// GORM等の実装"}</span>
            {"\n"}
            {"internal/application    "}
            <span className="text-background/50">{"// ユースケース"}</span>
            {"\n"}
            {"internal/domain         "}
            <span className="text-background/50">{"// 集約・値オブジェクト等"}</span>
          </>
        ),
        note: "internal/domainはどのパッケージもimportしない。外側の層だけが、より内側の層をimportしてよい。",
      }}
      summary="依存はいつも内向き。ドメインは誰にも頼らない、いちばん自由な層。"
    >
      <Reveal>
        <figure className="mx-auto max-w-xl">
          <OnionDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        外側(presentation/infrastructure)は内側(application)を知っていてよいが、内側は外側を知らない。いちばん内側のdomainは何にも依存しない。
      </p>
    </ConceptLayout>
  );
}
