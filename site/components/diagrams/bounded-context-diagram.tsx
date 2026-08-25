/** 境界づけられたコンテキストページの図解(旧静的サイトより移植、色は原本のまま) */
export function BoundedContextDiagram() {
  return (
    <svg viewBox="0 0 700 500" role="img" aria-labelledby="bc-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="bc-title">
        お店の中に7つの部屋(catalog・cart・order・customer・inventory・shipping・coupon)が並ぶ見取り図
      </title>

      <rect x="20" y="60" width="660" height="380" rx="28" fill="#FFFFFF" stroke="#EAD9C2" strokeWidth="2" />
      <text x="350" y="38" textAnchor="middle" fontWeight="800" fontSize="20">
        お店(このリポジトリ)
      </text>

      {/* 1行目: catalog / cart / order */}
      <g>
        <rect x="40" y="85" width="200" height="150" rx="18" fill="#F6E7CC" stroke="#B87A1E" strokeWidth="2" />
        <text x="140" y="120" textAnchor="middle" fontWeight="800" fontSize="18">
          商品カタログ
        </text>
        <text x="140" y="142" textAnchor="middle" className="font-mono" fontSize="13" fill="#5A4A3D">
          catalog
        </text>
        <rect x="80" y="185" width="120" height="30" rx="15" fill="#fff" stroke="#B87A1E" strokeWidth="1.5" />
        <text x="140" y="205" textAnchor="middle" className="font-mono" fontSize="12" fill="#B87A1E">
          商品ID
        </text>
      </g>

      <g>
        <rect x="255" y="85" width="200" height="150" rx="18" fill="#DCE7F5" stroke="#3E6FB0" strokeWidth="2" />
        <text x="355" y="120" textAnchor="middle" fontWeight="800" fontSize="18">
          カート
        </text>
        <text x="355" y="142" textAnchor="middle" className="font-mono" fontSize="13" fill="#5A4A3D">
          cart
        </text>
        <rect x="295" y="185" width="120" height="30" rx="15" fill="#fff" stroke="#3E6FB0" strokeWidth="1.5" />
        <text x="355" y="205" textAnchor="middle" className="font-mono" fontSize="12" fill="#3E6FB0">
          商品ID
        </text>
      </g>

      <g>
        <rect x="470" y="85" width="190" height="150" rx="18" fill="#FBE4D8" stroke="#E8672E" strokeWidth="2" />
        <text x="565" y="120" textAnchor="middle" fontWeight="800" fontSize="18">
          注文
        </text>
        <text x="565" y="142" textAnchor="middle" className="font-mono" fontSize="13" fill="#5A4A3D">
          order
        </text>
      </g>

      {/* 2行目: customer / inventory / shipping / coupon */}
      <g>
        <rect x="40" y="255" width="150" height="150" rx="18" fill="#F7E1E7" stroke="#C4577B" strokeWidth="2" />
        <text x="115" y="295" textAnchor="middle" fontWeight="800" fontSize="17">
          顧客
        </text>
        <text x="115" y="317" textAnchor="middle" className="font-mono" fontSize="12" fill="#5A4A3D">
          customer
        </text>
      </g>

      <g>
        <rect x="205" y="255" width="150" height="150" rx="18" fill="#E1EFDC" stroke="#5B8F4E" strokeWidth="2" />
        <text x="280" y="290" textAnchor="middle" fontWeight="800" fontSize="17">
          在庫
        </text>
        <text x="280" y="312" textAnchor="middle" className="font-mono" fontSize="12" fill="#5A4A3D">
          inventory
        </text>
        <rect x="230" y="345" width="100" height="30" rx="15" fill="#fff" stroke="#5B8F4E" strokeWidth="1.5" />
        <text x="280" y="365" textAnchor="middle" className="font-mono" fontSize="12" fill="#5B8F4E">
          商品ID
        </text>
      </g>

      <g>
        <rect x="370" y="255" width="150" height="150" rx="18" fill="#E9E1F3" stroke="#7A5CA8" strokeWidth="2" />
        <text x="445" y="295" textAnchor="middle" fontWeight="800" fontSize="17">
          配送
        </text>
        <text x="445" y="317" textAnchor="middle" className="font-mono" fontSize="12" fill="#5A4A3D">
          shipping
        </text>
      </g>

      <g>
        <rect x="535" y="255" width="125" height="150" rx="18" fill="#DBF0EE" stroke="#2E8F86" strokeWidth="2" />
        <text x="597" y="295" textAnchor="middle" fontWeight="800" fontSize="17">
          クーポン
        </text>
        <text x="597" y="317" textAnchor="middle" className="font-mono" fontSize="12" fill="#5A4A3D">
          coupon
        </text>
      </g>
    </svg>
  );
}
