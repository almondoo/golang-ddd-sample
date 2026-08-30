/** オニオンアーキテクチャページの図解(実体はpublic/diagrams/onion-diagram.svg、本コンポーネントはそれを読み込む薄いラッパー) */
export function OnionDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/onion-diagram.svg"
      alt="中心のdomainを、application・infrastructure/presentationの同心円が取り囲み、依存の矢印はすべて内向きであることを示す図"
      width={520}
      height={520}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
