/**
 * 楽観的ロックページの図解(新規作成)。
 * 版数札(version)つきの荷物(Stock)を2つのトランザクションが同時に読み、
 * 先に書き込んだ方だけが成功し、後発は版数の不一致で弾かれることを示す。
 */
export function OptimisticLockingDiagram() {
  return (
    <svg viewBox="0 0 560 440" role="img" aria-labelledby="ol-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="ol-title">
        トランザクションAとBが同じStock(version=1)を同時に読み込み、先にUPDATEしたAはversion=2で成功、後発のBはWHERE version=1に一致する行がなく shared.ErrConflict になることを示す図
      </title>
      <defs>
        <marker id="ol-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--border)" />
        </marker>
      </defs>

      {/* 版数札つきの荷物(Stock) */}
      <rect x="210" y="25" width="140" height="85" rx="18" fill="var(--primary)" />
      <text x="280" y="60" textAnchor="middle" fontWeight="800" fontSize="16" fill="var(--primary-foreground)">Stock</text>
      <text x="280" y="82" textAnchor="middle" className="font-mono" fontSize="11" fill="var(--primary-foreground)" opacity="0.9">quantity=10</text>
      <circle cx="335" cy="30" r="20" fill="var(--accent)" />
      <text x="335" y="35" textAnchor="middle" fontWeight="800" fontSize="13" fill="var(--accent-foreground)">v1</text>

      {/* トランザクションA・B(人物) */}
      <g transform="translate(90,190)">
        <circle cx="0" cy="0" r="24" fill="var(--chart-2)" />
        <circle cx="-7" cy="-3" r="3" fill="var(--foreground)" />
        <circle cx="7" cy="-3" r="3" fill="var(--foreground)" />
        <path d="M -8 8 Q 0 15 8 8" stroke="var(--foreground)" strokeWidth="2.5" fill="none" strokeLinecap="round" />
      </g>
      <text x="90" y="240" textAnchor="middle" fontWeight="700" fontSize="13" fill="var(--foreground)">トランザクションA</text>

      <g transform="translate(470,190)">
        <circle cx="0" cy="0" r="24" fill="var(--chart-3)" />
        <circle cx="-7" cy="-3" r="3" fill="var(--foreground)" />
        <circle cx="7" cy="-3" r="3" fill="var(--foreground)" />
        <path d="M -8 8 Q 0 15 8 8" stroke="var(--foreground)" strokeWidth="2.5" fill="none" strokeLinecap="round" />
      </g>
      <text x="470" y="240" textAnchor="middle" fontWeight="700" fontSize="13" fill="var(--foreground)">トランザクションB</text>

      {/* 両方が version=1 を読む */}
      <line x1="105" y1="170" x2="230" y2="115" stroke="var(--border)" strokeWidth="3" strokeDasharray="6 5" markerEnd="url(#ol-arrow)" />
      <line x1="455" y1="170" x2="330" y2="115" stroke="var(--border)" strokeWidth="3" strokeDasharray="6 5" markerEnd="url(#ol-arrow)" />
      <text x="280" y="150" textAnchor="middle" fontSize="12" fill="var(--muted-foreground)">両方とも version=1 を読む</text>

      {/* 結果 */}
      <line x1="105" y1="215" x2="105" y2="300" stroke="var(--primary)" strokeWidth="4" markerEnd="url(#ol-arrow)" />
      <line x1="470" y1="215" x2="470" y2="300" stroke="var(--destructive)" strokeWidth="4" markerEnd="url(#ol-arrow)" />

      <rect x="15" y="305" width="185" height="90" rx="18" fill="var(--card)" stroke="var(--primary)" strokeWidth="3" />
      <path d="M 45 350 L 60 365 L 90 330" stroke="var(--primary)" strokeWidth="5" fill="none" strokeLinecap="round" strokeLinejoin="round" />
      <text x="145" y="345" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">UPDATE ...</text>
      <text x="145" y="360" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">version=2</text>
      <text x="145" y="380" textAnchor="middle" fontWeight="800" fontSize="12" fill="var(--primary)">成功</text>

      <rect x="360" y="305" width="185" height="90" rx="18" fill="var(--muted)" stroke="var(--destructive)" strokeWidth="3" strokeDasharray="7 6" />
      <circle cx="390" cy="345" r="14" fill="none" stroke="var(--destructive)" strokeWidth="3.5" />
      <line x1="381" y1="336" x2="399" y2="354" stroke="var(--destructive)" strokeWidth="3.5" strokeLinecap="round" />
      <text x="475" y="340" textAnchor="middle" fontSize="11" fill="var(--foreground)">WHERE version=1</text>
      <text x="475" y="356" textAnchor="middle" fontSize="11" fill="var(--foreground)">一致する行なし</text>
      <text x="475" y="378" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="11" fill="var(--destructive)">
        shared.ErrConflict
      </text>

      <text x="280" y="425" textAnchor="middle" fontSize="12" fill="var(--muted-foreground)">
        版数が合わない=誰かが先に書き換えた、というサイン
      </text>
    </svg>
  );
}
