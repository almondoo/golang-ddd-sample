/**
 * 共有カーネルページの図解(新規作成)。
 * 7つのコンテキストがすべて中央のshared(道具箱)を一方向に参照し、
 * 逆方向の矢印が存在しないことを示す。
 */
export function SharedKernelDiagram() {
  const contexts = [
    { name: "catalog", x: 280, y: 65, sx: 280, sy: 87, ex: 280, ey: 155 },
    { name: "cart", x: 409, y: 127, sx: 392, sy: 141, ex: 356, ey: 170 },
    { name: "order", x: 441, y: 267, sx: 420, sy: 262, ex: 375, ey: 252 },
    { name: "customer", x: 352, y: 379, sx: 342, sy: 359, ex: 313, ey: 298 },
    { name: "inventory", x: 208, y: 379, sx: 218, sy: 359, ex: 248, ey: 298 },
    { name: "shipping", x: 119, y: 267, sx: 141, sy: 262, ex: 185, ey: 252 },
    { name: "coupon", x: 151, y: 127, sx: 168, sy: 141, ex: 204, ey: 170 },
  ];

  return (
    <svg viewBox="0 0 560 440" role="img" aria-labelledby="sk-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="sk-title">
        7つのコンテキストが、中央のshared(共有カーネル)を一方向に参照している図。矢印はすべて内向きで、逆方向は存在しない
      </title>
      <defs>
        <marker id="sk-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--primary)" />
        </marker>
      </defs>

      <rect x="10" y="10" width="540" height="420" rx="28" fill="var(--muted)" />

      {contexts.map((c) => (
        <g key={c.name}>
          <line x1={c.sx} y1={c.sy} x2={c.ex} y2={c.ey} stroke="var(--primary)" strokeWidth="4" strokeLinecap="round" markerEnd="url(#sk-arrow)" />
          <rect x={c.x - 48} y={c.y - 20} width="96" height="40" rx="14" fill="var(--card)" stroke="var(--border)" strokeWidth="2" />
          <text x={c.x} y={c.y + 5} textAnchor="middle" className="font-mono" fontWeight="700" fontSize="13" fill="var(--foreground)">
            {c.name}
          </text>
        </g>
      ))}

      {/* 中央: shared(共有カーネル) */}
      <rect x="195" y="175" width="170" height="110" rx="24" fill="var(--accent)" />
      <text x="280" y="205" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="18" fill="var(--accent-foreground)">
        shared
      </text>
      <text x="280" y="228" textAnchor="middle" className="font-mono" fontSize="12" fill="var(--accent-foreground)" opacity="0.9">
        Money
      </text>
      <text x="280" y="246" textAnchor="middle" className="font-mono" fontSize="12" fill="var(--accent-foreground)" opacity="0.9">
        ErrNotFound
      </text>
      <text x="280" y="264" textAnchor="middle" className="font-mono" fontSize="12" fill="var(--accent-foreground)" opacity="0.9">
        ErrConflict
      </text>

      <text x="280" y="420" textAnchor="middle" fontSize="14" fill="var(--foreground)">
        矢印はいつも内向き。sharedは誰にも依存しない
      </text>
    </svg>
  );
}
