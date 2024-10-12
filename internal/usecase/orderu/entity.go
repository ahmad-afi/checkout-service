package orderu

import "time"

type OrderCreateReq struct {
	Data []OrderData `json:"data" validate:"required,dive"`
}

type OrderData struct {
	ProductID string `json:"productid" validate:"required"`
	Qty       int    `json:"qty" validate:"required,gte=1"`
}

type CheckOrderRes struct {
	TotalAmount float64 `json:"totalAmount"`
	Discount    float64 `json:"discount"`
}

type OrderHistory struct {
	ID          string    `json:"id" db:"id"`
	OrderDate   time.Time `json:"orderDate" db:"order_date"`
	TotalAmount float64   `json:"totalAmount" db:"total_amount"`
}

type OrderDetail struct {
	OrderHistory
	Data          []OrderDetailItem `json:"data"`
	PromotionList []string          `json:"promotionList"`
}

type OrderDetailItem struct {
	ID        string  `json:"id" db:"id"`
	OrderID   string  `json:"order_id" db:"order_id"`
	ProductID string  `json:"product_id" db:"product_id"`
	SKU       string  `db:"sku" json:"sku"`
	Name      string  `db:"name" json:"name"`
	Qty       int     `json:"qty" db:"qty"`
	Price     float64 `db:"price" json:"price"`
	Discount  float64 `json:"discount" db:"discount"`
}
