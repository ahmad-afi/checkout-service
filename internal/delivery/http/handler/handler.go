package handler

import "checkout-service/internal/infrastructure/container"

type handler struct {
	OrderHandler   orderHandler
	ProductHandler productHandler
}

func SetupHandler(cont container.Container) handler {
	return handler{
		OrderHandler:   NewOrderHandler(cont.OrderUsc),
		ProductHandler: NewProductHandler(cont.ProductUsc),
	}
}
