/**
 * アプリケーションサービスページの図解。
 * 実体はpublic/diagrams/application-service-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function ApplicationServiceDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/application-service-diagram.svg"
      alt="段取り係(PlaceOrderUseCase)が、1トランザクションの枠の中で6つの手順を順番に指示していく図"
      width={520}
      height={630}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
