-- +goose Up
create table tasks (
    id serial primary key,
    team_id bigint unsigned not null references teams(id) on delete cascade,
    title varchar(255) not null,
    description text,
    status varchar(50) not null default 'pending',
    created_by bigint unsigned not null references users(id) on delete cascade,
    assignee_id bigint unsigned references users(id) on delete set null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    closed_at timestamp,
    version varchar(50) not null default "1"
);

-- +goose Down
drop table if exists tasks;
