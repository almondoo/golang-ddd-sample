/**
 * 集約ページの図解。
 * 実体はpublic/diagrams/aggregate-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function AggregateDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/aggregate-diagram.svg"
      alt="おもちゃ箱(Order)の中身(OrderItem)を直接さわろうとすると止められ、AddItem()経由なら受け付けられることを示す図"
      width={560}
      height={380}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
