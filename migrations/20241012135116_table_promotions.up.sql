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


INSERT INTO promotions (id, name,  description, type, promotiondetail)
VALUES
('01JA0MPMD87DTKV9E4X129WEFN', 'MacBook Pro + Free Raspberry Pi B', 'Each sale of a MacBook Pro comes with a free Raspberry Pi B',
'bundle' , '{"freeItemProductID": "01HKBSMBE8DVW9RVT6WBWWDNRS", "threshold" : 1, "getfree": 1}'),
('01JA0MPMD8M9CKN7VSEPS5W6RM', 'Buy 3 Google Homes, Pay for 2', 'Buy 3 Google Homes, Pay for 2',
'buy_x_pay_y','{ "buy": 3, "payFor": 2}'),
('01JA0MPMD8GW32VG4FQKEXWRDG', '10% Off for More than 3 Alexa Speakers', 'Buying more than 3 Alexa Speakers will get a 10% discount on all Alexa speakers',
'discount', '{"type": "percentage", "threshold": 3, "discount": 10}');

INSERT INTO product_promotions (id,promotion_id, product_id)
VALUES 
('01JA0MPMD9KW7HZ7M7ZJNJKA2T', '01JA0MPMD87DTKV9E4X129WEFN', '01HKBSM317D1K9JPBKSAT9QVY9'),
('01JA0MPMD8ZAXE8HVBN1BNMYQE', '01JA0MPMD8M9CKN7VSEPS5W6RM', '01HKBSKNHV6XYAF55NSAK940ZK'),
('01JA0MPMD84NGPX9XZV1F4335H', '01JA0MPMD8GW32VG4FQKEXWRDG', '01HKBSMH3S0ADPWX2D7QA9PAZY');
