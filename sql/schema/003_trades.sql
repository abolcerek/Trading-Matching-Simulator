-- +goose Up
CREATE TABLE trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    maker_order_id UUID REFERENCES orders(order_id),
    taker_order_id UUID REFERENCES orders(order_id),
    maker_user_id UUID REFERENCES users(id),
    taker_user_id UUID REFERENCES users(id),
    price BIGINT NOT NULL,
    quantity BIGINT NOT NULL,
    trade_seq BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE trades;