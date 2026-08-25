import Link from "next/link";
import { ArrowLeft, ArrowRight, LayoutList } from "lucide-react";
import { Button } from "@/components/ui/button";
import { getAdjacentPages, pageHref } from "@/lib/pages-data";

/**
 * 概念ページ下部の「前へ / 目次へ / 次へ」ナビゲーション。
 * ページ順はlib/pages-data.tsのPAGE_ORDERを唯一の情報源とする。
 */
export function PageNav({ currentSlug }: { currentSlug: string }) {
  const { prev, next } = getAdjacentPages(currentSlug);

  return (
    <nav
      aria-label="ページ送り"
      className="border-b border-border/70 bg-card px-6 py-5"
    >
      <div className="mx-auto flex max-w-3xl flex-wrap items-center justify-between gap-3.5">
        {prev ? (
          <Button
            variant="outline"
            size="lg"
            className="order-1 h-auto rounded-full px-4 py-2.5"
            render={<Link href={pageHref(prev)} />}
          >
            <ArrowLeft aria-hidden="true" />
            前へ: {prev.title}
          </Button>
        ) : (
          <span className="order-1" />
        )}

        <Button
          size="lg"
          className="order-3 h-auto rounded-full px-4 py-2.5 sm:order-2 sm:mx-auto"
          render={<Link href="/" />}
        >
          <LayoutList aria-hidden="true" />
          目次へ
        </Button>

        {next ? (
          <Button
            variant="outline"
            size="lg"
            className="order-2 h-auto rounded-full px-4 py-2.5 sm:order-3"
            render={<Link href={pageHref(next)} />}
          >
            次へ: {next.title}
            <ArrowRight aria-hidden="true" />
          </Button>
        ) : (
          <span className="order-2 sm:order-3" />
        )}
      </div>
    </nav>
  );
}
