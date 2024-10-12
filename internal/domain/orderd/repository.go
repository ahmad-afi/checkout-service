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
	// GetOrder(ctx context.Context) (res []ListOrderEntity, err error)
	// GetOrderByEmail(ctx context.Context, email string) (res OrderEntity, err error)
	// UpdateOrder(ctx context.Context, params OrderEntity) (err error)
	// UpdateLastAccess(ctx context.Context, id int) (err error)
	// DeleteOrder(ctx context.Context, id int) (err error)

	// GetRoleRight(ctx context.Context, roleid int) (res RoleRight, err error)
}
