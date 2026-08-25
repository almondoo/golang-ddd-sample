/**
 * トップページのヒーロー図解(旧静的サイトのindex.htmlより移植)。
 * 色はオリジナルの絵本パレットをそのまま保持している(サイト全体のUIトークンとは別の、挿絵専用の配色)。
 */
export function HeroDiagram() {
  return (
    <svg viewBox="0 0 640 380" role="img" aria-labelledby="hero-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="hero-title">
        現場で使われる言葉が、DDDを通ってそのままきれいなコードになる図
      </title>

      {/* 左: 現場でとびかう言葉(バラバラなタグ) */}
      <g transform="translate(70,80) rotate(-8)">
        <rect x="-46" y="-18" width="92" height="36" rx="18" fill="#C4577B" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="700" fontSize="17">
          注文
        </text>
      </g>
      <g transform="translate(190,55) rotate(6)">
        <rect x="-46" y="-18" width="92" height="36" rx="18" fill="#5B8F4E" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="700" fontSize="17">
          在庫
        </text>
      </g>
      <g transform="translate(60,175) rotate(9)">
        <rect x="-46" y="-18" width="92" height="36" rx="18" fill="#7A5CA8" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="700" fontSize="17">
          配送
        </text>
      </g>
      <g transform="translate(195,165) rotate(-7)">
        <rect x="-46" y="-18" width="92" height="36" rx="18" fill="#3E6FB0" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="700" fontSize="17">
          顧客
        </text>
      </g>
      <g transform="translate(85,260) rotate(-5)">
        <rect x="-46" y="-18" width="92" height="36" rx="18" fill="#B87A1E" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="700" fontSize="17">
          商品
        </text>
      </g>
      <g transform="translate(210,255) rotate(8)">
        <rect x="-52" y="-18" width="104" height="36" rx="18" fill="#2E8F86" />
        <text x="0" y="6" textAnchor="middle" fill="#fff" fontWeight="700" fontSize="17">
          カート
        </text>
      </g>

      {/* 中央の矢印とDDDバッジ */}
      <polygon points="270,180 355,180 355,158 400,197 355,236 355,214 270,214" fill="#E8672E" />
      <circle cx="325" cy="130" r="34" fill="#E8672E" />
      <text x="325" y="137" textAnchor="middle" fill="#fff" fontWeight="800" fontSize="19">
        DDD
      </text>

      {/* 右: きれいに整理されたコード */}
      <rect x="410" y="45" width="210" height="295" rx="24" fill="#FFFFFF" stroke="#EAD9C2" strokeWidth="2" />
      <text x="515" y="80" textAnchor="middle" fontWeight="700" fontSize="15" fill="#5A4A3D">
        きれいなコード
      </text>

      <rect x="425" y="98" width="180" height="42" rx="10" fill="#FBE4D8" />
      <text x="435" y="124" className="font-mono" fontSize="15" fill="#2B2117">
        type Order
      </text>

      <rect x="425" y="150" width="180" height="42" rx="10" fill="#DCE7F5" />
      <text x="435" y="176" className="font-mono" fontSize="15" fill="#2B2117">
        type Cart
      </text>

      <rect x="425" y="202" width="180" height="42" rx="10" fill="#F6E7CC" />
      <text x="435" y="228" className="font-mono" fontSize="15" fill="#2B2117">
        type Product
      </text>

      <rect x="425" y="254" width="180" height="42" rx="10" fill="#F7E1E7" />
      <text x="435" y="280" className="font-mono" fontSize="15" fill="#2B2117">
        type Customer
      </text>
    </svg>
  );
}
