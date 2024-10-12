package orderd

import "time"

type OrderEntity struct {
	ID          string    `json:"id" db:"id"`
	OrderDate   time.Time `json:"orderDate" db:"order_date"`
	TotalAmount float64   `json:"totalAmount" db:"total_amount"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type OrderItemEntity struct {
	ID        string    `json:"id" db:"id"`
	OrderID   string    `json:"order_id" db:"order_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	SKU       string    `db:"sku" json:"sku"`
	Name      string    `db:"name" json:"name"`
	Qty       int       `json:"qty" db:"qty"`
	Price     float64   `db:"price" json:"price"`
	Discount  float64   `json:"discount" db:"discount"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
