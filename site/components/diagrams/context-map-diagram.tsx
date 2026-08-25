/**
 * コンテキストマップページの図解(新規作成)。
 * order(application層)を中心に、PlaceOrderが直接呼び出す5コンテキスト
 * (customer/cart/catalog/inventory/coupon)への実線矢印と、
 * ShipOrderが呼ぶshippingへの薄い破線矢印を放射状に配置する。
 * 右上の吹き出しで、これが教科書のACL/OHSではなく直接呼び出しであることを明示する。
 */
export function ContextMapDiagram() {
  const primary = [
    { key: "customer", label: "customer", note: "実在確認", x: 280, y: 70 },
    { key: "cart", label: "cart", note: "カート読書", x: 460, y: 150 },
    { key: "catalog", label: "catalog", note: "価格取得", x: 460, y: 330 },
    { key: "inventory", label: "inventory", note: "在庫引当", x: 280, y: 410 },
    { key: "coupon", label: "coupon", note: "クーポン消費", x: 100, y: 330 },
  ];

  return (
    <svg viewBox="0 0 700 480" role="img" aria-labelledby="cm-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="cm-title">
        中央のorder(application層)から、customer・cart・catalog・inventory・couponへ実線の矢印が伸び、shippingへは薄い破線の矢印が伸びる図。右上の吹き出しでACL/OHSのような標準パターンではなく直接呼び出しであることを説明している
      </title>
      <defs>
        <marker id="cm-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--primary)" />
        </marker>
        <marker id="cm-arrow-sub" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--muted-foreground)" />
        </marker>
      </defs>

      <rect x="10" y="10" width="680" height="460" rx="28" fill="var(--muted)" />

      {/* secondary: shipping(ShipOrder、薄い破線) */}
      <line x1="245" y1="215" x2="185" y2="165" stroke="var(--muted-foreground)" strokeWidth="3" strokeDasharray="6 5" markerEnd="url(#cm-arrow-sub)" />
      <rect x="45" y="127" width="110" height="46" rx="14" fill="var(--card)" stroke="var(--muted-foreground)" strokeWidth="2" strokeDasharray="5 4" />
      <text x="100" y="155" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="13" fill="var(--muted-foreground)">shipping</text>
      <text x="100" y="115" textAnchor="middle" fontSize="10.5" fill="var(--muted-foreground)">ShipOrderのみ</text>

      {/* primary: 5コンテキストへの実線 */}
      {primary.map((n) => {
        const dx = n.x - 280;
        const dy = n.y - 240;
        const len = Math.hypot(dx, dy);
        const ux = dx / len;
        const uy = dy / len;
        const x1 = 280 + ux * 70;
        const y1 = 240 + uy * 42;
        const x2 = n.x - ux * 60;
        const y2 = n.y - uy * 30;
        return (
          <g key={n.key}>
            <line x1={x1} y1={y1} x2={x2} y2={y2} stroke="var(--primary)" strokeWidth="3.5" markerEnd="url(#cm-arrow)" />
            <rect x={n.x - 55} y={n.y - 23} width="110" height="46" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
            <text x={n.x} y={n.y - 2} textAnchor="middle" className="font-mono" fontWeight="700" fontSize="13" fill="var(--foreground)">
              {n.label}
            </text>
            <text x={n.x} y={n.y + 14} textAnchor="middle" fontSize="10" fill="var(--muted-foreground)">
              {n.note}
            </text>
          </g>
        );
      })}

      {/* 中央: order(application層) */}
      <rect x="215" y="207" width="130" height="66" rx="20" fill="var(--primary)" />
      <text x="280" y="234" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="15" fill="var(--primary-foreground)">
        order
      </text>
      <text x="280" y="253" textAnchor="middle" fontSize="10.5" fill="var(--primary-foreground)" opacity="0.9">
        application層
      </text>

      {/* 右上: ACL/OHSではなく直接呼び出しである旨の吹き出し */}
      <rect x="530" y="30" width="150" height="92" rx="16" fill="var(--card)" stroke="var(--accent)" strokeWidth="2.5" strokeDasharray="6 5" />
      <text x="605" y="56" textAnchor="middle" fontWeight="700" fontSize="11.5" fill="var(--accent)">
        Customer-Supplier
      </text>
      <text x="605" y="72" textAnchor="middle" fontWeight="700" fontSize="11.5" fill="var(--accent)">
        ACL・OHSではない
      </text>
      <text x="605" y="93" textAnchor="middle" fontSize="10.5" fill="var(--foreground)">
        ただの直接呼び出し
      </text>
      <text x="605" y="109" textAnchor="middle" fontSize="10.5" fill="var(--foreground)">
        (正直な trade-off)
      </text>
    </svg>
  );
}
