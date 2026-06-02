-- +goose Up
CREATE TABLE orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sequence_num BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    user_id UUID REFERENCES users(id),
    side TEXT NOT NULL CHECK (side IN ('buy', 'sell')),
    type TEXT NOT NULL CHECK (type IN ('market', 'limit')),
    price BIGINT CHECK ((type = 'limit' AND price IS NOT NULL) OR (type = 'market' AND price IS NULL)),
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    remaining_quantity BIGINT NOT NULL CHECK (remaining_quantity >= 0 AND remaining_quantity <= quantity),
    status TEXT NOT NULL CHECK (status IN ('pending', 'open', 'partially_filled', 'filled', 'canceled')),
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE orders;