/** 値オブジェクトページの図解(旧静的サイトより移植、色は原本のまま) */
export function ValueObjectDiagram() {
  return (
    <svg viewBox="0 0 520 320" role="img" aria-labelledby="vo-title" className="mx-auto block h-auto w-full max-w-full">
      <title id="vo-title">
        2つの「1000円」が、別々のものでも中身が同じなら同じ扱いになることを示す図
      </title>

      <rect x="20" y="20" width="480" height="280" rx="24" fill="#F6E7CC" />

      <rect x="55" y="90" width="130" height="90" rx="16" fill="#fff" stroke="#B87A1E" strokeWidth="3" />
      <text x="120" y="145" textAnchor="middle" fontWeight="800" fontSize="22" fill="#B87A1E">
        ¥1000
      </text>

      <text x="260" y="145" textAnchor="middle" fontSize="30" fontWeight="800" fill="#B87A1E">
        =
      </text>

      <rect x="335" y="90" width="130" height="90" rx="16" fill="#fff" stroke="#B87A1E" strokeWidth="3" />
      <text x="400" y="145" textAnchor="middle" fontWeight="800" fontSize="22" fill="#B87A1E">
        ¥1000
      </text>

      <text x="260" y="230" textAnchor="middle" fontSize="15" fill="#5A4A3D">
        財布のどのお札でも、金額が同じなら同じ価値
      </text>
      <text x="260" y="260" textAnchor="middle" fontSize="15" fill="#5A4A3D">
        IDで区別しない・後から書き換えない
      </text>
    </svg>
  );
}
