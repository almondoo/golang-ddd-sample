/** リポジトリページの図解(実体はpublic/diagrams/repository-diagram.svg、本コンポーネントはそれを読み込む薄いラッパー) */
export function RepositoryDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/repository-diagram.svg"
      alt="呼び出し側が倉庫係に「保存して」「持ってきて」とだけ話し、SQLはカーテンの向こうに隠れている図"
      width={560}
      height={380}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
