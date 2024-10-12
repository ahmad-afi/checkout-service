package productd

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	GetListProduct(ctx context.Context, listid []string) (res []ProductEntity, err error)
	GetListProductForUpdate(ctx context.Context, tx *sqlx.Tx, listid []string) (res []ProductEntity, err error)
	UpdateProductQtyBatch(ctx context.Context, tx *sqlx.Tx, data []UpdateProduct) (err error)
}
