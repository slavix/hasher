-- +goose Up
-- +goose StatementBegin
CREATE TABLE hashes
(
    id   serial      not null unique,
    hash varchar(64) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE hashes;
-- +goose StatementEnd
