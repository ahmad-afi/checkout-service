CREATE TABLE promotions (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL, -- discount, bundle, buy_x_pay_y
    "description" text NOT NULL,
    promotiondetail text not null, -- berbentuk object nanti
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

CREATE TABLE product_promotions (
    id VARCHAR(26) PRIMARY KEY,
    promotion_id VARCHAR(26) REFERENCES promotions(id) NOT NULL,
    product_id VARCHAR(26) REFERENCES products(id) NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);


-- for history
CREATE TABLE order_promotions (
    id VARCHAR(26) PRIMARY KEY,
    promotion_id VARCHAR(26) REFERENCES promotions(id) NOT NULL,
    order_id VARCHAR(26) REFERENCES orders(id) NOT NULL,
    name VARCHAR(255) NOT NULL,
    "refdata" text NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);
