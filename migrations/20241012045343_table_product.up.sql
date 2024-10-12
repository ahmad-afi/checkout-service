CREATE TABLE products (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    qty INT NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

INSERT INTO products(id, name, price, sku, qty)
VALUES
    ('01HKBSKNHV6XYAF55NSAK940ZK', 'Google Home', 49.99, '120P90', 10),
    ('01HKBSM317D1K9JPBKSAT9QVY9', 'Macbook Pro', 5399.99, '43N23P', 5),
    ('01HKBSMH3S0ADPWX2D7QA9PAZY', 'Alexa Speaker', 109.50, 'A304SD', 10),
    ('01HKBSMBE8DVW9RVT6WBWWDNRS', 'Raspberry Pi B', 30.00, '234234', 2)
ON CONFLICT (id) DO NOTHING;
