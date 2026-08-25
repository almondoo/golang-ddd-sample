import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { SharedKernelDiagram } from "@/components/diagrams/shared-kernel-diagram";

export const metadata: Metadata = {
  title: "共有カーネルとは",
  description:
    "みんなで使う共通の道具箱。共有カーネル(Shared Kernel)の考え方をGo DDDサンプルのinternal/domain/sharedで図解します。",
};

export default function SharedKernelPage() {
  return (
    <ConceptLayout
      slug="shared-kernel"
      eyebrow="05 / 共有カーネル"
      title="みんなで使う、共通の道具箱"
      lead="部屋(コンテキスト)ごとに言葉は分けるけど、本当に共通なものだけは1つの道具箱を共有する。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/shared/ (実装は3ファイル・約200行)\n"}
            </span>
            {"\n"}
            {"money.go   "}
            <span className="text-background/50">{"// Money値オブジェクト(負値・通貨不一致を拒否)"}</span>
            {"\n"}
            {"errors.go  "}
            <span className="text-background/50">{"// ErrNotFound / ErrConflict / DomainRuleError"}</span>
            {"\n"}
            {"id.go      "}
            <span className="text-background/50">{"// NewID(UUID生成のラッパー)"}</span>
          </>
        ),
        note: (
          <>
            catalog / cart / order など7つのコンテキストがすべて{" "}
            <code className="font-mono">shared</code>{" "}
            を参照するが、逆方向(sharedが特定コンテキストに依存する)矢印は存在しない。
          </>
        ),
      }}
      summary="共有カーネルは「本当にどこでも同じ意味を持つ」ものだけを置く、小さな共通の道具箱。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <SharedKernelDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        CouponCodeのようにコンテキストごとに意味が変わりうる型は、あえて共有せず各コンテキストの中に留める。迷ったら「共有しない」が本リポジトリの一貫した判断。
      </p>
    </ConceptLayout>
  );
}
