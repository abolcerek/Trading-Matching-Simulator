-- +goose Up
CREATE TABLE trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    maker_order_id UUID REFERENCES orders(order_id) NOT NULL,
    taker_order_id UUID REFERENCES orders(order_id) NOT NULL,
    maker_user_id UUID REFERENCES users(id) NOT NULL,
    taker_user_id UUID REFERENCES users(id) NOT NULL,
    price BIGINT NOT NULL,
    quantity BIGINT NOT NULL,
    trade_seq BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE trades;