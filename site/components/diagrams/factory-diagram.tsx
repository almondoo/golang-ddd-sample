/**
 * ファクトリページの図解。
 * 実体はpublic/diagrams/factory-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function FactoryDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/factory-diagram.svg"
      alt="新規生成はNewProductが検証してから*Productを作り、再構築はReconstructProductがDBの検証済みデータをそのまま復元する、2つの経路を示す図"
      width={600}
      height={400}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
