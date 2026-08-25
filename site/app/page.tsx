import type { Metadata } from "next";
import { Reveal } from "@/components/reveal";
import { TocGrid } from "@/components/toc-grid";
import { HeroDiagram } from "@/components/diagrams/hero-diagram";

export const metadata: Metadata = {
  title: "DDDってなに? - 図解でわかるドメイン駆動設計",
  description:
    "Goで書かれたECサイトのサンプルコードを教材に、ドメイン駆動設計(DDD)の考え方を絵で説明する絵本サイトです。",
};

export default function Home() {
  return (
    <main>
      <section className="border-b border-border/70 bg-gradient-to-b from-background to-secondary/15 px-6 py-20 text-center sm:py-24">
        <div className="mx-auto max-w-3xl">
          <h1 className="text-5xl leading-tight font-extrabold sm:text-6xl">
            DDDって
            <br />
            なに?
          </h1>
          <p className="mt-5 text-lg font-bold text-primary sm:text-xl">
            業務の言葉のまま、ソフトウェアを作る方法。
          </p>
          <p className="mx-auto mt-4 max-w-xl text-sm text-muted-foreground sm:text-base">
            「注文」「在庫」みたいな、現場でふだん使う言葉をそのままコードの名前にする。それだけです。このサイトでは、Goで書かれたECサイトのサンプルコードを教材に、DDDの考え方を1ページ1テーマの絵本で見ていきます。
          </p>

          <Reveal className="mt-12">
            <figure className="mx-auto max-w-2xl">
              <HeroDiagram />
            </figure>
          </Reveal>

          <p className="mt-10 inline-block text-sm text-muted-foreground">
            ↓ 目次からお好きなページへ
          </p>
        </div>
      </section>

      <section className="px-6 py-16 sm:py-20" id="toc">
        <div className="mx-auto max-w-5xl">
          <Reveal className="mx-auto max-w-xl text-center">
            <h2 className="text-3xl font-extrabold sm:text-4xl">目次</h2>
            <p className="mt-3 text-lg font-medium">
              19の絵本ページで、DDDの考え方をひとつずつ見ていく。
            </p>
          </Reveal>

          <TocGrid />
        </div>
      </section>
    </main>
  );
}
