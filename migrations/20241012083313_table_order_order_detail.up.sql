CREATE TABLE orders (
    id VARCHAR(26) PRIMARY KEY,
    order_date TIMESTAMP DEFAULT NOW(),
    total_amount DECIMAL(10, 2) NOT NULL,
    original_amount DECIMAL(10, 2) NOT NULL,
    total_discount DECIMAL(10, 2) NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

CREATE TABLE order_items (
    id VARCHAR(26) PRIMARY KEY,
    order_id VARCHAR(26) REFERENCES orders(id) NOT NULL,
    product_id VARCHAR(26) REFERENCES products(id) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    qty INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

ALTER TABLE order_items
    ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id)
        REFERENCES "orders" (id) on update cascade;
ALTER TABLE order_items
    ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id)
        REFERENCES "products" (id) on update cascade;