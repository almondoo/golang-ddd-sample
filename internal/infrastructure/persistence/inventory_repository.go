package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/almondoo/golang-ddd-sample/internal/domain/inventory"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// StockRepository は inventory.Repository の GORM 実装である。
//
// アプリケーション層は inventory.Repository（インターフェース）にのみ
// 依存しており、この構造体（具体的な GORM 実装）の存在を知らない。
type StockRepository struct {
	db *gorm.DB
}

// NewStockRepository は StockRepository を生成する。
func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

// コンパイル時に StockRepository が inventory.Repository を満たすことを保証する。
var _ inventory.Repository = (*StockRepository)(nil)

// FindByProductID は指定商品の在庫を取得する。
func (r *StockRepository) FindByProductID(ctx context.Context, id inventory.ProductID) (*inventory.Stock, error) {
	db := DBFromContext(ctx, r.db)

	var model StockModel
	if err := db.WithContext(ctx).First(&model, "product_id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock for product %s: %w", id, shared.ErrNotFound)
		}
		return nil, err
	}

	return model.toDomain()
}

// Save は在庫を永続化する（新規作成・更新の両方を担う upsert）。
//
// Stock 集約は並行するトランザクション同士の競合（lost update）を
// 防ぐため、version カラムを使った楽観ロックで更新する。GORM の Save
// （主キーが存在すれば UPDATE、存在しなければ INSERT）をそのまま使わず、
// 新規作成と更新を明示的に分けているのはこのためである
// （catalog.ProductRepository.Save は楽観ロックを行わないため単純な Save で
// 済むが、Stock はここが異なる）。
//
// 以前の実装は「まず product_id の COUNT で存在確認し、結果に応じて
// INSERT/UPDATE を分岐する」という find-then-branch 方式だったが、これは
// TOCTOU（Time-Of-Check-Time-Of-Use）競合を引き起こしていた。存在確認から
// 実際の INSERT/UPDATE までの間に他のトランザクションが割り込むと、
//   - 2つのトランザクションがともに「存在しない」と判定して INSERT を
//     試み、2件目が unique 制約違反で失敗する（本来は 409 を返すべきだが
//     未分類のドライバエラーがそのまま漏れて 500 になっていた）。
//   - あるいは、findByProductID で NotFound を受けた側が NewStock で
//     version=1 の集約を新規作成した直後に、別のトランザクションが先に
//     INSERT で version=1 の行を作ってコミットしてしまうと、後発側の
//     存在確認は「存在する」と判定して UPDATE 分岐に進み、
//     WHERE version = 1 が偶然マッチしてしまう（NewStock も version=1
//     から始まるため）。その結果、本来検出すべき競合が
//     RowsAffected > 0 として通ってしまい、先発側の書き込みを黙って
//     上書きする（lost update が楽観ロックをすり抜ける）。
//
// これを避けるため、DB への問い合わせ（存在確認）を挟まず、集約自身が
// 持つ version の値だけで INSERT/UPDATE を判定する方式に変更した。
// NewStock は version=0（「未永続化」を表す規約、stock.go の version
// フィールドのコメントを参照）から始まり、ReconstructStock は DB から
// 読み込んだ version（常に 1 以上）をそのまま引き継ぐため、この判定は
// 集約のインメモリな状態だけで完結し、DB の現在状態を読みに行く必要が
// ない。「読んでから分岐する」という TOCTOU の温床そのものを取り除いて
// いる点が、単に存在確認とINSERT/UPDATEを同一トランザクションにする
// （それでも読み取りと書き込みの間に競合の余地が残る）のとは異なる。
//
//   - version == 0（未永続化）: version を 1 として INSERT する。
//     PostgreSQL の product_id 主キー制約により、同時に2つのトランザク
//     ションが同じ product_id を INSERT しようとした場合は片方が
//     一意制約違反で失敗する。この場合は shared.ErrConflict を返す
//     （詳細は insertNew のコメントを参照）。
//   - version >= 1（永続化済み）: UPDATE ... WHERE product_id = ? AND
//     version = ? という条件付き更新を行い、version は +1 する。この
//     Stock インスタンスが読み込まれた時点（s.Version()）から他の
//     トランザクションが先に更新していれば WHERE 句にマッチする行が
//     0 件になる（RowsAffected == 0）。この場合、行そのものが存在しない
//     のかバージョンが不一致なのかを厳密に切り分ける必要はなく、一律で
//     shared.ErrConflict を返す（読み込み時点から状態が変わっている、
//     という事実だけが呼び出し側にとって重要なため）。
func (r *StockRepository) Save(ctx context.Context, s *inventory.Stock) error {
	db := DBFromContext(ctx, r.db).WithContext(ctx)
	model := stockModelFromDomain(s)

	if s.Version() == 0 {
		return r.insertNew(db, &model)
	}

	result := db.Model(&StockModel{}).
		Where("product_id = ? AND version = ?", model.ProductID, s.Version()).
		Updates(map[string]any{
			"quantity": model.Quantity,
			"reserved": model.Reserved,
			"version":  s.Version() + 1,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("stock update conflict: product_id=%s: %w", model.ProductID, shared.ErrConflict)
	}
	return nil
}

// insertNew は未永続化（version == 0）の Stock を version=1 として INSERT する。
//
// product_id は主キー（unique 制約）であるため、同じ商品に対する INSERT が
// 同時に2件発生した場合は DB が一意制約違反で片方を拒否する。これは
// 「同じ商品の在庫行を2つのトランザクションが同時に新規作成しようとした」
// という楽観ロック競合の一種であるため、一律で shared.ErrConflict に
// 変換する（呼び出し側にとっては UPDATE 分岐の競合と同じく「読み直して
// リトライが必要」という結論になるため区別しない）。
//
// gorm.ErrDuplicatedKey による判定は、provideDB（persistence.NewDB）で
// gorm.Config.TranslateError を有効にしていることが前提となる。これが
// 無効だと *pgconn.PgError 等のドライバ固有エラーがそのまま返り、
// この判定に失敗して未分類のエラー（HTTP 500 相当）として扱われてしまう。
func (r *StockRepository) insertNew(db *gorm.DB, model *StockModel) error {
	model.Version = 1
	if err := db.Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("stock already exists: product_id=%s: %w", model.ProductID, shared.ErrConflict)
		}
		return err
	}
	return nil
}
