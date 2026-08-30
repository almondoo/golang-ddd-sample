/**
 * 共有カーネルページの図解。
 * 実体はpublic/diagrams/shared-kernel-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function SharedKernelDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/shared-kernel-diagram.svg"
      alt="7つのコンテキストが、中央のshared(共有カーネル)を一方向に参照している図。矢印はすべて内向きで、逆方向は存在しない"
      width={560}
      height={440}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
