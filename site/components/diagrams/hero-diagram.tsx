/**
 * トップページのヒーロー図解。
 * 実体はpublic/diagrams/hero-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 * トップページの折り返し線より上に表示されるため、遅延読み込み(loading="lazy")は付けない。
 */
export function HeroDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/hero-diagram.svg"
      alt="現場で使われる言葉が、DDDを通ってそのままきれいなコードになる図"
      width={640}
      height={380}
      className="mx-auto block h-auto w-full max-w-full"
    />
  );
}
