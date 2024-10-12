package handler

import (
	"checkout-service/internal/helper"
	"checkout-service/internal/usecase/productu"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	productUsc productu.ProductUsc
}

func NewProductHandler(productUsc productu.ProductUsc) productHandler {
	return productHandler{productUsc: productUsc}
}

func (h *productHandler) GetProductList(ctx *fiber.Ctx) error {
	res, err := h.productUsc.GetListProduct(ctx.Context())
	if err != nil {
		return helper.BuildResponse(ctx, false, err.Message, nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, res, http.StatusOK)
}
