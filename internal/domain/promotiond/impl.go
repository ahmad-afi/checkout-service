package promotiond

import (
	"checkout-service/internal/utils"
	"context"

	"github.com/jmoiron/sqlx"
)

type PromotionDomain struct {
	pg *sqlx.DB
	utils.SQLTransaction
}

func NewPromotionDomain(pg *sqlx.DB, transaction utils.SQLTransaction) PromotionRepo {
	return &PromotionDomain{pg: pg, SQLTransaction: transaction}
}

func (d *PromotionDomain) GetListPromotion(ctx context.Context) (res []PromotionEntity, err error) {
	query := `select p.id,p.name, p.type, p.description, p.promotiondetail, p.created_at, p.updated_at from promotions p where deleted_at is null order by name`
	err = d.pg.SelectContext(ctx, &res, query)
	return
}

func (d *PromotionDomain) GetListPromotionByProductID(ctx context.Context, listProductid []string) (res []ListPromotion, err error) {
	var inputArgs []any
	query := `select p.id,p.name, pp.product_id, p.type, p.description, p.promotiondetail, p.created_at, p.updated_at from promotions p 
	join product_promotions pp
	on p.id = pp.promotion_id where p.deleted_at is null and pp.deleted_at is null `

	if len(listProductid) > 0 {
		query += " and pp.product_id in (?) "
		inputArgs = append(inputArgs, listProductid)
	}

	query, inputArgs, err = sqlx.In(query, inputArgs...)
	if err != nil {
		return nil, err
	}
	query = d.pg.Rebind(query)
	query += " order by p.created_at desc"
	err = d.pg.SelectContext(ctx, &res, query, inputArgs...)
	return
}

func (d *PromotionDomain) InsertBatchOrderPromotion(ctx context.Context, tx *sqlx.Tx, params []OrderPromotionEntity) (err error) {
	_, err = tx.NamedExecContext(ctx, `
	INSERT INTO order_promotions (
		id, 
		promotion_id, 
		order_id, 
		name, 
		refdata
	) VALUES (
		:id, 
		:promotion_id, 
		:order_id, 
		:name, 
		:refdata
	)`, params)
	return
}

func (d *PromotionDomain) GetListOrderPromotion(ctx context.Context, orderid string) (res []OrderPromotionEntity, err error) {
	query := `select id, 
	promotion_id, 
	order_id, 
	name, 
	refdata, 
	created_at, 
	updated_at
	from order_promotions p where deleted_at is null and order_id = $1 order by name`
	err = d.pg.SelectContext(ctx, &res, query, orderid)
	return
}
