/**
 * サイト内の全20ページ(トップ + 19個の概念ページ)の順序・タイトル等を
 * 一箇所で管理するデータ。TocGrid(目次)とPageNav(前へ/次へ)の
 * 両方がこの配列を単一の情報源として参照する。
 */

export type ConceptPage = {
  /** ルートセグメント。トップページは空文字列で表す */
  slug: string;
  /** 目次での通し番号(01〜19) */
  num: string;
  title: string;
  description: string;
  /** 目次をゆるく束ねる見出し(TocGridでのみ使用) */
  group: "ことばと地図" | "つくりかた" | "うごかしかた";
};

export const HOME_PAGE: ConceptPage = {
  slug: "",
  num: "00",
  title: "DDD絵本",
  description: "目次",
  group: "ことばと地図",
};

export const CONCEPT_PAGES: ConceptPage[] = [
  {
    slug: "domain",
    num: "01",
    title: "ドメイン",
    description: "ソフトの前に、業務の世界がある。",
    group: "ことばと地図",
  },
  {
    slug: "ubiquitous-language",
    num: "02",
    title: "ユビキタス言語",
    description: "みんな、同じ言葉で話す。",
    group: "ことばと地図",
  },
  {
    slug: "bounded-context",
    num: "03",
    title: "境界づけられたコンテキスト",
    description: "お店を小部屋に分ける。",
    group: "ことばと地図",
  },
  {
    slug: "context-map",
    num: "04",
    title: "コンテキストマップ",
    description: "つきあい方を、1枚の地図にする。",
    group: "ことばと地図",
  },
  {
    slug: "shared-kernel",
    num: "05",
    title: "共有カーネル",
    description: "みんなで使う共通の道具箱。",
    group: "ことばと地図",
  },
  {
    slug: "entity",
    num: "06",
    title: "エンティティ",
    description: "名前が変わってもあなたはあなた。",
    group: "つくりかた",
  },
  {
    slug: "value-object",
    num: "07",
    title: "値オブジェクト",
    description: "1000円はどの1000円でも同じ。",
    group: "つくりかた",
  },
  {
    slug: "aggregate",
    num: "08",
    title: "集約",
    description: "おもちゃ箱ごと出し入れする。",
    group: "つくりかた",
  },
  {
    slug: "domain-service",
    num: "09",
    title: "ドメインサービス",
    description: "集約をまたぐ計算の置き場所。",
    group: "つくりかた",
  },
  {
    slug: "repository",
    num: "10",
    title: "リポジトリ",
    description: "倉庫係にお願いするだけ。",
    group: "つくりかた",
  },
  {
    slug: "factory",
    num: "11",
    title: "ファクトリ",
    description: "正しい作り方を知っている工房。",
    group: "つくりかた",
  },
  {
    slug: "application-service",
    num: "12",
    title: "アプリケーションサービス",
    description: "段取り係が手順をまとめる。",
    group: "うごかしかた",
  },
  {
    slug: "domain-event",
    num: "13",
    title: "ドメインイベント",
    description: "「起きたこと」を知らせる仕組み。",
    group: "うごかしかた",
  },
  {
    slug: "cqrs",
    num: "14",
    title: "CQRS(軽量版)",
    description: "書く道と読む道を分ける。",
    group: "うごかしかた",
  },
  {
    slug: "query-service",
    num: "15",
    title: "クエリサービス",
    description: "見るだけなら、正式ルートはいらない。",
    group: "うごかしかた",
  },
  {
    slug: "onion-architecture",
    num: "16",
    title: "オニオンアーキテクチャ",
    description: "まん中を、まわりが守る。",
    group: "うごかしかた",
  },
  {
    slug: "optimistic-locking",
    num: "17",
    title: "楽観的ロック",
    description: "版数札つきの荷物。",
    group: "うごかしかた",
  },
  {
    slug: "specification",
    num: "18",
    title: "仕様(Specification)",
    description: "条件そのものを部品にする。",
    group: "うごかしかた",
  },
  {
    slug: "place-order",
    num: "19",
    title: "実例で見る: 注文確定の旅",
    description: "PlaceOrderの一部始終。",
    group: "うごかしかた",
  },
];

/** ページ送り(PageNav)で使う、index始まりの全ページ順序 */
export const PAGE_ORDER: ConceptPage[] = [HOME_PAGE, ...CONCEPT_PAGES];

export function pageHref(page: ConceptPage): string {
  return page.slug === "" ? "/" : `/${page.slug}`;
}

/**
 * 指定したslugの前後ページを返す。
 * トップページ(slug === "")の前は存在しない。
 */
export function getAdjacentPages(currentSlug: string): {
  prev: ConceptPage | null;
  next: ConceptPage | null;
} {
  const index = PAGE_ORDER.findIndex((p) => p.slug === currentSlug);
  if (index === -1) {
    return { prev: null, next: null };
  }
  return {
    prev: index > 0 ? PAGE_ORDER[index - 1] : null,
    next: index < PAGE_ORDER.length - 1 ? PAGE_ORDER[index + 1] : null,
  };
}
