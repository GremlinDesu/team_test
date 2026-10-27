-- +goose Up
create table team_members (
    id serial primary key,
    team_id bigint unsigned not null references teams(id) on delete cascade,
    user_id bigint unsigned not null references users(id) on delete cascade,
    role varchar(255) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

-- +goose Down
drop table if exists team_members;
