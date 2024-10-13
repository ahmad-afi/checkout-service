package handler

import "checkout-service/internal/infrastructure/container"

type handler struct {
	OrderHandler     orderHandler
	ProductHandler   productHandler
	PromotionHandler promotionHandler
}

func SetupHandler(cont container.Container) handler {
	return handler{
		OrderHandler:     NewOrderHandler(cont.OrderUsc),
		ProductHandler:   NewProductHandler(cont.ProductUsc),
		PromotionHandler: NewPromotionHandler(cont.PromotionUsc),
	}
}
