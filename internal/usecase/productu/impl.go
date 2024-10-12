package productu

import (
	"context"

	"checkout-service/internal/domain/productd"
	"checkout-service/internal/helper"
)

type ProductUsecase struct {
	productRepo productd.ProductRepo
}

func NewProductUsecase(productRepo productd.ProductRepo) ProductUsc {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) GetListProduct(ctx context.Context) (res []GetListProduct, err *helper.ErrorStruct) {
	resOrders, errRepo := u.productRepo.GetListProduct(ctx, nil)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at GetListProduct", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	for _, v := range resOrders {
		res = append(res, GetListProduct{
			ProductEntity: v,
		})
	}
	return
}
