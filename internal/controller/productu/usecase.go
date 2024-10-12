package productc

import (
	"context"
)

type ProductUsc interface {
	GetListProduct(ctx context.Context) (res []GetListProduct, err error)
}
