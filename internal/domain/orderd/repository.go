package orderd

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderRepo interface {
	WrapperTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error)
	CreateOrder(ctx context.Context, tx *sqlx.Tx, params OrderEntity) (err error)
	CreateOrderItems(ctx context.Context, tx *sqlx.Tx, params []OrderItemEntity) (err error)
	GetHistoryOrder(ctx context.Context, id string) (res []OrderEntity, err error)
	GetOrderItemsByOrderID(ctx context.Context, orderID string) (res []OrderItemEntity, err error)
}
