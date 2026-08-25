/**
 * 仕様(Specification)ページの図解(新規作成)。
 * 教科書的な一般像として、判定条件をオブジェクト化したSpecificationが
 * 全件から条件に合うものだけを絞り込む様子を示す(本リポジトリでの実例はない)。
 */
export function SpecificationDiagram() {
  return (
    <svg viewBox="0 0 560 380" role="img" aria-labelledby="spec-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="spec-title">
        「在庫あり」and「特定カテゴリ」という条件を1つのSpecificationオブジェクトにまとめ、全商品から条件に合うものだけを絞り込む図(教科書的な一般像)
      </title>

      <rect x="10" y="10" width="540" height="360" rx="28" fill="var(--muted)" />

      {/* 左: 全商品 */}
      <rect x="35" y="60" width="140" height="260" rx="20" fill="var(--card)" stroke="var(--border)" strokeWidth="2" />
      {[0, 1, 2, 3, 4].map((i) => (
        <rect key={i} x="55" y={85 + i * 45} width="100" height="30" rx="8" fill={i % 2 === 0 ? "var(--chart-2)" : "var(--chart-4)"} />
      ))}
      <text x="105" y="345" textAnchor="middle" fontWeight="700" fontSize="13" fill="var(--foreground)">全商品</text>

      {/* 条件チップ */}
      <rect x="215" y="55" width="130" height="34" rx="17" fill="var(--card)" stroke="var(--primary)" strokeWidth="2" />
      <text x="280" y="77" textAnchor="middle" fontSize="12" fontWeight="700" fill="var(--foreground)">在庫あり</text>
      <rect x="215" y="98" width="130" height="34" rx="17" fill="var(--card)" stroke="var(--primary)" strokeWidth="2" />
      <text x="280" y="120" textAnchor="middle" fontSize="12" fontWeight="700" fill="var(--foreground)">特定カテゴリ</text>
      <text x="280" y="145" textAnchor="middle" fontSize="12" fontWeight="800" fill="var(--primary)">AND / OR / NOT で合成</text>

      {/* Specification(じょうご) */}
      <polygon points="230,165 330,165 300,230 260,230" fill="var(--accent)" />
      <rect x="260" y="230" width="40" height="35" fill="var(--accent)" />
      <text x="280" y="195" textAnchor="middle" fontWeight="800" fontSize="12" fill="var(--accent-foreground)">Specification</text>

      {/* 矢印: 全商品 -> じょうご -> 結果 */}
      <line x1="175" y1="190" x2="228" y2="190" stroke="var(--foreground)" strokeWidth="4" strokeLinecap="round" />
      <polygon points="228,190 210,182 210,198" fill="var(--foreground)" />
      <path d="M 280 265 C 280 300, 340 300, 385 300" stroke="var(--foreground)" strokeWidth="4" fill="none" strokeLinecap="round" />
      <polygon points="385,300 368,292 368,308" fill="var(--foreground)" />

      {/* 右: 絞り込み結果 */}
      <rect x="390" y="60" width="140" height="260" rx="20" fill="var(--card)" stroke="var(--border)" strokeWidth="2" />
      <rect x="410" y="120" width="100" height="30" rx="8" fill="var(--chart-2)" />
      <rect x="410" y="165" width="100" height="30" rx="8" fill="var(--chart-2)" />
      <text x="460" y="345" textAnchor="middle" fontWeight="700" fontSize="13" fill="var(--foreground)">絞り込み結果</text>
    </svg>
  );
}
