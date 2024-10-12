package productd

type ProductEntity struct {
	ID    string  `db:"id" json:"id"`
	SKU   string  `db:"sku" json:"sku"`
	Name  string  `db:"name" json:"name"`
	Price float64 `db:"price" json:"price"`
	Qty   int     `db:"qty" json:"qty"`
}

type UpdateProduct struct {
	ID  string `db:"id" json:"id"`
	Qty int    `db:"qty" json:"qty"`
}
