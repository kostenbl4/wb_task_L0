-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature TEXT,
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(50) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL
);

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(50) NOT NULL,
    request_id VARCHAR(50),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(50) NOT NULL,
    status INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
