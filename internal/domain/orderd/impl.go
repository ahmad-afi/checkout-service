package orderd

import (
	"checkout-service/internal/utils"
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderDomain struct {
	pg *sqlx.DB
	utils.SQLTransaction
}

func NewOrderDomain(pg *sqlx.DB, transaction utils.SQLTransaction) OrderRepo {
	return &OrderDomain{pg: pg, SQLTransaction: transaction}
}

func (d *OrderDomain) CreateOrder(ctx context.Context, tx *sqlx.Tx, params OrderEntity) (err error) {
	_, err = tx.NamedExecContext(ctx, `INSERT INTO orders (id, order_date, total_amount) VALUES(:id, :order_date, :total_amount)`, params)
	return
}

func (d *OrderDomain) CreateOrderItems(ctx context.Context, tx *sqlx.Tx, params []OrderItemEntity) (err error) {
	_, err = tx.NamedExecContext(ctx, `INSERT INTO order_items (id, order_id, product_id, sku, name, qty, price, discount) VALUES(:id, :order_id, :product_id, :sku, :name, :qty, :price, :discount)`, params)
	return
}

func (d *OrderDomain) GetHistoryOrder(ctx context.Context, id string) (res []OrderEntity, err error) {
	var inputArgs []any
	query := `select id, order_date, total_amount, created_at, updated_at from orders where deleted_at is null`

	if id != "" {
		query += " and id = ? "
		inputArgs = append(inputArgs, id)
	}

	query, inputArgs, err = sqlx.In(query, inputArgs...)
	if err != nil {
		return nil, err
	}
	query = d.pg.Rebind(query)
	query += " order by order_date desc "
	err = d.pg.SelectContext(ctx, &res, query, inputArgs...)
	return
}

func (d *OrderDomain) GetOrderItemsByOrderID(ctx context.Context, orderID string) (res []OrderItemEntity, err error) {
	query := `select id, order_id, product_id, sku, name, qty, price, discount, created_at, updated_at from order_items where deleted_at is null and order_id = $1 order by name`
	err = d.pg.SelectContext(ctx, &res, query, orderID)
	return
}

// func (d *OrderDomain) GetOrderByEmail(ctx context.Context, email string) (res OrderEntity, err error) {
// 	query := `select id,role_id, email, password, name, last_access, created_at, updated_at from users where deleted_at is null and email = $1`

// 	err = d.pg.GetContext(ctx, &res, query, email)
// 	return
// }

// func (d *OrderDomain) UpdateOrder(ctx context.Context, params OrderEntity) (err error) {
// 	query := `update users set name = $1,updated_at = now() where id = $2`

// 	_, err = d.pg.ExecContext(ctx, query, params.Name, params.ID)
// 	return
// }

// func (d *OrderDomain) UpdateLastAccess(ctx context.Context, id int) (err error) {
// 	query := `update users set last_access = now() where id = $1`

// 	_, err = d.pg.ExecContext(ctx, query, id)
// 	return
// }

// func (d *OrderDomain) DeleteOrder(ctx context.Context, id int) (err error) {
// 	query := `update users set deleted_at = now() where id = $1`

// 	_, err = d.pg.ExecContext(ctx, query, id)
// 	return
// }

// func (d *OrderDomain) GetRoleRight(ctx context.Context, roleid int) (res RoleRight, err error) {
// 	query := `SELECT id, role_id, route, "section", "path", r_create, r_read, r_update, r_delete, created_at, updated_at FROM role_rights WHERE role_id = $1`

// 	err = d.pg.GetContext(ctx, &res, query, roleid)
// 	return
// }
