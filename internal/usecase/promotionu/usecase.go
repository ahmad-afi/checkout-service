package promotionu

import (
	"checkout-service/internal/helper"
	"context"
)

type PromotionUsc interface {
	GetListPromotion(ctx context.Context) (res []GetListPromotion, err *helper.ErrorStruct)
}
