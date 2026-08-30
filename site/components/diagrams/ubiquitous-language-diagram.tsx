/** ユビキタス言語ページの図解(実体はpublic/diagrams/ubiquitous-language-diagram.svg、本コンポーネントはそれを読み込む薄いラッパー) */
export function UbiquitousLanguageDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/ubiquitous-language-diagram.svg"
      alt="店員・開発者・コードの3者が、みな「注文」という同じ言葉を話している図"
      width={620}
      height={320}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
