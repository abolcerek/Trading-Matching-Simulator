-- +goose Up
CREATE TABLE commands (
    Sequence_num BIGINT NOT NULL DEFAULT 0
);

INSERT INTO commands (Sequence_num)
VALUES (0);

-- +goose Down
DROP TABLE commands;