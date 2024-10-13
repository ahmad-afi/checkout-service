package handler

import (
	"checkout-service/internal/helper"
	"checkout-service/internal/usecase/orderu"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	orderusc orderu.OrderUsc
}

func NewOrderHandler(orderusc orderu.OrderUsc) orderHandler {
	return orderHandler{orderusc: orderusc}
}

func (h *orderHandler) CheckOrder(ctx *fiber.Ctx) error {
	data := new(orderu.OrderCreateReq)
	if err := ctx.BodyParser(data); err != nil {
		return helper.BuildResponse(ctx, false, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := h.orderusc.CheckOrder(ctx.Context(), *data)
	if err != nil {
		return helper.BuildResponse(ctx, false, err.Message, nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, res, http.StatusOK)
}

func (h *orderHandler) CreateOrder(ctx *fiber.Ctx) error {

	data := new(orderu.OrderCreateReq)
	if err := ctx.BodyParser(data); err != nil {
		return helper.BuildResponse(ctx, false, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := h.orderusc.CreateOrder(ctx.Context(), *data)
	if err != nil {
		return helper.BuildResponse(ctx, false, err.Message, nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, res, http.StatusCreated)
}

func (h *orderHandler) HistoryOrder(ctx *fiber.Ctx) error {
	res, err := h.orderusc.HistoryOrder(ctx.Context())
	if err != nil {
		return helper.BuildResponse(ctx, false, err.Message, nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, res, http.StatusOK)
}

func (h *orderHandler) OrderDetail(ctx *fiber.Ctx) error {
	orderid := ctx.Params("orderid")
	if orderid == "" {
		return helper.BuildResponse(ctx, false, "orderid cant be empty", nil, http.StatusBadRequest)
	}
	res, err := h.orderusc.OrderDetail(ctx.Context(), orderid)
	if err != nil {
		return helper.BuildResponse(ctx, false, err.Message, nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, res, http.StatusOK)
}
