package orderu

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"checkout-service/internal/domain/orderd"
	"checkout-service/internal/domain/productd"
	"checkout-service/internal/domain/promotiond"
	"checkout-service/internal/helper"
	"checkout-service/internal/utils"

	"github.com/jmoiron/sqlx"
)

type OrderUsecase struct {
	orderRepo     orderd.OrderRepo
	productRepo   productd.ProductRepo
	promotionRepo promotiond.PromotionRepo
}

func NewOrderUsecase(orderRepo orderd.OrderRepo, productRepo productd.ProductRepo, promotionRepo promotiond.PromotionRepo) OrderUsc {
	return &OrderUsecase{orderRepo, productRepo, promotionRepo}
}

func (u *OrderUsecase) CheckOrder(ctx context.Context, params OrderCreateReq) (res CheckOrderRes, err *helper.ErrorStruct) {
	if errValidate := utils.Validator(params); errValidate != nil {
		helper.Logger(helper.LoggerLevelError, "Error at len(listProductID) <0", errValidate.Err)
		err = helper.HelperErrorResponse(errValidate.Err, errValidate.Message)
		return
	}

	listProductID := []string{}
	mapProductQty := map[string]int{} // productid : qty
	var originalAmount float64 = 0
	var mapProductDetail = map[string]mapProduct{} // productid : productdetail

	for _, v := range params.Data {
		listProductID = append(listProductID, v.ProductID)

		// kalau ada data kembar dijadikan satu
		if totQty, exist := mapProductQty[v.ProductID]; exist {
			totQty += v.Qty
			mapProductQty[v.ProductID] = totQty
		} else {
			mapProductQty[v.ProductID] = v.Qty
		}
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
			err = helper.HelperErrorResponse(errRepo, errRepo.Error())
			return
		} else {
			totalAmount := helper.RoundNumber(float64(mapProductQty[v.ID])*v.Price, 2)

			mapProductDetail[v.ID] = mapProduct{
				ID:          v.ID,
				SKU:         v.SKU,
				Name:        v.Name,
				Price:       v.Price,
				Qty:         v.Qty,
				QtyToBuy:    mapProductQty[v.ID],
				TotalAmount: totalAmount,
			}
		}
	}

	// #3 Check Promotions
	resPromotion, errRepo := u.promotionRepo.GetListPromotionByProductID(ctx, listProductID)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at promotionRepo.GetListPromotionByProductID", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	discount, validPromotion, errCalculate := u.calculateDiscount(ctx, mapProductDetail, resPromotion)
	if errCalculate != nil {
		helper.Logger(helper.LoggerLevelError, "Error at calculateDiscount", errCalculate)
		err = helper.HelperErrorResponse(errCalculate)
		return
	}

	// newDataProduct, mungkin diatas ada penyesuaian lagi
	var newOrderData = make([]OrderDataRes, 0)
	for _, v := range mapProductDetail {
		// pehitungan original amount, siapa tau ada penambahan barang di case bundle
		originalAmount += v.TotalAmount

		newOrderData = append(newOrderData, OrderDataRes{
			ProductID:     v.ID,
			Qty:           v.QtyToBuy,
			Name:          v.Name,
			TotalAmount:   v.TotalAmount,
			TotalDiscount: v.TotalDiscount,
		})
	}
	sort.Slice(newOrderData, func(i, j int) bool {
		return newOrderData[i].Name < newOrderData[j].Name
	})

	res.PromotionList = make([]string, 0)
	for _, v := range validPromotion {
		res.PromotionList = append(res.PromotionList, v.Name)
	}

	res.Discount = helper.RoundNumber(discount, 2)
	res.OriginalAmount = helper.RoundNumber(originalAmount, 2)
	res.TotalAmount = helper.RoundNumber(originalAmount-discount, 2)
	res.Data = newOrderData
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
	var originalAmount, discount, totalAmount float64
	var mapProductDetail = map[string]mapProduct{} // productid : productdetail
	// untuk respon dan final data
	var newOrderData = make([]OrderDataRes, 0)

	for _, v := range params.Data {
		listProductID = append(listProductID, v.ProductID)

		// kalau ada data kembar dijadikan satu
		if totQty, exist := mapProductQty[v.ProductID]; exist {
			totQty += v.Qty
			mapProductQty[v.ProductID] = totQty
		} else {
			mapProductQty[v.ProductID] = v.Qty
		}
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

		// create orderid
		orderID, errIDGenerator := utils.IDGenerator()
		if errIDGenerator != nil {
			helper.Logger(helper.LoggerLevelError, "Error at errIDGenerator", errIDGenerator)
			err = helper.HelperErrorResponse(errIDGenerator)
			return
		}

		for _, v := range resProduct {
			if mapProductQty[v.ID] > v.Qty {
				errRepo = fmt.Errorf("invalid qty for product %s", v.Name)
				helper.Logger(helper.LoggerLevelError, "Error at Check valid qty", errRepo)
				err = helper.HelperErrorResponse(errRepo, errRepo.Error())
				return
			} else {
				// totalAmount := float64(mapProductQty[v.ID]) * v.Price
				totalAmount := helper.RoundNumber(float64(mapProductQty[v.ID])*v.Price, 2)
				mapProductDetail[v.ID] = mapProduct{
					ID:          v.ID,
					SKU:         v.SKU,
					Name:        v.Name,
					Price:       v.Price,
					Qty:         v.Qty,
					QtyToBuy:    mapProductQty[v.ID],
					TotalAmount: totalAmount,
				}
			}
		}

		// #3 Check Promotions
		resPromotion, errRepo := u.promotionRepo.GetListPromotionByProductID(ctx, listProductID)
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at promotionRepo.GetListPromotionByProductID", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		disc, validPromotion, errCalculate := u.calculateDiscount(ctx, mapProductDetail, resPromotion)
		if errCalculate != nil {
			helper.Logger(helper.LoggerLevelError, "Error at calculateDiscount", errCalculate)
			err = helper.HelperErrorResponse(errCalculate)
			return
		}

		res.PromotionList = make([]string, 0)
		for _, v := range validPromotion {
			res.PromotionList = append(res.PromotionList, v.Name)
		}

		discount = disc

		// #4 Update product
		// update data from discount

		listProductUpdate := make([]productd.UpdateProduct, 0)
		listProductItem := make([]orderd.OrderItemEntity, 0)
		for _, v := range mapProductDetail {
			newQty := v.Qty - v.QtyToBuy
			if newQty < 0 {
				errQty := fmt.Errorf("invalid qty for product %s")
				helper.Logger(helper.LoggerLevelError, "Error at qty", errQty)
				err = helper.HelperErrorResponse(errQty, errQty.Error())
				return
			}

			listProductUpdate = append(listProductUpdate, productd.UpdateProduct{
				ID:  v.ID,
				Qty: newQty,
			})

			// pehitungan original amount, siapa tau ada penambahan barang di case bundle
			originalAmount += v.TotalAmount

			newOrderData = append(newOrderData, OrderDataRes{
				ProductID:     v.ID,
				Qty:           v.QtyToBuy,
				Name:          v.Name,
				TotalAmount:   v.TotalAmount,
				TotalDiscount: v.TotalDiscount,
			})

			orderItem, errIDGenerator := utils.IDGenerator()
			if errIDGenerator != nil {
				helper.Logger(helper.LoggerLevelError, "Error at errIDGenerator orderItem", errIDGenerator)
				err = helper.HelperErrorResponse(errIDGenerator)
				return
			}
			listProductItem = append(listProductItem, orderd.OrderItemEntity{
				ID:          orderItem,
				OrderID:     orderID,
				ProductID:   v.ID,
				Qty:         v.QtyToBuy,
				Price:       v.Price,
				Discount:    v.TotalDiscount,
				SKU:         v.SKU,
				Name:        v.Name,
				TotalAmount: v.TotalAmount,
			})
		}
		sort.Slice(newOrderData, func(i, j int) bool {
			return newOrderData[i].Name < newOrderData[j].Name
		})
		errRepo = u.productRepo.UpdateProductQtyBatch(ctx, tx, listProductUpdate)
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at productRepo.UpdateProductQtyBatch", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}
		totalAmount = originalAmount - discount

		// #5 Create Order
		errRepo = u.orderRepo.CreateOrder(ctx, tx, orderd.OrderEntity{
			ID:             orderID,
			OrderDate:      time.Now(),
			TotalAmount:    originalAmount,
			OriginalAmount: totalAmount,
			TotalDiscount:  discount,
		})
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at orderRepo.CreateOrder", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		// #6 Create Order Item
		errRepo = u.orderRepo.CreateOrderItems(ctx, tx, listProductItem)
		if errRepo != nil {
			helper.Logger(helper.LoggerLevelError, "Error at orderRepo.CreateOrderItems", errRepo)
			err = helper.HelperErrorResponse(errRepo)
			return
		}

		// #7 insert into order promotion
		if len(validPromotion) > 0 {
			dataOrderPromotion := []promotiond.OrderPromotionEntity{}
			for _, v := range validPromotion {
				orderPromotionID, errIDGenerator := utils.IDGenerator()
				if errIDGenerator != nil {
					helper.Logger(helper.LoggerLevelError, "Error at errIDGenerator orderPromotionID", errIDGenerator)
					err = helper.HelperErrorResponse(errIDGenerator)
					return
				}

				byteRefdata, errMarshal := json.Marshal(v)
				if errMarshal != nil {
					helper.Logger(helper.LoggerLevelError, "Error at errMarshal orderPromotionID", errMarshal)
					err = helper.HelperErrorResponse(errMarshal)
					return
				}

				dataOrderPromotion = append(dataOrderPromotion, promotiond.OrderPromotionEntity{
					ID:          orderPromotionID,
					PromotionID: v.ID,
					OrderID:     orderID,
					Name:        v.Name,
					RefData:     string(byteRefdata),
				})
			}
			errRepo = u.promotionRepo.InsertBatchOrderPromotion(ctx, tx, dataOrderPromotion)
			if errRepo != nil {
				helper.Logger(helper.LoggerLevelError, "Error at promotionRepo.InsertBatchOrderPromotion", errRepo)
				err = helper.HelperErrorResponse(errRepo)
				return
			}
		}

		return
	})

	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.WrapperTransaction", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	res.Discount = 0
	res.Discount = helper.RoundNumber(discount, 2)
	res.OriginalAmount = helper.RoundNumber(originalAmount, 2)
	res.TotalAmount = helper.RoundNumber(totalAmount, 2)
	res.Data = newOrderData

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
			ID:             v.ID,
			OrderDate:      v.OrderDate,
			TotalAmount:    v.TotalAmount,
			OriginalAmount: v.OriginalAmount,
			TotalDiscount:  v.TotalDiscount,
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
		err = helper.HelperErrorResponse(errRepo, errRepo.Error())
		return
	}

	orderItem, errRepo := u.orderRepo.GetOrderItemsByOrderID(ctx, orderid)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at orderRepo.GetHistoryOrder", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	resOrderPromotion, errRepo := u.promotionRepo.GetListOrderPromotion(ctx, orderid)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at promotionRepo.GetListOrderPromotion", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	res.PromotionList = make([]string, 0)
	for _, v := range resOrderPromotion {
		res.PromotionList = append(res.PromotionList, v.Name)
	}

	res.ID = historyOrder[0].ID
	res.OrderDate = historyOrder[0].OrderDate
	res.TotalAmount = historyOrder[0].TotalAmount
	res.OriginalAmount = historyOrder[0].OriginalAmount
	res.TotalDiscount = historyOrder[0].TotalDiscount

	for _, v := range orderItem {
		res.Data = append(res.Data, OrderDetailItem{
			ID:          v.ID,
			OrderID:     v.OrderID,
			ProductID:   v.ProductID,
			SKU:         v.SKU,
			Name:        v.Name,
			Qty:         v.Qty,
			Price:       v.Price,
			Discount:    v.Discount,
			TotalAmount: v.TotalAmount,
		})
	}

	return
}

