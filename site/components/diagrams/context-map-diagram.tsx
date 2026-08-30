/**
 * コンテキストマップページの図解。
 * 実体はpublic/diagrams/context-map-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function ContextMapDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/context-map-diagram.svg"
      alt="中央のorder(application層)から、customer・cart・catalog・inventory・couponへ実線の矢印が伸び、shippingへは薄い破線の矢印が伸びる図。右上の吹き出しでACL/OHSのような標準パターンではなく直接呼び出しであることを説明している"
      width={700}
      height={480}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
