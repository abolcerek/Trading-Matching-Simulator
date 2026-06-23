-- +goose Up
CREATE SEQUENCE command_seq;

-- +goose Down
DROP SEQUENCE command_seq;
