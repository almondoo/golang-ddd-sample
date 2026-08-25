/** オニオンアーキテクチャページの図解(新規作成)。サイトのデザイントークンを色に使う */
export function OnionDiagram() {
  return (
    <svg viewBox="0 0 520 520" role="img" aria-labelledby="onion-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="onion-title">
        中心のdomainを、application・infrastructure/presentationの同心円が取り囲み、依存の矢印はすべて内向きであることを示す図
      </title>

      {/* 外側の輪: presentation / infrastructure */}
      <circle cx="260" cy="260" r="230" fill="var(--secondary)" />
      {/* 中間の輪: application */}
      <circle cx="260" cy="260" r="152" fill="var(--primary)" />
      {/* 中心: domain */}
      <circle cx="260" cy="260" r="76" fill="var(--accent)" />

      {/* 内向きの依存矢印(3方向)。1本目は「アプリケーション」ラベルとの重なりを避けるため35度回転させて配置している */}
      <g transform="rotate(35 260 260)">
        <line x1="260" y1="55" x2="260" y2="165" stroke="var(--card)" strokeWidth="6" strokeLinecap="round" />
        <polygon points="260,180 250,158 270,158" fill="var(--card)" stroke="var(--foreground)" strokeWidth="1.5" />
      </g>
      <g transform="rotate(120 260 260)">
        <line x1="260" y1="55" x2="260" y2="165" stroke="var(--card)" strokeWidth="6" strokeLinecap="round" />
        <polygon points="260,180 250,158 270,158" fill="var(--card)" stroke="var(--foreground)" strokeWidth="1.5" />
      </g>
      <g transform="rotate(240 260 260)">
        <line x1="260" y1="55" x2="260" y2="165" stroke="var(--card)" strokeWidth="6" strokeLinecap="round" />
        <polygon points="260,180 250,158 270,158" fill="var(--card)" stroke="var(--foreground)" strokeWidth="1.5" />
      </g>

      {/* 中心: domain ラベル */}
      <text x="260" y="252" textAnchor="middle" fontWeight="800" fontSize="19" fill="var(--background)">
        ドメイン
      </text>
      <text x="260" y="274" textAnchor="middle" className="font-mono" fontSize="12" fill="var(--background)" opacity="0.9">
        internal/domain
      </text>

      {/* 中間の輪: application ラベル */}
      <text x="260" y="122" textAnchor="middle" fontWeight="800" fontSize="16" fill="var(--background)">
        アプリケーション
      </text>
      <text x="260" y="140" textAnchor="middle" className="font-mono" fontSize="11" fill="var(--background)" opacity="0.9">
        internal/application
      </text>

      {/* 外側の輪ラベル: presentation(上)・infrastructure(下) */}
      <text x="260" y="45" textAnchor="middle" fontWeight="800" fontSize="15" fill="var(--foreground)">
        プレゼンテーション
      </text>
      <text x="260" y="62" textAnchor="middle" className="font-mono" fontSize="11" fill="var(--foreground)">
        internal/presentation
      </text>
      <text x="260" y="462" textAnchor="middle" fontWeight="800" fontSize="15" fill="var(--foreground)">
        インフラストラクチャ
      </text>
      <text x="260" y="479" textAnchor="middle" className="font-mono" fontSize="11" fill="var(--foreground)">
        internal/infrastructure
      </text>
    </svg>
  );
}
