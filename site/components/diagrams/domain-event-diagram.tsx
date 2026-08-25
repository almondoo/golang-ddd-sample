/**
 * ドメインイベントページの図解(新規作成)。
 * 「以前(イベントバス経由、不採用)」と「現在(直接呼び出し)」を
 * 上下に並べて対比させる。
 */
export function DomainEventDiagram() {
  return (
    <svg viewBox="0 0 600 440" role="img" aria-labelledby="de-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="de-title">
        以前はイベントバス経由でCartハンドラが購読していたが、現在は学習コストを下げるためPlaceOrderUseCaseがorderRepoとcartRepoを直接呼び出す、という変遷を示す図
      </title>
      <defs>
        <marker id="de-arrow" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--muted-foreground)" />
        </marker>
        <marker id="de-arrow-live" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--primary)" />
        </marker>
      </defs>

      {/* 以前(不採用) */}
      <text x="300" y="20" textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--muted-foreground)">
        以前(不採用): イベントバス経由
      </text>
      <g opacity="0.55">
        <rect x="20" y="35" width="150" height="55" rx="14" fill="var(--muted)" stroke="var(--muted-foreground)" strokeWidth="2" />
        <text x="95" y="67" textAnchor="middle" fontSize="12" fontWeight="700" fill="var(--muted-foreground)">PlaceOrderUseCase</text>

        <line x1="170" y1="62" x2="220" y2="62" stroke="var(--muted-foreground)" strokeWidth="3" markerEnd="url(#de-arrow)" />
        <rect x="225" y="35" width="140" height="55" rx="14" fill="var(--muted)" stroke="var(--muted-foreground)" strokeWidth="2" />
        <text x="295" y="58" textAnchor="middle" fontSize="12" fontWeight="700" fill="var(--muted-foreground)">イベントバス</text>
        <text x="295" y="74" textAnchor="middle" fontSize="10" fill="var(--muted-foreground)">Publish(OrderPlaced)</text>

        <line x1="365" y1="62" x2="415" y2="62" stroke="var(--muted-foreground)" strokeWidth="3" markerEnd="url(#de-arrow)" />
        <rect x="420" y="35" width="160" height="55" rx="14" fill="var(--muted)" stroke="var(--muted-foreground)" strokeWidth="2" />
        <text x="500" y="58" textAnchor="middle" fontSize="12" fontWeight="700" fill="var(--muted-foreground)">Cartハンドラ</text>
        <text x="500" y="74" textAnchor="middle" fontSize="10" fill="var(--muted-foreground)">購読 → cart.Clear()</text>
      </g>
      {/* 禁止アイコン */}
      <circle cx="580" cy="20" r="14" fill="none" stroke="var(--muted-foreground)" strokeWidth="3" />
      <line x1="570" y1="10" x2="590" y2="30" stroke="var(--muted-foreground)" strokeWidth="3" strokeLinecap="round" />

      {/* 変更の矢印 */}
      <line x1="300" y1="105" x2="300" y2="175" stroke="var(--foreground)" strokeWidth="4" markerEnd="url(#de-arrow-live)" />
      <text x="315" y="130" fontSize="12" fontWeight="700" fill="var(--foreground)">学習コストを下げるため</text>
      <text x="315" y="147" fontSize="12" fontWeight="700" fill="var(--foreground)">直接呼び出しに変更</text>

      {/* 現在 */}
      <text x="300" y="205" textAnchor="middle" fontWeight="800" fontSize="14" fill="var(--foreground)">
        現在: 直接呼び出し
      </text>
      <rect x="180" y="220" width="240" height="55" rx="14" fill="var(--primary)" />
      <text x="300" y="253" textAnchor="middle" fontWeight="800" fontSize="13" fill="var(--primary-foreground)">
        PlaceOrderUseCase
      </text>

      <line x1="230" y1="275" x2="130" y2="330" stroke="var(--primary)" strokeWidth="4" markerEnd="url(#de-arrow-live)" />
      <line x1="370" y1="275" x2="470" y2="330" stroke="var(--primary)" strokeWidth="4" markerEnd="url(#de-arrow-live)" />

      <rect x="30" y="335" width="180" height="55" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
      <text x="120" y="368" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="13" fill="var(--foreground)">
        orderRepo.Save
      </text>

      <rect x="390" y="335" width="180" height="55" rx="14" fill="var(--card)" stroke="var(--primary)" strokeWidth="2.5" />
      <text x="480" y="368" textAnchor="middle" className="font-mono" fontWeight="700" fontSize="13" fill="var(--foreground)">
        cartRepo.Save
      </text>

      <text x="300" y="425" textAnchor="middle" fontSize="12" fill="var(--muted-foreground)">
        同一トランザクション内、order→cartを直接import
      </text>
    </svg>
  );
}
