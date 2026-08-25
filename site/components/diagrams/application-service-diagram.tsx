/**
 * アプリケーションサービスページの図解(新規作成)。
 * 段取り係(usecase)が、1トランザクションの中で6つの手順を
 * 順番に指示していく様子を示す。
 */
export function ApplicationServiceDiagram() {
  const steps = [
    { label: "顧客の実在確認", detail: "customerRepo.FindByID" },
    { label: "カートの読み取り", detail: "cartRepo.FindByCustomerID" },
    { label: "価格取得+在庫引当", detail: "catalogRepo / inventoryRepo" },
    { label: "Order集約の生成", detail: "domainorder.NewOrder" },
    { label: "クーポン適用(あれば)", detail: "couponRepo → o.ApplyDiscount" },
    { label: "保存してカートを空に", detail: "orderRepo.Save / cart.Clear" },
  ];
  const stepTop = (i: number) => 110 + i * 78;

  return (
    <svg viewBox="0 0 520 630" role="img" aria-labelledby="as-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="as-title">
        段取り係(PlaceOrderUseCase)が、1トランザクションの枠の中で6つの手順を順番に指示していく図
      </title>
      <defs>
        <marker id="as-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--border)" />
        </marker>
      </defs>

      {/* トランザクション境界 */}
      <rect x="20" y="70" width="480" height="540" rx="28" fill="none" stroke="var(--primary)" strokeWidth="3" strokeDasharray="10 8" />
      <rect x="150" y="40" width="220" height="40" rx="20" fill="var(--primary)" />
      <text x="260" y="65" textAnchor="middle" className="font-mono" fontWeight="800" fontSize="13" fill="var(--primary-foreground)">
        tx.Manager.Do(1トランザクション)
      </text>

      {/* 段取り係(人物) */}
      <g transform="translate(70,102)">
        <circle cx="0" cy="0" r="24" fill="var(--accent)" />
        <circle cx="-7" cy="-3" r="3" fill="var(--accent-foreground)" />
        <circle cx="7" cy="-3" r="3" fill="var(--accent-foreground)" />
        <path d="M -8 8 Q 0 15 8 8" stroke="var(--accent-foreground)" strokeWidth="2.5" fill="none" strokeLinecap="round" />
      </g>
      <text x="70" y="150" textAnchor="middle" fontWeight="700" fontSize="12" fill="var(--foreground)">
        段取り係
      </text>
      <text x="70" y="165" textAnchor="middle" className="font-mono" fontSize="10" fill="var(--muted-foreground)">
        usecase
      </text>

      {steps.map((s, i) => {
        const top = stepTop(i);
        return (
          <g key={s.label}>
            {i > 0 ? (
              <line x1="270" y1={stepTop(i - 1) + 60} x2="270" y2={top} stroke="var(--border)" strokeWidth="3" markerEnd="url(#as-arrow)" />
            ) : null}
            <circle cx="140" cy={top + 30} r="16" fill="var(--primary)" />
            <text x="140" y={top + 35} textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--primary-foreground)">
              {i + 1}
            </text>
            <rect x="170" y={top} width="300" height="60" rx="16" fill="var(--card)" stroke="var(--border)" strokeWidth="2" />
            <text x="190" y={top + 26} fontWeight="700" fontSize="14" fill="var(--foreground)">
              {s.label}
            </text>
            <text x="190" y={top + 46} className="font-mono" fontSize="11" fill="var(--muted-foreground)">
              {s.detail}
            </text>
          </g>
        );
      })}
    </svg>
  );
}
