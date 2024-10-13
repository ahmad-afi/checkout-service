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
	OriginalAmount float64        `json:"originalAmount"`
	TotalAmount    float64        `json:"totalAmount"`
	Discount       float64        `json:"discount"`
	Data           []OrderDataRes `json:"data" validate:"required,dive"`
	PromotionList  []string       `json:"promotionList"`
}

type OrderDataRes struct {
	ProductID     string  `json:"productid"`
	Qty           int     `json:"qty"`
	Name          string  `json:"name"`
	TotalAmount   float64 `json:"totalAmount"`
	TotalDiscount float64 `json:"totalDiscount"`
}

type OrderHistory struct {
	ID             string    `json:"id"`
	OrderDate      time.Time `json:"orderDate"`
	OriginalAmount float64   `json:"originalAmount"`
	TotalAmount    float64   `json:"totalAmount"`
	TotalDiscount  float64   `json:"totalDiscount"`
}

type OrderDetail struct {
	OrderHistory
	Data          []OrderDetailItem `json:"data"`
	PromotionList []string          `json:"promotionList"`
}

type OrderDetailItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   string  `json:"product_id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Qty         int     `json:"qty"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	TotalAmount float64 `json:"total_amount"`
}

type mapProduct struct {
	ID            string  `json:"id"`
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Qty           int     `json:"qty"`
	QtyToBuy      int     `json:"qtyToBuy"`
	TotalAmount   float64 `json:"totalAmount"`
	TotalDiscount float64 `json:"totalDiscount"`
}
