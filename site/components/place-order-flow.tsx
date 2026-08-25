import { Reveal } from "@/components/reveal";

type Step = {
  title: string;
  detail: string;
};

const STEPS: Step[] = [
  {
    title: "HTTP",
    detail: "POST /orders リクエストが届く",
  },
  {
    title: "Controller",
    detail: "OrderController.handlePlaceOrder がJSONを読み取る",
  },
  {
    title: "UseCase開始(1トランザクション)",
    detail: "PlaceOrderUseCase.Execute が txManager.Do の中で処理をまとめる",
  },
  {
    title: "カート読込",
    detail: "cartRepo.FindByCustomerID でカートの中身を取得する",
  },
  {
    title: "在庫確保",
    detail: "stock.Reserve で数量を仮引当てする(実在庫はまだ減らさない)",
  },
  {
    title: "クーポン適用",
    detail: "指定があれば coupon.Use → Order.ApplyDiscount で割引を反映",
  },
  {
    title: "注文保存",
    detail: "orderRepo.Save → カートを空にして保存。201 Createdを返す",
  },
];

type ErrorRow = {
  situation: string;
  example: string;
  status: string;
};

const ERROR_ROWS: ErrorRow[] = [
  {
    situation: "shared.ErrNotFound",
    example: "顧客が見つからない・カートが未作成 など",
    status: "404 Not Found",
  },
  {
    situation: "shared.ErrConflict",
    example: "楽観ロックの競合(同時更新)",
    status: "409 Conflict",
  },
  {
    situation: "NewDomainRuleError",
    example: "カートが空・在庫不足・クーポン無効 など",
    status: "422 Unprocessable Entity",
  },
];

/** 「注文確定の旅」ページ専用: 番号付きステップフローとエラー対応表 */
export function PlaceOrderFlow() {
  return (
    <div className="space-y-14">
      <ol className="relative mx-auto max-w-xl border-l-2 border-border pl-8">
        {STEPS.map((step, index) => (
          <li key={step.title} className="relative pb-10 last:pb-0">
            <Reveal margin="-5%">
              <span
                className="absolute top-0 -left-[calc(2rem+13px)] flex size-8 items-center justify-center rounded-full bg-primary text-sm font-extrabold text-primary-foreground shadow-clay"
                aria-hidden="true"
              >
                {index + 1}
              </span>
              <p className="font-mono text-base font-bold text-foreground">{step.title}</p>
              <p className="mt-1 text-sm text-muted-foreground">{step.detail}</p>
            </Reveal>
          </li>
        ))}
      </ol>

      <Reveal>
        <div className="mx-auto max-w-xl overflow-x-auto rounded-3xl border border-border bg-card shadow-clay">
          <table className="w-full border-collapse text-left text-sm">
            <caption className="px-6 pt-5 text-left text-base font-extrabold text-foreground">
              エラー → HTTPステータス対応
            </caption>
            <thead>
              <tr className="border-b border-border text-xs text-muted-foreground">
                <th scope="col" className="px-6 py-3 font-bold">
                  状況
                </th>
                <th scope="col" className="px-6 py-3 font-bold">
                  具体例
                </th>
                <th scope="col" className="px-6 py-3 font-bold">
                  HTTPステータス
                </th>
              </tr>
            </thead>
            <tbody>
              {ERROR_ROWS.map((row) => (
                <tr key={row.situation} className="border-b border-border last:border-b-0">
                  <td className="px-6 py-3 font-mono text-foreground">{row.situation}</td>
                  <td className="px-6 py-3 text-muted-foreground">{row.example}</td>
                  <td className="px-6 py-3 font-bold text-accent">{row.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Reveal>
    </div>
  );
}
