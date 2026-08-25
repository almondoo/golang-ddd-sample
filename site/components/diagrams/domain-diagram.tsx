/**
 * ドメインページの図解(新規作成)。
 * 左に「業務の世界」(EC: 商品・カート・注文・配送のミニ情景)を描き、
 * 矢印でそれが internal/domain/ 配下のパッケージへ写し取られることを示す。
 */
export function DomainDiagram() {
  const scenes = [
    { x: 118, y: 118, label: "商品", fill: "var(--chart-3)" },
    { x: 224, y: 118, label: "カート", fill: "var(--chart-2)" },
    { x: 118, y: 224, label: "注文", fill: "var(--primary)" },
    { x: 224, y: 224, label: "配送", fill: "var(--chart-5)" },
  ];

  return (
    <svg viewBox="0 0 620 400" role="img" aria-labelledby="domain-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="domain-title">
        商品・カート・注文・配送のミニ情景を抱えるEC業務の世界が、矢印を通ってinternal/domain/配下のコンテキストのコード箱へ写し取られる図
      </title>
      <defs>
        <marker id="domain-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--foreground)" />
        </marker>
      </defs>

      {/* 左: 業務の世界(EC) */}
      <circle cx="171" cy="171" r="150" fill="var(--secondary)" opacity="0.35" />
      <circle cx="171" cy="171" r="150" fill="none" stroke="var(--primary)" strokeWidth="3" />
      <text x="171" y="45" textAnchor="middle" fontWeight="800" fontSize="16" fill="var(--foreground)">
        業務の世界(EC)
      </text>

      {scenes.map((s) => (
        <g key={s.label}>
          <rect x={s.x - 44} y={s.y - 34} width="88" height="68" rx="16" fill={s.fill} />
          <text x={s.x} y={s.y + 6} textAnchor="middle" fontWeight="800" fontSize="15" fill="var(--background)">
            {s.label}
          </text>
        </g>
      ))}

      <text x="171" y="345" textAnchor="middle" fontSize="13" fill="var(--foreground)">
        コードの都合ではなく、この世界を中心に据える
      </text>

      {/* 矢印 */}
      <line x1="330" y1="200" x2="392" y2="200" stroke="var(--foreground)" strokeWidth="5" strokeLinecap="round" markerEnd="url(#domain-arrow)" />

      {/* 右: internal/domain/ へ写し取られたコード */}
      <rect x="400" y="60" width="200" height="280" rx="20" fill="var(--foreground)" />
      <text x="500" y="90" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="13" fill="var(--background)">
        internal/domain/
      </text>
      {["catalog", "cart", "order", "customer", "inventory", "shipping", "coupon"].map((name, i) => (
        <rect
          key={name}
          x="416"
          y={105 + i * 33}
          width="168"
          height="25"
          rx="8"
          fill="var(--card)"
          opacity="0.92"
        />
      ))}
      {["catalog", "cart", "order", "customer", "inventory", "shipping", "coupon"].map((name, i) => (
        <text
          key={name}
          x="500"
          y={105 + i * 33 + 17}
          textAnchor="middle"
          className="font-mono"
          fontSize="12"
          fill="var(--foreground)"
        >
          {name}
        </text>
      ))}
    </svg>
  );
}
