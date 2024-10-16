package promotiond

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PromotionRepo interface {
	WrapperTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error)
	GetListPromotion(ctx context.Context) (res []PromotionEntity, err error)
	GetListPromotionByProductID(ctx context.Context, listProductid []string) (res []ListPromotion, err error)
	InsertBatchOrderPromotion(ctx context.Context, tx *sqlx.Tx, params []OrderPromotionEntity) (err error)
	GetListOrderPromotion(ctx context.Context, orderid string) (res []OrderPromotionEntity, err error)
}
