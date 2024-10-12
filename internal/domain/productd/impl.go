package productd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductDomain struct {
	pg *sqlx.DB
}

func NewProductDomain(pg *sqlx.DB) ProductRepo {
	return &ProductDomain{pg: pg}
}

func (d *ProductDomain) GetListProduct(ctx context.Context, listid []string) (res []ProductEntity, err error) {
	var inputArgs []any
	query := `select id, sku, name, price, qty from products where deleted_at is null `

	if len(listid) > 0 {
		query += " and id in (?) "
		inputArgs = append(inputArgs, listid)
	}
	query, inputArgs, err = sqlx.In(query, inputArgs...)
	if err != nil {
		return nil, err
	}
	query = d.pg.Rebind(query)

	query += " order by name asc "
	err = d.pg.SelectContext(ctx, &res, query, inputArgs...)
	return
}

func (d *ProductDomain) GetListProductForUpdate(ctx context.Context, tx *sqlx.Tx, listid []string) (res []ProductEntity, err error) {
	var inputArgs []any
	query := `select id, sku, name, price, qty from products where deleted_at is null`
	if len(listid) > 0 {
		query += " and id in (?) "
		inputArgs = append(inputArgs, listid)
	}
	query, inputArgs, err = sqlx.In(query, inputArgs...)
	if err != nil {
		return nil, err
	}
	query = d.pg.Rebind(query)
	query += " order by name asc for update"
	err = tx.SelectContext(ctx, &res, query, inputArgs...)
	return
}

func (d *ProductDomain) UpdateProductQtyBatch(ctx context.Context, tx *sqlx.Tx, data []UpdateProduct) (err error) {
	var listValues = make([]string, 0, 10)
	var inputArgs []any
	for _, v := range data {
		listValues = append(listValues, "(?, ?)")
		inputArgs = append(inputArgs, v.ID, v.Qty)
	}
	query := fmt.Sprintf(`UPDATE products AS p
	SET qty = v.qty :: int, 
		updated_at = NOW()
	FROM (VALUES %s) AS v(id, qty)
	WHERE p.id = v.id;
	`, strings.Join(listValues, ","))

	query = tx.Rebind(query)

	res, err := tx.ExecContext(ctx, query, inputArgs...)
	if err != nil {
		return
	}

	affected, _ := res.RowsAffected()
	if affected < 1 {
		err = fmt.Errorf("failed to update data")
		return
	}
	return
}
