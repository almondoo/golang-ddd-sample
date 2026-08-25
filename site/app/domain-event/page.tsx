import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { NotAdoptedBadge } from "@/components/not-adopted-badge";
import { DomainEventDiagram } from "@/components/diagrams/domain-event-diagram";

export const metadata: Metadata = {
  title: "ドメインイベントとは",
  description:
    "「起きたこと」を知らせる仕組み。ドメインイベント(Domain Event)の考え方と、本リポジトリで一度採用してからやめた経緯を図解します。",
};

export default function DomainEventPage() {
  return (
    <ConceptLayout
      slug="domain-event"
      eyebrow="13 / ドメインイベント"
      title="「起きたこと」を、知らせる仕組み"
      lead="発行する側は、受け取る側を知らないまま連携できる。ただし本リポジトリは一度使って、やめた。"
      badge={<NotAdoptedBadge />}
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/application/usecase/order/place_order.go のコメントより"}
            </span>
            {"\n"}
            <span className="text-background/50">
              {"// 以前はOrderPlacedイベントを発行し、カート側のハンドラが\n"}
              {"// 購読して空にする方式だったが、仕組みの理解コストを下げる\n"}
              {"// ため直接呼び出しに変更した。"}
            </span>
            {"\n\n"}
            {"uc.orderRepo.Save(ctx, order)\n"}
            {"cart.Clear()\n"}
            {"uc.cartRepo.Save(ctx, cart) "}
            <span className="text-background/50">{"// 直接呼び出し"}</span>
          </>
        ),
        note: "internal/infrastructure/README.mdにも「そのためinfrastructure層にイベントバスの実装は存在しない」と明記されている。",
      }}
      summary="「不採用・直接呼び出し」は意図して選んだトレードオフ。結合度は上がるが、処理の流れは一目で追える。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <DomainEventDiagram />
        </figure>
      </Reveal>

      <div className="mx-auto mt-8 max-w-xl rounded-3xl border-2 border-border bg-card px-6 py-5">
        <h3 className="mb-2 text-lg font-extrabold">なぜ不採用?</h3>
        <p className="text-sm text-muted-foreground">
          反応するコンテキストが増えるほど、イベント+購読の疎結合さは効いてくる。だが本リポジトリは学習コストを優先し、いったん導入したイベントバスを取り下げて直接呼び出しに戻した。将来ポイント付与など反応する処理が増えたら、ドメインイベント(あるいはOutboxパターン)への発展を検討するのが自然な流れ。
        </p>
      </div>
    </ConceptLayout>
  );
}
