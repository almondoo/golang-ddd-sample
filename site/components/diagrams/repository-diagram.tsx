/** リポジトリページの図解(新規作成)。サイトのデザイントークンを色に使う */
export function RepositoryDiagram() {
  return (
    <svg viewBox="0 0 560 380" role="img" aria-labelledby="repo-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="repo-title">
        呼び出し側が倉庫係に「保存して」「持ってきて」とだけ話し、SQLはカーテンの向こうに隠れている図
      </title>

      <rect x="20" y="20" width="520" height="340" rx="28" fill="var(--muted)" />

      {/* 左: 呼び出し側(UseCase) */}
      <g transform="translate(130,235)">
        <rect x="-30" y="-10" width="60" height="70" rx="22" fill="var(--primary)" />
        <circle cx="0" cy="-38" r="26" fill="var(--background)" />
        <circle cx="-8" cy="-40" r="3" fill="var(--foreground)" />
        <circle cx="8" cy="-40" r="3" fill="var(--foreground)" />
        <path d="M -9 -28 Q 0 -21 9 -28" stroke="var(--foreground)" strokeWidth="2.6" fill="none" strokeLinecap="round" />
      </g>
      <text x="130" y="330" textAnchor="middle" fontWeight="700" fontSize="14" fill="var(--foreground)">
        呼び出す側(UseCase)
      </text>

      {/* 吹き出し */}
      <rect x="185" y="130" width="150" height="40" rx="18" fill="var(--card)" stroke="var(--primary)" strokeWidth="2" />
      <text x="260" y="155" textAnchor="middle" fontSize="14" fontWeight="700" fill="var(--foreground)">
        保存して
      </text>
      <rect x="175" y="185" width="200" height="40" rx="18" fill="var(--card)" stroke="var(--primary)" strokeWidth="2" />
      <text x="275" y="210" textAnchor="middle" fontSize="13" fontWeight="700" fill="var(--foreground)">
        ID: P-001を持ってきて
      </text>

      {/* 右: 倉庫係(Repository実装) */}
      <g transform="translate(420,235)">
        <rect x="-30" y="-10" width="60" height="70" rx="22" fill="var(--accent)" />
        <circle cx="0" cy="-38" r="26" fill="var(--background)" />
        <circle cx="-8" cy="-40" r="3" fill="var(--foreground)" />
        <circle cx="8" cy="-40" r="3" fill="var(--foreground)" />
        <path d="M -9 -28 Q 0 -21 9 -28" stroke="var(--foreground)" strokeWidth="2.6" fill="none" strokeLinecap="round" />
      </g>
      <text x="420" y="330" textAnchor="middle" fontWeight="700" fontSize="14" fill="var(--foreground)">
        倉庫係(Repository)
      </text>

      {/* 倉庫係の裏に隠れたSQL(カーテンの向こう) */}
      <rect x="470" y="90" width="70" height="130" rx="14" fill="var(--foreground)" />
      {/* 南京錠アイコン(絵文字ではなく図形で描く) */}
      <path
        d="M 495 140 a 10 10 0 0 1 20 0 v 8 h -20 z"
        fill="none"
        stroke="var(--background)"
        strokeWidth="3.5"
        opacity="0.9"
      />
      <rect x="490" y="148" width="30" height="24" rx="5" fill="var(--background)" opacity="0.9" />
      <text x="505" y="195" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="15" fill="var(--background)" opacity="0.9">
        SQL
      </text>
    </svg>
  );
}
