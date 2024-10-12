package productu

import (
	"checkout-service/internal/helper"
	"context"
)

type ProductUsc interface {
	GetListProduct(ctx context.Context) (res []GetListProduct, err *helper.ErrorStruct)
}
