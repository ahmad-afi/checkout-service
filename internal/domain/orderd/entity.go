package orderd

import "time"

type OrderEntity struct {
	ID             string    `json:"id" db:"id"`
	OrderDate      time.Time `json:"orderDate" db:"order_date"`
	TotalAmount    float64   `json:"totalAmount" db:"total_amount"`
	OriginalAmount float64   `json:"originalAmount" db:"original_amount"`
	TotalDiscount  float64   `json:"totalDiscount" db:"total_discount"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type OrderItemEntity struct {
	ID          string    `db:"id"`
	OrderID     string    `db:"order_id"`
	ProductID   string    `db:"product_id" `
	SKU         string    `db:"sku"`
	Name        string    `db:"name"`
	Qty         int       `db:"qty"`
	Price       float64   `db:"price"`
	TotalAmount float64   `db:"total_amount"`
	Discount    float64   `db:"discount"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
