package promotionu

import (
	"context"

	"checkout-service/internal/domain/promotiond"
	"checkout-service/internal/helper"
)

type PromotionUsecase struct {
	promtionRepo promotiond.PromotionRepo
}

func NewPromotionUsecase(promtionRepo promotiond.PromotionRepo) PromotionUsc {
	return &PromotionUsecase{promtionRepo: promtionRepo}
}

func (u *PromotionUsecase) GetListPromotion(ctx context.Context) (res []GetListPromotion, err *helper.ErrorStruct) {
	resOrders, errRepo := u.promtionRepo.GetListPromotion(ctx)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, "Error at GetListPromotion", errRepo)
		err = helper.HelperErrorResponse(errRepo)
		return
	}

	for _, v := range resOrders {
		res = append(res, GetListPromotion{
			ID:          v.ID,
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
		})
	}
	return
}
