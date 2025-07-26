-- +goose Up
-- +goose StatementBegin
BEGIN;

INSERT INTO users DEFAULT VALUES; -- 1
INSERT INTO users DEFAULT VALUES; -- 2
INSERT INTO users DEFAULT VALUES; -- 3
INSERT INTO users DEFAULT VALUES; -- 4
INSERT INTO users DEFAULT VALUES; -- 5
INSERT INTO users DEFAULT VALUES; -- 6
INSERT INTO users DEFAULT VALUES; -- 7
INSERT INTO users DEFAULT VALUES; -- 8
INSERT INTO users DEFAULT VALUES; -- 9
INSERT INTO users DEFAULT VALUES; -- 10

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DELETE FROM users;

COMMIT;
-- +goose StatementEnd
