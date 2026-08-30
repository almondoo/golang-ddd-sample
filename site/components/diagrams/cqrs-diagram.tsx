/**
 * CQRSページの図解。
 * 実体はpublic/diagrams/cqrs-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function CqrsDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/cqrs-diagram.svg"
      alt="commandはUseCase→Repository→Product集約(検証)→GORM実装という経路でドメインを通り、queryはUseCase→QueryService→GORM実装(DTO直行)という経路でドメインを迂回し、両方とも同じPostgreSQLへ書き込み・問い合わせすることを示す図"
      width={560}
      height={440}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
