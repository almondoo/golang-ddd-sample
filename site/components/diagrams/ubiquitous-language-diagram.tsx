/** ユビキタス言語ページの図解(旧静的サイトより移植、色は原本のまま) */
export function UbiquitousLanguageDiagram() {
  return (
    <svg viewBox="0 0 620 320" role="img" aria-labelledby="ul-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="ul-title">
        店員・開発者・コードの3者が、みな「注文」という同じ言葉を話している図
      </title>

      {/* 3つの吹き出し(すべて「注文」) */}
      <g transform="translate(100,60)">
        <rect x="-70" y="-30" width="140" height="56" rx="20" fill="#C4577B" />
        <polygon points="-10,26 10,26 0,44" fill="#C4577B" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="800" fontSize="20">
          注文
        </text>
      </g>
      <g transform="translate(310,50)">
        <rect x="-70" y="-30" width="140" height="56" rx="20" fill="#3E6FB0" />
        <polygon points="-10,26 10,26 0,44" fill="#3E6FB0" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="800" fontSize="20">
          注文
        </text>
      </g>
      <g transform="translate(520,60)">
        <rect x="-70" y="-30" width="140" height="56" rx="20" fill="#E8672E" />
        <polygon points="-10,26 10,26 0,44" fill="#E8672E" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="800" fontSize="20">
          注文
        </text>
      </g>

      {/* 店員キャラクター */}
      <g transform="translate(100,190)">
        <rect x="-32" y="-14" width="64" height="72" rx="24" fill="#C4577B" />
        <rect x="-32" y="-14" width="64" height="18" rx="9" fill="#a5486a" />
        <circle cx="0" cy="-42" r="28" fill="#F7E1E7" />
        <circle cx="-9" cy="-44" r="3.2" fill="#2B2117" />
        <circle cx="9" cy="-44" r="3.2" fill="#2B2117" />
        <path d="M -10 -32 Q 0 -24 10 -32" stroke="#2B2117" strokeWidth="3" fill="none" strokeLinecap="round" />
      </g>
      <text x="100" y="270" textAnchor="middle" fontWeight="700" fontSize="16">
        店員さん
      </text>

      {/* 開発者キャラクター */}
      <g transform="translate(310,190)">
        <rect x="-32" y="-14" width="64" height="72" rx="24" fill="#3E6FB0" />
        <circle cx="0" cy="-42" r="28" fill="#DCE7F5" />
        <circle cx="-9" cy="-44" r="3.2" fill="#2B2117" />
        <circle cx="9" cy="-44" r="3.2" fill="#2B2117" />
        <path d="M -10 -32 Q 0 -24 10 -32" stroke="#2B2117" strokeWidth="3" fill="none" strokeLinecap="round" />
        <rect x="-26" y="30" width="52" height="34" rx="6" fill="#2B2117" />
        <rect x="-22" y="34" width="44" height="24" rx="3" fill="#9fd0ff" />
      </g>
      <text x="310" y="270" textAnchor="middle" fontWeight="700" fontSize="16">
        開発者
      </text>

      {/* コード(モニター) */}
      <g transform="translate(520,190)">
        <rect x="-46" y="-30" width="92" height="66" rx="10" fill="#2B2117" />
        <rect x="-38" y="-22" width="76" height="50" rx="4" fill="#FBE4D8" />
        <text x="0" y="0" textAnchor="middle" className="font-mono" fontSize="12" fill="#E8672E">
          Order
        </text>
        <text x="0" y="14" textAnchor="middle" className="font-mono" fontSize="10" fill="#5A4A3D">
          {"{ ... }"}
        </text>
        <rect x="-10" y="36" width="20" height="14" fill="#8a7a6a" />
        <rect x="-30" y="50" width="60" height="8" rx="4" fill="#8a7a6a" />
      </g>
      <text x="520" y="270" textAnchor="middle" fontWeight="700" fontSize="16">
        コード
      </text>
    </svg>
  );
}
