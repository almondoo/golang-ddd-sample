/** 値オブジェクトページの図解(実体はpublic/diagrams/value-object-diagram.svg、本コンポーネントはそれを読み込む薄いラッパー) */
export function ValueObjectDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/value-object-diagram.svg"
      alt="2つの「1000円」が、別々のものでも中身が同じなら同じ扱いになることを示す図"
      width={520}
      height={320}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
