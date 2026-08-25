import Link from "next/link";

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-10 flex items-center justify-center border-b border-border/70 bg-background/85 px-6 py-3.5 backdrop-blur-sm">
      <Link
        href="/"
        className="inline-flex items-baseline gap-1.5 text-base font-extrabold tracking-wide text-primary focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50 rounded-md"
      >
        DDD絵本
        <span className="text-sm font-medium text-foreground">
          図解でわかるドメイン駆動設計
        </span>
      </Link>
    </header>
  );
}
