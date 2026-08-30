/**
 * 楽観的ロックページの図解。
 * 実体はpublic/diagrams/optimistic-locking-diagram.svgとして配信される静的SVGで、
 * 本コンポーネントはそれを読み込む薄いラッパーに過ぎない。
 */
export function OptimisticLockingDiagram() {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- 静的SVGのため最適化不要
    <img
      src="/diagrams/optimistic-locking-diagram.svg"
      alt="トランザクションAとBが同じStock(version=1)を同時に読み込み、先にUPDATEしたAはversion=2で成功、後発のBはWHERE version=1に一致する行がなく shared.ErrConflict になることを示す図"
      width={560}
      height={440}
      className="mx-auto block h-auto w-full max-w-full"
      loading="lazy"
    />
  );
}
