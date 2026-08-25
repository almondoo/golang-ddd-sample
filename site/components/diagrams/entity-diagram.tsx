/** エンティティページの図解(旧静的サイトより移植、色は原本のまま) */
export function EntityDiagram() {
  return (
    <svg viewBox="0 0 520 360" role="img" aria-labelledby="entity-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="entity-title">
        名前が「田中さん」から「田中一郎さん」に変わっても、ID: C-001は変わらないことを示す図
      </title>

      <rect x="20" y="20" width="480" height="320" rx="24" fill="#F7E1E7" />

      <g transform="translate(260,120)">
        <rect x="-34" y="-14" width="68" height="76" rx="26" fill="#C4577B" />
        <circle cx="0" cy="-44" r="30" fill="#fff" />
        <circle cx="-10" cy="-46" r="3.4" fill="#2B2117" />
        <circle cx="10" cy="-46" r="3.4" fill="#2B2117" />
        <path d="M -11 -34 Q 0 -25 11 -34" stroke="#2B2117" strokeWidth="3.2" fill="none" strokeLinecap="round" />
      </g>

      <rect x="80" y="200" width="140" height="40" rx="20" fill="#fff" stroke="#C4577B" strokeWidth="2" />
      <text x="150" y="225" textAnchor="middle" fontSize="16" fill="#2B2117">
        田中さん
      </text>
      <text x="260" y="225" textAnchor="middle" fontSize="22" fill="#C4577B">
        →
      </text>
      <rect x="300" y="200" width="140" height="40" rx="20" fill="#fff" stroke="#C4577B" strokeWidth="2" />
      <text x="370" y="225" textAnchor="middle" fontSize="16" fill="#2B2117">
        田中一郎さん
      </text>

      <rect x="185" y="264" width="150" height="44" rx="22" fill="#C4577B" />
      <text x="260" y="292" textAnchor="middle" fill="#fff" fontWeight="700" className="font-mono" fontSize="15">
        ID: C-001
      </text>

      <text x="260" y="330" textAnchor="middle" fontSize="14" fill="#5A4A3D">
        名前は変わったけど、IDは変わらない = 同じ人
      </text>
    </svg>
  );
}
