package promotiond

import "time"

// Promotion represents the promotions table
type PromotionEntity struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Type            string    `json:"type" db:"type"` // discount, bundle, buy_x_pay_y
	Description     string    `json:"description" db:"description"`
	PromotionDetail string    `json:"promotiondetail" db:"promotiondetail"` // This will be an object in JSON format
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type BundlePromotion struct {
	FreeItemProductID string `json:"freeItemProductID"`
	Threshold         int    `json:"threshold"`
	Getfree           int    `json:"getfree"`
}

type BuyPayPromotion struct {
	Buy    int `json:"buy"`
	PayFor int `json:"payFor"`
}

type DiscountPromotion struct {
	Type      string  `json:"type"` // fixed / percentage
	Threshold int     `json:"threshold"`
	Discount  float64 `json:"discount"`
}

// example
// INSERT INTO promotions (id, name,  description, type, promotiondetail)
// VALUES
// ('01JA0MPMD87DTKV9E4X129WEFN', 'MacBook Pro + Free Raspberry Pi B', 'Each sale of a MacBook Pro comes with a free Raspberry Pi B',
// 'bundle' , '{"freeItemProductID": "01HKBSMBE8DVW9RVT6WBWWDNRS", "threshold" : 1, "getfree": 1}'),
// ('01JA0MPMD8M9CKN7VSEPS5W6RM', 'Buy 3 Google Homes, Pay for 2', 'Buy 3 Google Homes, Pay for 2',
// 'buy_x_pay_y','{ "buy": 3, "payFor": 2}'),
// ('01JA0MPMD8GW32VG4FQKEXWRDG', '10% Off for More than 3 Alexa Speakers', 'Buying more than 3 Alexa Speakers will get a 10% discount on all Alexa speakers',
// 'discount', '{"type": "percentage", "threshold": 3, "discount": 10}');

// ProductPromotion represents the product_promotions table
type ProductPromotionEntity struct {
	ID          int       `json:"id" db:"id"`
	PromotionID string    `json:"promotion_id" db:"promotion_id"`
	ProductID   string    `json:"product_id" db:"product_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// INSERT INTO product_promotions (id,promotion_id, product_id)
// VALUES
// ('01JA0MPMD9KW7HZ7M7ZJNJKA2T', '01JA0MPMD87DTKV9E4X129WEFN', '01HKBSM317D1K9JPBKSAT9QVY9'),
// ('01JA0MPMD8ZAXE8HVBN1BNMYQE', '01JA0MPMD8M9CKN7VSEPS5W6RM', '01HKBSKNHV6XYAF55NSAK940ZK'),
// ('01JA0MPMD84NGPX9XZV1F4335H', '01JA0MPMD8GW32VG4FQKEXWRDG', '01HKBSMH3S0ADPWX2D7QA9PAZY');

// OrderPromotion represents the order_promotions table
type OrderPromotionEntity struct {
	ID          int       `json:"id" db:"id"`
	PromotionID string    `json:"promotion_id" db:"promotion_id"`
	OrderID     string    `json:"order_id" db:"order_id"`
	RefData     string    `json:"refdata" db:"refdata"` // Additional reference data as text
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ListPromotion struct {
	PromotionEntity
	ProductID string `json:"product_id" db:"product_id"`
}
