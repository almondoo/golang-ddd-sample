/** 境界づけられたコンテキストページの図解(実体はpublic/diagrams/bounded-context-diagram.svg、本コンポーネントはそれを読み込む薄いラッパー) */
export function BoundedContextDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/bounded-context-diagram.svg"
      alt="お店の中に7つの部屋(catalog・cart・order・customer・inventory・shipping・coupon)が並ぶ見取り図"
      width={700}
      height={500}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
