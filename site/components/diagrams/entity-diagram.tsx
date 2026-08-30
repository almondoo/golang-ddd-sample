/** エンティティページの図解(実体はpublic/diagrams/entity-diagram.svg、本コンポーネントはそれを読み込む薄いラッパー) */
export function EntityDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/entity-diagram.svg"
      alt="名前が「田中さん」から「田中一郎さん」に変わっても、ID: C-001は変わらないことを示す図"
      width={520}
      height={360}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
