package http

import (
	"checkout-service/internal/delivery/http/handler"
	"checkout-service/internal/infrastructure/container"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(f *fiber.App, cont container.Container) {
	h := handler.SetupHandler(cont)

	f.Get("", healthCheck)
	v1api := f.Group("/v1")

	productGroup := v1api.Group("/product")
	{
		productGroup.Get("", h.ProductHandler.GetProductList)
	}
	promotionGroup := v1api.Group("/promotion")
	{
		promotionGroup.Get("", h.PromotionHandler.GetPromotionList)
	}
	orderGroup := v1api.Group("/order")
	{
		orderGroup.Post("/check", h.OrderHandler.CheckOrder)
		orderGroup.Post("/confirm", h.OrderHandler.CreateOrder)
		orderGroup.Get("", h.OrderHandler.HistoryOrder)
		orderGroup.Get("/:orderid", h.OrderHandler.OrderDetail)
	}
}

func healthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": "Server is up and running",
	})
}
