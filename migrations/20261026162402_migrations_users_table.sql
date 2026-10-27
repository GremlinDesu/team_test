-- +goose Up
CREATE TABLE users (
    id serial PRIMARY KEY,
    email varchar(255) unique not null,
    password_hash text not null,
    name varchar(255) not null,
    created_at timestamp not null default current_timestamp
);
-- +goose Down
DROP TABLE if exists users;
