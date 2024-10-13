package handler

import (
	"checkout-service/internal/helper"
	"checkout-service/internal/usecase/promotionu"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type promotionHandler struct {
	promotionUsc promotionu.PromotionUsc
}

func NewPromotionHandler(promotionUsc promotionu.PromotionUsc) promotionHandler {
	return promotionHandler{promotionUsc: promotionUsc}
}

func (h *promotionHandler) GetPromotionList(ctx *fiber.Ctx) error {
	res, err := h.promotionUsc.GetListPromotion(ctx.Context())
	if err != nil {
		return helper.BuildResponse(ctx, false, err.Message, nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, res, http.StatusOK)
}
