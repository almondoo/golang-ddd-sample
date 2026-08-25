import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { OptimisticLockingDiagram } from "@/components/diagrams/optimistic-locking-diagram";

export const metadata: Metadata = {
  title: "楽観的ロックとは",
  description:
    "版数札つきの荷物。楽観的ロック(Optimistic Lock)の考え方をGo DDDサンプルのinventory.Stockで図解します。",
};

export default function OptimisticLockingPage() {
  return (
    <ConceptLayout
      slug="optimistic-locking"
      eyebrow="17 / 楽観的ロック"
      title="版数札つきの、荷物"
      lead="読んだ時の札(version)と、今の札が違ったら書き込ませない。ロックはかけず、書く瞬間だけ確認する。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/infrastructure/persistence/inventory_repository.go"}
            </span>
            {"\n"}
            {"// version == 0(未永続化)なら version=1 で INSERT\n"}
            {"// version >= 1 なら 条件付き UPDATE\n"}
            {"UPDATE stocks SET reserved = ?, version = version + 1\n"}
            {"WHERE product_id = ? AND version = ?\n\n"}
            <span className="text-background/50">
              {"// RowsAffected == 0 → 他のtxが先に更新した(競合)"}
            </span>
            {"\n"}
            {"if rowsAffected == 0 {\n"}
            {"    return shared.ErrConflict\n"}
            {"}"}
          </>
        ),
        note: (
          <>
            <code className="font-mono">shared.ErrConflict</code>{" "}
            はプレゼンテーション層でHTTP 409に変換される。「読み直して再試行すれば解消しうる」クライアント側の問題として扱う。
          </>
        ),
      }}
      summary="楽観的ロックは競合を「起こさない」仕組みではなく「検出する」仕組み。競合したら読み直して再試行する。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <OptimisticLockingDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        本リポジトリはStockだけにこの仕組みを実装している。「同じ商品を複数の注文が同時に引き当てる」lost updateが最も起きやすく、実害(在庫超過)も分かりやすい箇所だから。
      </p>
    </ConceptLayout>
  );
}
