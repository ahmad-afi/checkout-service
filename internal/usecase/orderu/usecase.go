package orderu

import (
	"checkout-service/internal/helper"
	"context"
)

type OrderUsc interface {
	CheckOrder(ctx context.Context, params OrderCreateReq) (res CheckOrderRes, err *helper.ErrorStruct)
	CreateOrder(ctx context.Context, params OrderCreateReq) (res CheckOrderRes, err *helper.ErrorStruct)
	HistoryOrder(ctx context.Context) (res []OrderHistory, err *helper.ErrorStruct)
	OrderDetail(ctx context.Context, orderid string) (res OrderDetail, err *helper.ErrorStruct)
}
