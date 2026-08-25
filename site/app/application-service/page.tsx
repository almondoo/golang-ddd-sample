import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { ApplicationServiceDiagram } from "@/components/diagrams/application-service-diagram";

export const metadata: Metadata = {
  title: "アプリケーションサービスとは",
  description:
    "段取り係が手順をまとめる。アプリケーションサービス(usecase)の考え方をGo DDDサンプルのPlaceOrderUseCaseで図解します。",
};

export default function ApplicationServicePage() {
  return (
    <ConceptLayout
      slug="application-service"
      eyebrow="12 / アプリケーションサービス"
      title="段取り係が、手順をまとめる"
      lead="自分ではルールを判断しない。集約に指示を出し、1つのトランザクションにまとめるだけ。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/application/usecase/order/place_order.go"}
            </span>
            {"\n"}
            {"func (uc *PlaceOrderUseCase) Execute(ctx context.Context, in Input) (*Output, error) {\n"}
            {"    "}
            <span className="text-background/50">{"// 手順の組み立てだけ。ルール判断はしない"}</span>
            {"\n"}
            {"    err := uc.txManager.Do(ctx, func(ctx context.Context) error {\n"}
            {"        "}
            <span className="text-background/50">{"// 顧客確認 → カート読込 → 在庫引当"}</span>
            {"\n"}
            {"        "}
            <span className="text-background/50">{"// → Order生成 → クーポン適用 → 保存 → カートを空に"}</span>
            {"\n"}
            {"        ...\n"}
            {"        return nil\n"}
            {"    })\n"}
            {"    ...\n"}
            {"}"}
          </>
        ),
        note: (
          <>
            在庫が足りるか、クーポンが有効か、といった判断はすべて呼び出し先の集約のメソッドが返す。
            <code className="font-mono">PlaceOrderUseCase</code>{" "}
            自身はその結果をそのまま伝えるだけ。
          </>
        ),
      }}
      summary="アプリケーションサービスは「調整役」。ルールは持たず、集約への指示と1トランザクションの組み立てに徹する。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <ApplicationServiceDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        トランザクションを開始・コミットするかどうかを決めるのは、常にアプリケーション層。リポジトリ自身が独自にトランザクションを張ることはない。
      </p>
    </ConceptLayout>
  );
}
