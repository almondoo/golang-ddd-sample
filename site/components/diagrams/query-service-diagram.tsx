/**
 * クエリサービスページの図解(新規作成)。
 * 「書く窓口」(UseCase→集約→Repository)と「見る窓口」(UseCase→QueryService→DTO直行)の
 * 2つの窓口を並べ、見る窓口の脇にドメイン箱を薄く・バツ印つきで置いて
 * 「素通りされる」ことを視覚化する。
 */
export function QueryServiceDiagram() {
  return (
    <svg viewBox="0 0 600 460" role="img" aria-labelledby="qs-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="qs-title">
        書く窓口はUseCaseから集約(ドメイン)を経てRepositoryへ進みDBに書き込むのに対し、見る窓口はUseCaseからQueryServiceへ進みドメイン箱を素通りしてDTOを直接DBから組み立てることを示す図
      </title>
      <defs>
        <marker id="qs-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--foreground)" />
        </marker>
      </defs>

      <text x="150" y="26" textAnchor="middle" fontWeight="800" fontSize="15" fill="var(--primary)">書く窓口</text>
      <text x="450" y="26" textAnchor="middle" fontWeight="800" fontSize="15" fill="var(--chart-2)">見る窓口</text>

      {/* 書く窓口カウンター */}
      <rect x="50" y="40" width="200" height="52" rx="14" fill="var(--primary)" />
      <text x="150" y="71" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--primary-foreground)">ChangePriceUseCase</text>
      <line x1="150" y1="92" x2="150" y2="122" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />

      <rect x="50" y="122" width="200" height="52" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
      <text x="150" y="153" textAnchor="middle" fontWeight="700" fontSize="13" fill="var(--foreground)">集約(ドメイン)</text>
      <text x="150" y="168" textAnchor="middle" className="font-mono" fontSize="10" fill="var(--foreground)" opacity="0.85">Product / Cart</text>
      <line x1="150" y1="174" x2="150" y2="204" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />

      <rect x="50" y="204" width="200" height="52" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
      <text x="150" y="235" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">catalog.Repository(IF)</text>
      <line x1="150" y1="256" x2="150" y2="286" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />

      <rect x="50" y="286" width="200" height="52" rx="14" fill="var(--primary)" />
      <text x="150" y="317" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--primary-foreground)">正式ルートを通って保存</text>

      {/* 見る窓口カウンター */}
      <rect x="350" y="40" width="200" height="52" rx="14" fill="var(--chart-2)" />
      <text x="450" y="71" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">ListProductsUseCase</text>
      <line x1="450" y1="92" x2="450" y2="122" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />

      <rect x="350" y="122" width="200" height="52" rx="14" fill="var(--card)" stroke="var(--chart-2)" strokeWidth="2.5" />
      <text x="450" y="148" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="10.5" fill="var(--foreground)">ProductQueryService(IF)</text>
      <text x="450" y="163" textAnchor="middle" className="font-mono" fontSize="10" fill="var(--foreground)" opacity="0.8">CartQueryService(IF)</text>

      {/* 素通りされるドメイン箱(薄く・バツ印) */}
      <g opacity="0.5">
        <rect x="245" y="128" width="90" height="40" rx="10" fill="none" stroke="var(--muted-foreground)" strokeWidth="2.5" strokeDasharray="6 5" />
        <text x="290" y="152" textAnchor="middle" fontSize="11" fontWeight="700" fill="var(--muted-foreground)">ドメイン</text>
        <line x1="252" y1="135" x2="328" y2="161" stroke="var(--destructive)" strokeWidth="3" strokeLinecap="round" />
        <line x1="328" y1="135" x2="252" y2="161" stroke="var(--destructive)" strokeWidth="3" strokeLinecap="round" />
      </g>
      <text x="290" y="115" textAnchor="middle" fontSize="11" fill="var(--muted-foreground)">素通り</text>

      <line x1="450" y1="174" x2="450" y2="204" stroke="var(--foreground)" strokeWidth="3" strokeDasharray="7 6" markerEnd="url(#qs-arrow)" />

      <rect x="350" y="204" width="200" height="52" rx="14" fill="var(--card)" stroke="var(--chart-2)" strokeWidth="2.5" />
      <text x="450" y="227" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="11" fill="var(--foreground)">persistence.*Query</text>
      <text x="450" y="243" textAnchor="middle" fontSize="10" fill="var(--foreground)" opacity="0.85">DTOに直行</text>
      <line x1="450" y1="256" x2="450" y2="286" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />

      <rect x="350" y="286" width="200" height="52" rx="14" fill="var(--chart-2)" />
      <text x="450" y="317" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--foreground)">DTOを直接返す</text>

      {/* DB */}
      <line x1="150" y1="338" x2="245" y2="378" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />
      <line x1="450" y1="338" x2="355" y2="378" stroke="var(--foreground)" strokeWidth="3" markerEnd="url(#qs-arrow)" />
      <ellipse cx="300" cy="398" rx="120" ry="18" fill="var(--accent)" />
      <rect x="180" y="398" width="240" height="30" fill="var(--accent)" />
      <ellipse cx="300" cy="428" rx="120" ry="18" fill="var(--accent)" />
      <text x="300" y="417" textAnchor="middle" fontWeight="800" fontSize="13" fill="var(--accent-foreground)">
        PostgreSQL(同一DB)
      </text>
    </svg>
  );
}
