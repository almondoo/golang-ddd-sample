/**
 * ドメインページの図解。
 * 実体はpublic/diagrams/domain-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function DomainDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/domain-diagram.svg"
      alt="商品・カート・注文・配送のミニ情景を抱えるEC業務の世界が、矢印を通ってinternal/domain/配下のコンテキストのコード箱へ写し取られる図"
      width={620}
      height={400}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
