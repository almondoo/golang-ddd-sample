/**
 * ドメインサービスページの図解(新規作成)。
 * 「1つの集約で判断できるか」「同一コンテキストで複数集約をまたぐか」という
 * 判定フローと、本リポジトリでは後者(ドメインサービス)に実例が無いことを示す。
 */
export function DomainServiceDiagram() {
  return (
    <svg viewBox="0 0 660 520" role="img" aria-labelledby="ds-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="ds-title">
        ロジックの置き場所を判定するフローチャート。1つの集約で判断できれば集約のメソッドへ、同一コンテキストで複数集約をまたぐならドメインサービスへ(本リポジトリに実例なし)、コンテキストをまたぐならアプリケーションサービスへ進む
      </title>
      <defs>
        <marker id="ds-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--foreground)" />
        </marker>
      </defs>

      {/* Q1: 1つの集約で判断できるか */}
      <rect x="110" y="20" width="340" height="60" rx="20" fill="var(--muted)" stroke="var(--border)" strokeWidth="2" />
      <text x="280" y="45" textAnchor="middle" fontWeight="700" fontSize="15" fill="var(--foreground)">
        1つの集約(ルート+子)だけで
      </text>
      <text x="280" y="65" textAnchor="middle" fontWeight="700" fontSize="15" fill="var(--foreground)">
        判断できる?
      </text>

      {/* Yes -> 集約のメソッド */}
      <line x1="450" y1="50" x2="432" y2="140" stroke="var(--primary)" strokeWidth="4" markerEnd="url(#ds-arrow)" />
      <text x="450" y="105" fontSize="13" fontWeight="800" fill="var(--primary)">Yes</text>
      <rect x="430" y="110" width="200" height="70" rx="16" fill="var(--primary)" />
      <text x="530" y="140" textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--primary-foreground)">
        集約のメソッド
      </text>
      <text x="530" y="160" textAnchor="middle" className="font-mono" fontSize="9" fill="var(--primary-foreground)" opacity="0.9">
        Customer.ChangeDefaultAddress
      </text>

      {/* No -> Q2 */}
      <line x1="280" y1="80" x2="280" y2="210" stroke="var(--foreground)" strokeWidth="4" markerEnd="url(#ds-arrow)" />
      <text x="300" y="150" fontSize="13" fontWeight="800" fill="var(--foreground)">No</text>

      {/* Q2 */}
      <rect x="100" y="210" width="360" height="60" rx="20" fill="var(--muted)" stroke="var(--border)" strokeWidth="2" />
      <text x="280" y="235" textAnchor="middle" fontWeight="700" fontSize="15" fill="var(--foreground)">
        同一コンテキスト内で複数集約を
      </text>
      <text x="280" y="255" textAnchor="middle" fontWeight="700" fontSize="15" fill="var(--foreground)">
        またぐ?
      </text>

      {/* Yes -> ドメインサービス(不採用) */}
      <line x1="460" y1="240" x2="432" y2="345" stroke="var(--muted-foreground)" strokeWidth="4" markerEnd="url(#ds-arrow)" />
      <text x="465" y="290" fontSize="13" fontWeight="800" fill="var(--muted-foreground)">Yes</text>
      <rect x="430" y="310" width="200" height="90" rx="16" fill="var(--muted)" stroke="var(--muted-foreground)" strokeWidth="2.5" strokeDasharray="7 6" />
      <text x="530" y="340" textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--muted-foreground)">
        ドメインサービス
      </text>
      <text x="530" y="360" textAnchor="middle" fontSize="11" fill="var(--muted-foreground)">
        本リポジトリに
      </text>
      <text x="530" y="376" textAnchor="middle" fontSize="11" fill="var(--muted-foreground)">
        実例なし
      </text>
      {/* 禁止アイコン(実例なしを示す) */}
      <circle cx="610" cy="322" r="13" fill="none" stroke="var(--muted-foreground)" strokeWidth="3" />
      <line x1="601" y1="313" x2="619" y2="331" stroke="var(--muted-foreground)" strokeWidth="3" strokeLinecap="round" />

      {/* No -> アプリケーションサービス */}
      <line x1="280" y1="270" x2="280" y2="400" stroke="var(--accent)" strokeWidth="4" markerEnd="url(#ds-arrow)" />
      <text x="265" y="340" textAnchor="end" fontSize="13" fontWeight="800" fill="var(--accent)">No(コンテキストをまたぐ)</text>
      <rect x="130" y="400" width="300" height="80" rx="20" fill="var(--accent)" />
      <text x="280" y="435" textAnchor="middle" fontWeight="800" fontSize="15" fill="var(--accent-foreground)">
        アプリケーションサービス
      </text>
      <text x="280" y="458" textAnchor="middle" className="font-mono" fontSize="12" fill="var(--accent-foreground)" opacity="0.9">
        PlaceOrderUseCase
      </text>
    </svg>
  );
}