func (u *OrderUsecase) calculateDiscount(ctx context.Context, mapProductDetail map[string]mapProduct, dataPromotion []promotiond.ListPromotion) (discount float64, validPromotion []promotiond.ListPromotion, err error) {
	for _, v := range dataPromotion {
		switch v.Type {
		case "discount":
			var promotionDetail promotiond.DiscountPromotion
			err = json.Unmarshal([]byte(v.PromotionDetail), &promotionDetail)
			if err != nil {
				helper.Logger(helper.LoggerLevelError, "Error at json.Unmarsha discount", err)
				return
			}

			dataProductToBuy := mapProductDetail[v.ProductID]

			if dataProductToBuy.QtyToBuy >= promotionDetail.Threshold {
				// memasukan nama promo
				validPromotion = append(validPromotion, v)

				if promotionDetail.Type == "fixed" {
					disc := promotionDetail.Discount
					discount += disc

					dataProductToBuy.TotalDiscount = disc
					mapProductDetail[v.ProductID] = dataProductToBuy
				} else if promotionDetail.Type == "percentage" {
					disc := dataProductToBuy.TotalAmount * (promotionDetail.Discount / 100)
					disc = helper.RoundNumber(disc, 2)
					discount += disc

					dataProductToBuy.TotalDiscount = disc
					mapProductDetail[v.ProductID] = dataProductToBuy
				}
			}
		case "bundle":
			var promotionDetail promotiond.BundlePromotion
			err = json.Unmarshal([]byte(v.PromotionDetail), &promotionDetail)
			if err != nil {
				helper.Logger(helper.LoggerLevelError, "Error at json.Unmarsha bundle", err)
				return
			}

			dataProductToBuy := mapProductDetail[v.ProductID]
			if dataProductToBuy.QtyToBuy >= promotionDetail.Threshold {
				// memasukan nama promo
				// validPromotion = append(validPromotion, v)

				totalFreeProductQty := (dataProductToBuy.QtyToBuy / promotionDetail.Threshold) * promotionDetail.Getfree

				// get data free product
				dataProduct, errGetListProduct := u.productRepo.GetListProduct(ctx, []string{promotionDetail.FreeItemProductID})
				if errGetListProduct != nil {
					helper.Logger(helper.LoggerLevelError, "Error at GetListProduct", errGetListProduct)
				} else if len(dataProduct) < 1 {
					errGetListProduct = fmt.Errorf("data not found")
					helper.Logger(helper.LoggerLevelError, "Error at len(dataProduct) <1", errGetListProduct)
				} else {
					// kalau product nya habis, tidak dapat promo
					if dataProduct[0].Qty > 1 {
						// check apakah ada pembelian product free
						dataProductFree, exist := mapProductDetail[promotionDetail.FreeItemProductID]
						if exist {
							// case kalau product nya kurang
							if totalFreeProductQty > dataProduct[0].Qty {
								totalFreeProductQty = dataProduct[0].Qty
							}

							// kalau dia beli produknya lebih sedikit dari seharusnya promo
							// maka di pakasa jdi sesuai gratis promo
							if totalFreeProductQty > dataProductFree.QtyToBuy {
								dataProductFree.QtyToBuy = totalFreeProductQty
								dataProductFree.TotalAmount = dataProduct[0].Price * float64(totalFreeProductQty)
							}
							disc := dataProduct[0].Price * float64(totalFreeProductQty)
							disc = helper.RoundNumber(disc, 2)
							discount += disc
							dataProductFree.TotalDiscount = disc

							mapProductDetail[promotionDetail.FreeItemProductID] = dataProductFree
						} else {
							// case kalau product nya kurang
							if totalFreeProductQty > dataProduct[0].Qty {
								totalFreeProductQty = dataProduct[0].Qty
							}

							disc := dataProduct[0].Price * float64(totalFreeProductQty)
							disc = helper.RoundNumber(disc, 2)
							discount += disc
							mapProductDetail[promotionDetail.FreeItemProductID] = mapProduct{
								ID:            dataProduct[0].ID,
								SKU:           dataProduct[0].SKU,
								Name:          dataProduct[0].Name,
								Price:         dataProduct[0].Price,
								Qty:           dataProduct[0].Qty,
								QtyToBuy:      totalFreeProductQty,
								TotalAmount:   dataProduct[0].Price * float64(totalFreeProductQty),
								TotalDiscount: disc,
							}
						}

						// memasukan nama promo
						validPromotion = append(validPromotion, v)
					}
				}
			}

		case "buy_x_pay_y":
			var promotionDetail promotiond.BuyPayPromotion
			err = json.Unmarshal([]byte(v.PromotionDetail), &promotionDetail)
			if err != nil {
				helper.Logger(helper.LoggerLevelError, "Error at json.Unmarsha buy_x_pay_y", err)
				return
			}

			dataProductToBuy := mapProductDetail[v.ProductID]
			if dataProductToBuy.QtyToBuy >= promotionDetail.Buy {
				// memasukan nama promo
				validPromotion = append(validPromotion, v)

				// misal beli 3 dapat 2, maka dari pembagian berikut dpt hasil 1
				promoQty := dataProductToBuy.QtyToBuy / promotionDetail.Buy

				// maka dihitung 2 x harga normal
				var promoPrice float64 = float64(promoQty*promotionDetail.PayFor) * dataProductToBuy.Price

				// perhitungan sisa yg tidak termasuk promo
				remaining := dataProductToBuy.QtyToBuy % promotionDetail.Buy
				var remainingPrice float64 = float64(remaining) * dataProductToBuy.Price
				disc := dataProductToBuy.TotalAmount - (promoPrice + remainingPrice)
				disc = helper.RoundNumber(disc, 2)
				discount += disc

				dataProductToBuy.TotalDiscount = disc
				mapProductDetail[v.ProductID] = dataProductToBuy
			}
		}
	}

	return
}
