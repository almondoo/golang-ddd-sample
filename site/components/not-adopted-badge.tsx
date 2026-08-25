import { Ban } from "lucide-react";

/**
 * 「このサンプルでは不採用」を示すピル型バッジ。
 * domain-service / domain-event / specification の3ページで使う。
 * 正直に「使っていない」ことを伝える、否定ではなく学びの一部としてのトーン。
 */
export function NotAdoptedBadge() {
  return (
    <span className="inline-flex items-center gap-2 rounded-full border-2 border-muted-foreground/40 bg-muted px-4 py-2 text-sm font-bold text-muted-foreground">
      <Ban className="size-4 shrink-0" aria-hidden="true" />
      このサンプルでは使っていません
    </span>
  );
}
