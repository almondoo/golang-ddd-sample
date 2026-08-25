import type { CSSProperties } from "react";
import type { Metadata } from "next";
import { Baloo_2, M_PLUS_Rounded_1c } from "next/font/google";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import "./globals.css";

// 見出し用: 英数字・記号向けのポップな丸ゴシック
const baloo2 = Baloo_2({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700", "800"],
  variable: "--font-baloo",
  display: "swap",
});

// 本文・日本語見出し用の丸ゴシック
// (next/font/googleの型定義上 M PLUS Rounded 1c は"japanese"サブセットを
// 選択できないが、日本語グリフ自体はフォントに含まれておりlatinのみの
// 指定でも問題なく表示される)
const mplusRounded = M_PLUS_Rounded_1c({
  subsets: ["latin"],
  weight: ["400", "500", "700", "800"],
  variable: "--font-mplus",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "DDDってなに? - 図解でわかるドメイン駆動設計",
    template: "%s - DDD絵本",
  },
  description:
    "Goで書かれたECサイトのサンプルコードを教材に、ドメイン駆動設計(DDD)の考え方を絵で説明する絵本サイトです。",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="ja"
      data-scroll-behavior="smooth"
      className={`${baloo2.variable} ${mplusRounded.variable} h-full`}
      style={
        {
          "--font-sans":
            "var(--font-mplus), 'Noto Sans JP', 'Hiragino Kaku Gothic ProN', 'Hiragino Sans', 'Yu Gothic', Meiryo, -apple-system, BlinkMacSystemFont, sans-serif",
          "--font-heading":
            "var(--font-baloo), var(--font-mplus), 'Noto Sans JP', 'Hiragino Kaku Gothic ProN', 'Hiragino Sans', 'Yu Gothic', Meiryo, -apple-system, BlinkMacSystemFont, sans-serif",
        } as CSSProperties
      }
    >
      <body className="flex min-h-full flex-col antialiased">
        <SiteHeader />
        {children}
        <SiteFooter />
      </body>
    </html>
  );
}
