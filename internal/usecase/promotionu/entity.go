package promotionu

type GetListPromotion struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Type        string `json:"type" db:"type"` // discount, bundle, buy_x_pay_y
	Description string `json:"description" db:"description"`
}
