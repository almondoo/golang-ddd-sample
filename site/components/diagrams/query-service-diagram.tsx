/**
 * クエリサービスページの図解。
 * 実体はpublic/diagrams/query-service-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function QueryServiceDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/query-service-diagram.svg"
      alt="書く窓口はUseCaseから集約(ドメイン)を経てRepositoryへ進みDBに書き込むのに対し、見る窓口はUseCaseからQueryServiceへ進みドメイン箱を素通りしてDTOを直接DBから組み立てることを示す図"
      width={600}
      height={460}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
