import type { Metadata } from "next";
import { ConceptLayout } from "@/components/concept-layout";
import { Reveal } from "@/components/reveal";
import { RepositoryDiagram } from "@/components/diagrams/repository-diagram";

export const metadata: Metadata = {
  title: "リポジトリとは",
  description:
    "倉庫係にお願いするだけ。リポジトリ(Repository)の考え方をGo DDDサンプルのcatalog.Repositoryで図解します。",
};

export default function RepositoryPage() {
  return (
    <ConceptLayout
      slug="repository"
      eyebrow="10 / リポジトリ"
      title="倉庫係に、お願いするだけ"
      lead="「保存して」「ID◯◯を持ってきて」。それだけ話せばいい。SQLは倉庫係だけの秘密。"
      example={{
        code: (
          <>
            <span className="text-background/50">
              {"// internal/domain/catalog/repository.go"}
            </span>
            {"\n"}
            {"type Repository interface {\n"}
            {"    "}
            <span className="text-background/50">{"// idに対応する商品を取得する"}</span>
            {"\n"}
            {"    FindByID(ctx context.Context, id ProductID) (*Product, error)\n\n"}
            {"    "}
            <span className="text-background/50">
              {"// 商品を永続化する(新規作成・更新のいずれも担うupsert)"}
            </span>
            {"\n"}
            {"    Save(ctx context.Context, p *Product) error\n"}
            {"}"}
          </>
        ),
        note: "呼び出す側はFindByID・Saveという2つの窓口だけを知っていればいい。GORMやSQLの実装は、インフラストラクチャ層にあるリポジトリの実装が隠す。",
      }}
      summary="リポジトリは「取ってきて」「しまって」を頼む窓口。しまい方(SQL)は聞かない。"
    >
      <Reveal>
        <figure className="mx-auto max-w-2xl">
          <RepositoryDiagram />
        </figure>
      </Reveal>
      <p className="mx-auto mt-5 max-w-xl text-center text-sm text-muted-foreground">
        ドメイン層はRepositoryという「インターフェース(窓口)」だけを定義する。実際にSQLを書く倉庫係(実装)はインフラストラクチャ層にいる。
      </p>
    </ConceptLayout>
  );
}
