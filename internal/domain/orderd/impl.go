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
	_, err = tx.NamedExecContext(ctx, `INSERT INTO orders (id, order_date, total_amount, original_amount, total_discount) VALUES(:id, :order_date, :total_amount, :original_amount, :total_discount)`, params)
	return
}

func (d *OrderDomain) CreateOrderItems(ctx context.Context, tx *sqlx.Tx, params []OrderItemEntity) (err error) {
	_, err = tx.NamedExecContext(ctx, `INSERT INTO order_items (id, order_id, product_id, sku, name, qty, price, discount) VALUES(:id, :order_id, :product_id, :sku, :name, :qty, :price, :discount)`, params)
	return
}

func (d *OrderDomain) GetHistoryOrder(ctx context.Context, id string) (res []OrderEntity, err error) {
	var inputArgs []any
	query := `select id, order_date, total_amount, original_amount, total_discount, created_at, updated_at from orders where deleted_at is null`

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
