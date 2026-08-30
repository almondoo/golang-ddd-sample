/**
 * ドメインサービスページの図解。
 * 実体はpublic/diagrams/domain-service-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function DomainServiceDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/domain-service-diagram.svg"
      alt="ロジックの置き場所を判定するフローチャート。1つの集約で判断できれば集約のメソッドへ、同一コンテキストで複数集約をまたぐならドメインサービスへ(本リポジトリに実例なし)、コンテキストをまたぐならアプリケーションサービスへ進む"
      width={660}
      height={520}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
