/**
 * CQRSページの図解(新規作成)。
 * command(書き込み)はドメインを経由し、query(読み取り)はドメインを迂回して
 * どちらも同じDBへたどり着く、という非対称な2経路を示す。
 */
export function CqrsDiagram() {
  return (
    <svg viewBox="0 0 560 440" role="img" aria-labelledby="cqrs-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="cqrs-title">
        commandはUseCase→Repository→Product集約(検証)→GORM実装という経路でドメインを通り、queryはUseCase→QueryService→GORM実装(DTO直行)という経路でドメインを迂回し、両方とも同じPostgreSQLへ書き込み・問い合わせすることを示す図
      </title>
      <defs>
        <marker id="cqrs-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--foreground)" />
        </marker>
      </defs>

      <text x="150" y="22" textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--primary)">command(書き込み)</text>
      <text x="410" y="22" textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--chart-2)">query(読み取り)</text>

      {/* command column */}
      <rect x="60" y="35" width="180" height="50" rx="14" fill="var(--primary)" />
      <text x="150" y="65" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--primary-foreground)">RegisterProductUseCase</text>
      <line x1="150" y1="85" x2="150" y2="115" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#cqrs-arrow)" />

      <rect x="60" y="115" width="180" height="50" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
      <text x="150" y="145" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">catalog.Repository(IF)</text>
      <line x1="150" y1="165" x2="150" y2="195" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#cqrs-arrow)" />

      <rect x="60" y="195" width="180" height="50" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
      <text x="150" y="225" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--foreground)">Product集約(検証)</text>
      <line x1="150" y1="245" x2="150" y2="275" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#cqrs-arrow)" />

      <rect x="60" y="275" width="180" height="50" rx="14" fill="var(--primary)" />
      <text x="150" y="305" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--primary-foreground)">GORM実装(INSERT/UPDATE)</text>

      {/* query column */}
      <rect x="320" y="35" width="180" height="50" rx="14" fill="var(--chart-2)" />
      <text x="410" y="65" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">ListProductsUseCase</text>
      <line x1="410" y1="85" x2="410" y2="115" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#cqrs-arrow)" />

      <rect x="320" y="115" width="180" height="50" rx="14" fill="var(--card)" stroke="var(--chart-2)" strokeWidth="2.5" />
      <text x="410" y="145" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">ProductQueryService(IF)</text>
      <line x1="410" y1="165" x2="410" y2="275" stroke="var(--foreground)" strokeWidth="3" strokeDasharray="7 6" markerEnd="url(#cqrs-arrow)" />
      <text x="500" y="220" textAnchor="middle" fontSize="11" fill="var(--muted-foreground)">ドメインを迂回</text>

      <rect x="320" y="275" width="180" height="50" rx="14" fill="var(--chart-2)" />
      <text x="410" y="298" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--foreground)">GORM実装</text>
      <text x="410" y="315" textAnchor="middle" fontSize="10" fill="var(--foreground)" opacity="0.85">(DTOに直行)</text>

      {/* DB */}
      <line x1="150" y1="325" x2="230" y2="365" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#cqrs-arrow)" />
      <line x1="410" y1="325" x2="330" y2="365" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#cqrs-arrow)" />
      <ellipse cx="280" cy="385" rx="110" ry="18" fill="var(--accent)" />
      <rect x="170" y="385" width="220" height="30" fill="var(--accent)" />
      <ellipse cx="280" cy="415" rx="110" ry="18" fill="var(--accent)" />
      <text x="280" y="404" textAnchor="middle" fontWeight="800" fontSize="13" fill="var(--accent-foreground)">
        PostgreSQL(同一DB)
      </text>
    </svg>
  );
}
