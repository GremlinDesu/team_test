-- +goose Up
CREATE TABLE teams (
    id serial primary key,
    name varchar(255) not null,
    created_by bigint unsigned not null references users(id) on delete cascade,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
-- +goose Down
drop table if exists teams;
