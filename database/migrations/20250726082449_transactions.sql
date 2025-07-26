-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TYPE transaction_type AS ENUM ('bet', 'win');

CREATE TABLE users (
    id SERIAL PRIMARY KEY
);

CREATE TABLE transactions (
    user_id SERIAL REFERENCES users (id),
    transaction_type transaction_type NOT NULL,
    amount INTEGER NOT NULL CHECK (amount > 0),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE transactions;
DROP TABLE users;

DROP TYPE transaction_type;

COMMIT;
-- +goose StatementEnd
