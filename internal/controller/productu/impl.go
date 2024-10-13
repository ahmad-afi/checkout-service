package productc

import (
	"context"
	"log"

	"checkout-service/internal/domain/productd"
)

type ProductController struct {
	productRepo productd.ProductRepo
}

func NewProductController(productRepo productd.ProductRepo) ProductUsc {
	return &ProductController{productRepo: productRepo}
}

func (u *ProductController) GetListProduct(ctx context.Context) (res []GetListProduct, err error) {
	resOrders, err := u.productRepo.GetListProduct(ctx, nil)
	if err != nil {
		log.Println("error at GetOrder : ", err.Error())
		return
	}

	for _, v := range resOrders {
		res = append(res, GetListProduct{
			ProductEntity: v,
		})
	}
	return
}
