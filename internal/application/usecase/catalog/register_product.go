package catalog

import (
	"context"

	"github.com/almondoo/golang-ddd-sample/internal/application/tx"
	domaincatalog "github.com/almondoo/golang-ddd-sample/internal/domain/catalog"
	"github.com/almondoo/golang-ddd-sample/internal/domain/shared"
)

// RegisterProductInput はユースケースの入力データである。
//
// アプリケーション層は HTTP のリクエストボディをそのまま受け取らず、
// この専用の入力型に一度変換してから使う。こうすることで、ユースケース
// は「HTTP から呼ばれた」という事実を知らずに済み、将来 CLI やバッチ
// から同じユースケースを呼び出す場合にも流用できる。
//
// 通貨は本サンプルでは日本円（shared.JPY）に固定している。多通貨対応が
// 必要になった場合は Currency フィールドを追加することになる。
type RegisterProductInput struct {
	Name        string
	Description string
	PriceAmount int64
}

// RegisterProductOutput はユースケースの出力データである。
type RegisterProductOutput struct {
	ProductID string
}

// RegisterProductUseCase は新規商品を登録するユースケースである。
type RegisterProductUseCase struct {
	repo      domaincatalog.Repository
	txManager tx.Manager
}

// NewRegisterProductUseCase は RegisterProductUseCase を生成する。
func NewRegisterProductUseCase(repo domaincatalog.Repository, txManager tx.Manager) *RegisterProductUseCase {
	return &RegisterProductUseCase{repo: repo, txManager: txManager}
}

// Execute は商品登録ユースケースを実行する。
//
// 「1 ユースケース = 1 トランザクション」という原則に従い、ここで
// txManager.Do を呼んでトランザクション境界を明示する。ユースケースが
// トランザクションの開始・終了を決める責務を持つのは、ユースケースこそが
// 「1 つの業務操作として、どこからどこまでを不可分に実行すべきか」を
// 知っている層だからである（リポジトリやドメインはこの範囲を知らない）。
func (uc *RegisterProductUseCase) Execute(ctx context.Context, input RegisterProductInput) (RegisterProductOutput, error) {
	price, err := shared.NewMoney(input.PriceAmount, shared.JPY)
	if err != nil {
		return RegisterProductOutput{}, err
	}

	product, err := domaincatalog.NewProduct(input.Name, input.Description, price)
	if err != nil {
		return RegisterProductOutput{}, err
	}

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		return uc.repo.Save(ctx, product)
	})
	if err != nil {
		return RegisterProductOutput{}, err
	}

	return RegisterProductOutput{ProductID: product.ID().String()}, nil
}
