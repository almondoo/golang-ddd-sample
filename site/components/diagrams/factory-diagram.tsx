/**
 * ファクトリページの図解(新規作成)。
 * 「新規生成(New*、検証あり)」と「再構築(Reconstruct*、検証なし・リポジトリ専用)」の
 * 2つの経路が分かれていることを示す。
 */
export function FactoryDiagram() {
  return (
    <svg viewBox="0 0 600 400" role="img" aria-labelledby="factory-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="factory-title">
        新規生成はNewProductが検証してから*Productを作り、再構築はReconstructProductがDBの検証済みデータをそのまま復元する、2つの経路を示す図
      </title>
      <defs>
        <marker id="factory-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--foreground)" />
        </marker>
      </defs>

      <rect x="10" y="10" width="580" height="380" rx="28" fill="var(--muted)" />

      {/* 新規生成レーン */}
      <text x="35" y="45" fontWeight="800" fontSize="15" fill="var(--foreground)">新規生成</text>
      <rect x="30" y="60" width="150" height="50" rx="14" fill="var(--card)" stroke="var(--border)" strokeWidth="2" />
      <text x="105" y="90" textAnchor="middle" fontSize="12" fill="var(--foreground)">name / desc / price</text>
      <line x1="180" y1="85" x2="205" y2="85" stroke="var(--foreground)" strokeWidth="4" markerEnd="url(#factory-arrow)" />

      <rect x="210" y="60" width="230" height="70" rx="18" fill="var(--primary)" />
      <text x="325" y="90" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="15" fill="var(--primary-foreground)">
        NewProduct
      </text>
      <text x="325" y="110" textAnchor="middle" fontSize="12" fill="var(--primary-foreground)" opacity="0.9">
        (検証あり)
      </text>

      <line x1="440" y1="80" x2="462" y2="65" stroke="var(--primary)" strokeWidth="4" markerEnd="url(#factory-arrow)" />
      <rect x="470" y="30" width="110" height="50" rx="14" fill="var(--secondary)" />
      <text x="525" y="60" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--secondary-foreground)">
        *Product(新規)
      </text>

      <line x1="440" y1="110" x2="470" y2="150" stroke="var(--destructive)" strokeWidth="4" markerEnd="url(#factory-arrow)" />
      <circle cx="505" cy="165" r="16" fill="none" stroke="var(--destructive)" strokeWidth="4" />
      <line x1="494" y1="154" x2="516" y2="176" stroke="var(--destructive)" strokeWidth="4" strokeLinecap="round" />
      <text x="505" y="200" textAnchor="middle" fontSize="12" fontWeight="700" fill="var(--destructive)">
        検証NG→エラー
      </text>

      {/* 区切り線 */}
      <line x1="30" y1="225" x2="570" y2="225" stroke="var(--border)" strokeWidth="2" strokeDasharray="6 6" />

      {/* 再構築レーン */}
      <text x="35" y="255" fontWeight="800" fontSize="15" fill="var(--foreground)">再構築(リポジトリ専用)</text>
      <rect x="30" y="270" width="150" height="50" rx="14" fill="var(--card)" stroke="var(--border)" strokeWidth="2" />
      <text x="105" y="300" textAnchor="middle" fontSize="12" fill="var(--foreground)">DB行(検証済み)</text>
      <line x1="180" y1="295" x2="205" y2="295" stroke="var(--foreground)" strokeWidth="4" markerEnd="url(#factory-arrow)" />

      <rect x="210" y="270" width="260" height="70" rx="18" fill="var(--accent)" />
      <text x="340" y="300" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="14" fill="var(--accent-foreground)">
        ReconstructProduct
      </text>
      <text x="340" y="320" textAnchor="middle" fontSize="12" fill="var(--accent-foreground)" opacity="0.9">
        (検証なし)
      </text>

      <line x1="470" y1="305" x2="485" y2="305" stroke="var(--accent)" strokeWidth="4" markerEnd="url(#factory-arrow)" />
      <rect x="495" y="280" width="90" height="50" rx="14" fill="var(--secondary)" />
      <text x="540" y="310" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--secondary-foreground)">
        *Product(復元)
      </text>

      <text x="340" y="375" textAnchor="middle" fontSize="13" fill="var(--foreground)">
        Reconstruct*はRepositoryのFindByID内でのみ呼ばれる
      </text>
    </svg>
  );
}
