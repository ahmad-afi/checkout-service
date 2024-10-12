package orderu

import (
	"context"
	"fmt"
	"time"

	"checkout-service/internal/domain/orderd"
	"checkout-service/internal/domain/productd"
	"checkout-service/internal/helper"
	"checkout-service/internal/utils"

	"github.com/jmoiron/sqlx"
)

type OrderUsecase struct {
	orderRepo   orderd.OrderRepo
	productRepo productd.ProductRepo
}

func NewOrderUsecase(orderRepo orderd.OrderRepo, productRepo productd.ProductRepo) OrderUsc {
	return &OrderUsecase{orderRepo, productRepo}
}

func (u *OrderUsecase) CheckOrder(ctx context.Context, params OrderCreateReq) (res CheckOrderRes, err *helper.ErrorStruct) {
	if errValidate := utils.Validator(params); errValidate != nil {
		helper.Logger(helper.LoggerLevelError, "Error at len(listProductID) <0", errValidate.Err)
		err = helper.HelperErrorResponse(errValidate.Err, errValidate.Message)
		return
	}

	listProductID := []string{}
	mapProductQty := map[string]int{} // productid : qty
	var totalAmount float64 = 0

	for _, v := range params.Data {
		listProductID = append(listProductID, v.ProductID)
		mapProductQty[v.ProductID] = v.Qty
	}

	if len(listProductID) < 1 {
		errData := fmt.Errorf("invalid data product")
		helper.Logger(helper.LoggerLevelError, "Error at len(listProductID) <0", errData)
		err = helper.HelperErrorResponse(errData, errData.Error())
		return
	}

	// #1 Get List Data Product and for Locking
	resProduct, errRepo := u.productRepo.GetListProduct(ctx, listProductID)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetListProductForUpdate", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	// #2 Check valid qty
	for _, v := range resProduct {
		if mapProductQty[v.ID] > v.Qty {
			errRepo = fmt.Errorf("invalid qty for product %s", v.Name)
			helper.Logger(helper.LoggerLevelError, "Error at Check valid qty", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		} else {
			totalAmount += float64(mapProductQty[v.ID]) * v.Price
		}
	}

	res.Discount = 0
	res.TotalAmount = helper.RoundNumber(totalAmount, 2)

	return
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, params OrderCreateReq) (res CheckOrderRes, err *helper.ErrorStruct) {
	if errValidate := utils.Validator(params); errValidate != nil {
		helper.Logger(helper.LoggerLevelError, "Error at len(listProductID) <0", errValidate.Err)
		err = helper.HelperErrorResponse(errValidate.Err, errValidate.Message)
		return
	}
	listProductID := []string{}
	mapProductQty := map[string]int{} // productid : qty
	var totalAmount float64 = 0

	for _, v := range params.Data {
		listProductID = append(listProductID, v.ProductID)
		mapProductQty[v.ProductID] = v.Qty
	}

	if len(listProductID) < 1 {
		errData := fmt.Errorf("invalid data product")
		helper.Logger(helper.LoggerLevelError, "Error at len(listProductID) <0", errData)
		err = helper.HelperErrorResponse(errData)
		return
	}

	errRepo := u.orderRepo.WrapperTransaction(ctx, func(tx *sqlx.Tx) (errRepo error) {
		// #1 Get List Data Product and for Locking
		resProduct, errRepo := u.productRepo.GetListProductForUpdate(ctx, tx, listProductID)
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetListProductForUpdate", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		// #2 Check valid qty
		orderID, errIDGenerator := utils.IDGenerator()
		if errIDGenerator != nil {
			helper.Logger(helper.LoggerLevelError, "Error at errIDGenerator", errIDGenerator)
			err = helper.HelperErrorResponse(errIDGenerator)
			return
		}

		listProductUpdate := make([]productd.UpdateProduct, 0)
		listProductItem := make([]orderd.OrderItemEntity, 0)
		for _, v := range resProduct {
			if mapProductQty[v.ID] > v.Qty {
				errRepo = fmt.Errorf("invalid qty for product %s", v.Name)
				helper.Logger(helper.LoggerLevelError, "Error at Check valid qty", errRepo)
				err = helper.HelperErrorResponse(errRepo)
				return
			} else {

				newQty := v.Qty - mapProductQty[v.ID]
				listProductUpdate = append(listProductUpdate, productd.UpdateProduct{
					ID:  v.ID,
					Qty: newQty,
				})

				totalAmount += float64(mapProductQty[v.ID]) * v.Price
				orderItem, errIDGenerator := utils.IDGenerator()
				if errIDGenerator != nil {
					helper.Logger(helper.LoggerLevelError, "Error at errIDGenerator", errIDGenerator)
					err = helper.HelperErrorResponse(errIDGenerator)
					return
				}
				listProductItem = append(listProductItem, orderd.OrderItemEntity{
					ID:        orderItem,
					OrderID:   orderID,
					ProductID: v.ID,
					Qty:       mapProductQty[v.ID],
					Price:     v.Price,
					Discount:  0,
					SKU:       v.SKU,
					Name:      v.Name,
				})
			}
		}

		// #3 Update product
		errRepo = u.productRepo.UpdateProductQtyBatch(ctx, tx, listProductUpdate)
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at productRepo.UpdateProductQtyBatch", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		// #4 Create Order
		errRepo = u.orderRepo.CreateOrder(ctx, tx, orderd.OrderEntity{
			ID:          orderID,
			OrderDate:   time.Now(),
			TotalAmount: totalAmount,
		})
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at orderRepo.CreateOrder", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		// #4 Create Order Item
		errRepo = u.orderRepo.CreateOrderItems(ctx, tx, listProductItem)
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at orderRepo.CreateOrderItems", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		return
	})

	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.WrapperTransaction", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	res.Discount = 0
	res.TotalAmount = helper.RoundNumber(totalAmount, 2)

	return
}

func (u *OrderUsecase) HistoryOrder(ctx context.Context) (res []OrderHistory, err *helper.ErrorStruct) {
	listHistoryOrder, errRepo := u.orderRepo.GetHistoryOrder(ctx, "")
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetHistoryOrder", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	res = make([]OrderHistory, 0)
	for _, v := range listHistoryOrder {
		res = append(res, OrderHistory{
			ID:          v.ID,
			OrderDate:   v.OrderDate,
			TotalAmount: v.TotalAmount,
		})
	}
	return
}

func (u *OrderUsecase) OrderDetail(ctx context.Context, orderid string) (res OrderDetail, err *helper.ErrorStruct) {
	historyOrder, errRepo := u.orderRepo.GetHistoryOrder(ctx, orderid)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetHistoryOrder", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	if len(historyOrder) < 1 {
		errRepo = fmt.Errorf("orderid invalid")
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetHistoryOrder", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	orderItem, errRepo := u.orderRepo.GetOrderItemsByOrderID(ctx, orderid)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetHistoryOrder", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	res.ID = historyOrder[0].ID
	res.OrderDate = historyOrder[0].OrderDate
	res.TotalAmount = historyOrder[0].TotalAmount
	res.PromotionList = make([]string, 0)

	for _, v := range orderItem {
		res.Data = append(res.Data, OrderDetailItem{
			ID:        v.ID,
			OrderID:   v.OrderID,
			ProductID: v.ProductID,
			SKU:       v.SKU,
			Name:      v.Name,
			Qty:       v.Qty,
			Price:     v.Price,
			Discount:  v.Discount,
		})
	}

	return
}
