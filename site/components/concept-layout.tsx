import type { ReactNode } from "react";
import { Reveal } from "@/components/reveal";
import { PageNav } from "@/components/page-nav";

type ExampleContent = {
  code: ReactNode;
  note: ReactNode;
};

export function ConceptLayout({
  slug,
  eyebrow,
  title,
  lead,
  badge,
  children,
  example,
  summary,
}: {
  /** ルートセグメント(前へ/次への算出に使う) */
  slug: string;
  /** 例: "01 / ユビキタス言語" */
  eyebrow: string;
  title: ReactNode;
  lead: string;
  /** 「このサンプルでは不採用」等、ヒーロー直下に置くバッジ(不採用系ページのみ使用) */
  badge?: ReactNode;
  /** 図解(SVG)と、必要に応じた補足段落 */
  children: ReactNode;
  example?: ExampleContent;
  summary: string;
}) {
  return (
    <>
      <main>
        <section className="border-b border-border/70 bg-gradient-to-b from-background to-secondary/15 px-6 py-16 text-center sm:py-20">
          <div className="mx-auto max-w-3xl">
            <p className="mb-2.5 inline-flex items-center gap-2 text-sm font-bold tracking-wide text-muted-foreground">
              <span className="inline-block h-2.5 w-2.5 shrink-0 rounded-full bg-accent" aria-hidden="true" />
              {eyebrow}
            </p>
            <h1 className="text-4xl leading-tight font-extrabold sm:text-5xl">
              {title}
            </h1>
            <p className="mx-auto mt-5 max-w-xl text-lg font-medium sm:text-xl">
              {lead}
            </p>
            {badge ? <div className="mt-6 flex justify-center">{badge}</div> : null}
          </div>
        </section>

        <section className="px-6 py-16 sm:py-18">
          <div className="mx-auto max-w-3xl">{children}</div>
        </section>

        {example ? (
          <section className="border-t border-border/70 bg-card px-6 py-16 sm:py-18">
            <div className="mx-auto max-w-3xl">
              <h2 className="mb-9 text-center text-2xl font-extrabold sm:text-3xl">
                本リポジトリでの実例
              </h2>
              <Reveal>
                <pre className="mx-auto mb-5 max-w-2xl overflow-x-auto rounded-3xl bg-foreground px-6 py-6 text-sm leading-7 whitespace-pre text-background shadow-clay font-mono">
                  {example.code}
                </pre>
              </Reveal>
              <p className="mx-auto max-w-2xl text-center text-sm text-muted-foreground">
                {example.note}
              </p>
            </div>
          </section>
        ) : null}

        <section className="px-6 py-16 sm:py-18">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="mb-5 text-2xl font-extrabold sm:text-3xl">
              ひとことまとめ
            </h2>
            <p className="inline-block rounded-3xl border-2 border-primary bg-primary/10 px-6 py-4.5 text-lg font-bold text-foreground">
              {summary}
            </p>
          </div>
        </section>
      </main>

      <PageNav currentSlug={slug} />
    </>
  );
}
