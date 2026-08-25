/**
 * 集約ページの図解(新規作成)。
 * サイトのデザイントークン(var(--primary)等)をそのままSVGの色として使い、
 * 挿絵をUIパレットと統一している(旧ページの挿絵は原本の配色を維持している点と対照的)。
 */
export function AggregateDiagram() {
  return (
    <svg viewBox="0 0 560 380" role="img" aria-labelledby="agg-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="agg-title">
        おもちゃ箱(Order)の中身(OrderItem)を直接さわろうとすると止められ、AddItem()経由なら受け付けられることを示す図
      </title>

      <rect x="20" y="20" width="520" height="340" rx="28" fill="var(--muted)" />

      {/* おもちゃ箱(Order集約) */}
      <polygon points="180,170 215,120 345,120 380,170" fill="var(--accent)" opacity="0.7" />
      <rect x="180" y="170" width="200" height="140" rx="18" fill="var(--accent)" />

      {/* 箱の中の3つのOrderItem */}
      <rect x="196" y="205" width="46" height="80" rx="10" fill="var(--primary)" />
      <rect x="257" y="205" width="46" height="80" rx="10" fill="var(--secondary)" />
      <rect x="318" y="205" width="46" height="80" rx="10" fill="var(--primary)" />

      {/* 「Order」ネームプレート */}
      <rect x="228" y="146" width="104" height="34" rx="17" fill="var(--foreground)" />
      <text x="280" y="168" textAnchor="middle" fill="var(--background)" fontWeight="800" fontSize="16">
        Order
      </text>

      {/* 左: 直接さわろうとして止められる */}
      <line x1="55" y1="240" x2="150" y2="240" stroke="var(--destructive)" strokeWidth="7" strokeLinecap="round" />
      <polygon points="150,240 132,230 132,250" fill="var(--destructive)" />
      <circle cx="103" cy="240" r="36" fill="none" stroke="var(--destructive)" strokeWidth="6" />
      <line x1="79" y1="216" x2="127" y2="264" stroke="var(--destructive)" strokeWidth="6" strokeLinecap="round" />
      <text x="103" y="300" textAnchor="middle" fill="var(--foreground)" fontWeight="700" fontSize="14">
        中身を直接さわる
      </text>

      {/* 右: AddItem()経由なら受け付けられる */}
      <line x1="505" y1="110" x2="400" y2="150" stroke="var(--primary)" strokeWidth="7" strokeLinecap="round" />
      <polygon points="400,150 419,142 411,161" fill="var(--primary)" />
      <circle cx="470" cy="90" r="20" fill="var(--primary)" />
      <path d="M 461 90 L 468 98 L 481 82" stroke="var(--background)" strokeWidth="4" fill="none" strokeLinecap="round" strokeLinejoin="round" />
      <text x="470" y="130" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="14" fill="var(--foreground)">
        AddItem()
      </text>

      <text x="280" y="350" textAnchor="middle" fontSize="14" fill="var(--foreground)">
        不変条件は、箱(集約ルート)が守る
      </text>
    </svg>
  );
}
