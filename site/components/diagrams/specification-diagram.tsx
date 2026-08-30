/**
 * 仕様(Specification)ページの図解。
 * 実体はpublic/diagrams/specification-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function SpecificationDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/specification-diagram.svg"
      alt="「在庫あり」and「特定カテゴリ」という条件を1つのSpecificationオブジェクトにまとめ、全商品から条件に合うものだけを絞り込む図(教科書的な一般像)"
      width={560}
      height={380}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
