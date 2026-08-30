/**
 * ドメインイベントページの図解。
 * 実体はpublic/diagrams/domain-event-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function DomainEventDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/domain-event-diagram.svg"
      alt="以前はイベントバス経由でCartハンドラが購読していたが、現在は学習コストを下げるためPlaceOrderUseCaseがorderRepoとcartRepoを直接呼び出す、という変遷を示す図"
      width={600}
      height={440}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
